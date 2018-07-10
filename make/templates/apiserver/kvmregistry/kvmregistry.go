package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
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

func outputJSON(response http.ResponseWriter, input interface{}, contentType ...string) error {
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
	return nil
}

func registerJob(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		jobName := request.FormValue("job_name")
		outputJSON(response, register(jobName), "text/plain")
	}
}

func getJob(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		registerJob(response, request)
	}
}

func getAvailableNodes(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		outputJSON(response, availables())
	}
}

func getAllNodes(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		outputJSON(response, storage)
	}
}

func releaseNode(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		kvmName := request.FormValue("job_name")
		outputJSON(response, unregister(kvmName), "text/plain")
	}
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

	l, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(l, nil)
}
