package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	limitSize     int
	port          int
	kvmNamePrefix = "kvm-"
	storeFile     = "store.dat"
)

type node struct {
	KvmName string `json:"kvm_name"`
	JobName string `json:"job_name"`
	InUsed  bool   `json:"in_used"`
}

func (n *node) String() string {
	if n.JobName != "" {
		return fmt.Sprintf("Node: %s with job: %s, in-used: %v", n.KvmName, n.JobName, n.InUsed)
	}
	return fmt.Sprintf("Node: %s was not in-used.", n.KvmName)
}

var storage = make(map[string]*node)

func sync(process func(fh *os.File), accessOpt int) error {
	fh, err := os.OpenFile(storeFile, accessOpt, 0660)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Store file does not exist, create it ...")
			_, err = os.Create(storeFile)
			process(fh)
			return err
		}
		log.Printf("Failed to open file: %+v", err)
		return err
	}
	defer fh.Close()
	process(fh)
	return nil
}

func persist(fh *os.File) {
	w := bufio.NewWriter(fh)
	for n, s := range storage {
		w.WriteString(fmt.Sprintf("%s=%s\n", n, s.JobName))
	}
	w.Flush()
	output()
}

func load(fh *os.File) {
	if _, err := os.Stat(storeFile); !os.IsNotExist(err) {
		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), "=")
			nodeName := parts[0]
			jobName := parts[1]
			if n, ok := storage[nodeName]; jobName != "" && ok {
				n.KvmName = nodeName
				n.JobName = jobName
				n.InUsed = true
			}
		}
	}
	output()
}

func output() {
	for _, s := range storage {
		log.Printf("%+v", s)
	}
	log.Println("=============================")
}

func initilize() {
	log.Printf("Start initializing ...")
	for i := 0; i < limitSize; i++ {
		kvmName := kvmNamePrefix + strconv.Itoa(i+1)
		storage[kvmName] = &node{kvmName, "", false}
	}
	sync(load, os.O_RDONLY)

}

func register(jobName string) (nodeName string) {
	var hasRegistered bool
	for n, s := range storage {
		if s.JobName == jobName && s.InUsed {
			nodeName = n
			hasRegistered = true
			continue
		}
	}
	if !hasRegistered {
		availableNodes := availables()
		availableCount := len(availableNodes)
		if availableCount == 0 {
			nodeName = "FULL"
			log.Println("All nodes has been allocated, no available one currently.")
		} else {
			nodeName = availableNodes[rand.Intn(len(availableNodes))]
			storage[nodeName] = &node{KvmName: nodeName, JobName: jobName, InUsed: true}
			sync(persist, os.O_TRUNC|os.O_WRONLY)
		}
	}
	log.Printf("Current available nodes: %v", availables())
	return
}

func availables() []string {
	results := []string{}
	for n, s := range storage {
		if s == nil || s.JobName == "" || !s.InUsed {
			results = append(results, n)
		}
	}
	return results
}

func unregister(jobName string) (kvmName string) {
	kvmName = "NULL"
	for n, s := range storage {
		if s.InUsed && s.JobName == jobName {
			s.KvmName = n
			s.JobName = ""
			s.InUsed = false
			kvmName = n
			sync(persist, os.O_TRUNC|os.O_WRONLY)
			return
		}
	}
	return
}

func renderOutput(response http.ResponseWriter, statusCode int, input interface{}, contentType ...string) error {
	header := response.Header()
	if len(contentType) > 0 {
		header.Set("Context-type", contentType[0])
	} else {
		header.Set("Content-type", "application/json")
	}
	if data, ok := input.(string); ok {
		response.Write([]byte(data))
	} else {
		data, err := json.Marshal(input)
		if err != nil {
			log.Printf("Failed to marshal storage data: %+v", err)
			return err
		}
		response.Write(data)
	}
	if statusCode != http.StatusOK {
		response.WriteHeader(statusCode)
	}
	return nil
}

func renderJSON(resp http.ResponseWriter, data interface{}) {
	renderOutput(resp, http.StatusOK, data, "application/json")
}

func renderText(resp http.ResponseWriter, content string) {
	renderOutput(resp, http.StatusOK, content, "text/plain")
}

func internalError(resp http.ResponseWriter, err error) {
	renderOutput(resp, http.StatusInternalServerError, err.Error(), "text/plain")
}

func customAbort(resp http.ResponseWriter, statusCode int, err error) {
	renderOutput(resp, statusCode, err.Error(), "text/plain")
}

func registerJob(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		jobName := request.FormValue("job_name")
		renderText(response, register(jobName))
	}
}

func getJob(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		registerJob(response, request)
	}
}

func getAvailableNodes(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		renderJSON(response, availables())
	}
}

func getAllNodes(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		renderJSON(response, storage)
	}
}

func releaseNode(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		kvmName := request.FormValue("job_name")
		renderText(response, unregister(kvmName))
	}
}

func triggerScript(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		scriptName := request.FormValue("name")
		arg := request.FormValue("arg")
		output, err := executeScripts(scriptName, arg)
		if err != nil {
			internalError(response, err)
			return
		}
		renderText(response, output)
	}
}

func executeScripts(scriptName string, arg string) (string, error) {
	if _, err := os.Stat(scriptName); os.IsNotExist(err) {
		return fmt.Sprintf("Script file %s not found.", scriptName), err
	}
	cmd := exec.Command("chmod", "+x", scriptName)
	cmd = exec.Command("sh", scriptName, arg)
	var stdOutput bytes.Buffer
	cmd.Stdout = &stdOutput
	err := cmd.Run()
	if err != nil {
		return "Script executing failed.", err
	}
	return stdOutput.String(), nil
}

func main() {
	flag.IntVar(&limitSize, "size", 5, "Limit size of KVM node allocation.")
	flag.IntVar(&port, "port", 8899, "Specify the KVM-registry server port.")
	flag.Parse()
	log.Printf("KVM-registry limit size is %d.", limitSize)
	log.Printf("KVM-registry server listened on port %d.", port)
	initilize()
	http.HandleFunc("/register-job", registerJob)
	http.HandleFunc("/get-job", getJob)
	http.HandleFunc("/available-nodes", getAvailableNodes)
	http.HandleFunc("/nodes", getAllNodes)
	http.HandleFunc("/release-node", releaseNode)
	http.HandleFunc("/trigger-script", triggerScript)

	l, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(l, nil)
}
