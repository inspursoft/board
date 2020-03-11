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
	"time"

	"dao"

	_ "github.com/mattn/go-sqlite3/driver"
)

var (
	limitSize      int
	port           int
	affinityConfig = "affinity.ini"
)

func initLastAffinity() {
	err := dao.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize registry DB, error: %+v\n", err)
	}
	err = dao.DeleteLastAffinity()
	if err != nil {
		log.Printf("Failed to delete last affinity of KVMs, error: %+v\n", err)
	}
}

func loadKVMAffinity() {
	if _, err := os.Stat(affinityConfig); os.IsNotExist(err) {
		log.Fatalf("KVM affinity config not exists, error: %+v\n", err)
	}
	fh, err := os.Open(affinityConfig)
	if err != nil {
		log.Fatalf("Failed to open affinity config file, error: %+v\n", err)
	}
	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "=")
		err = dao.AddOrUpdateKVM(parts[0], parts[1])
		if err != nil {
			log.Printf("Failed to add or update KVM, error: %+v", err)
		}
	}
	output()
}

func output() {
	allKVMs, err := dao.GetKVMStatus()
	if err != nil {
		log.Printf("Failed to get available KVMs, error: %+v\n", err)
		return
	}
	for _, s := range allKVMs {
		log.Printf("%+v", s)
	}
	log.Println("=============================")
}

func initilize() {
	log.Println("Start loading and updating affinity configs ...")
	initLastAffinity()
	loadKVMAffinity()
}

func register(jobName string, affinity string) string {
	availableNodes := availableKVMWithAffinity(affinity)
	for len(availableNodes) == 0 {
		time.Sleep(time.Second * 2)
		log.Println("All nodes has been allocated, no available one currently.")
		availableNodes = availableKVMWithAffinity(affinity)
	}

	allocatableNodes := []string{}
	for _, registry := range availableNodes {
		if kvmName, allocated := registry[jobName]; allocated {
			return kvmName
		} else {
			allocatableNodes = append(allocatableNodes, registry["kvm_name"])
		}
	}

	kvmName := allocatableNodes[rand.Intn(len(allocatableNodes))]
	err := dao.AddOrUpdateRegistry(kvmName, jobName)
	if err != nil {
		log.Printf("Failed to add or update registry: %+v\n", err)
	}
	output()
	return kvmName
}

func updateRegistry(buildID, kvmName, jobName string) {
	err := dao.UpdateRegistry(buildID, kvmName, jobName)
	if err != nil {
		log.Printf("Failed to update registry: %+v with kvmName: %s, job name: %s, build ID: %s\n", err, kvmName, jobName, buildID)
	}
}

func availableKVMWithAffinity(affinity string) []map[string]string {
	results, err := dao.GetAvailableKVMWithAffinity(affinity)
	if err != nil {
		log.Printf("Failed to get available KVMs with affinity %s, error: %+v\n", affinity, err)
		return nil
	}
	return results
}

func allAvailables() []map[string]string {
	results, err := dao.GetKVMStatus()
	if err != nil {
		log.Printf("Failed to get available KVMs, error: %+v\n", err)
		return nil
	}
	return results
}

func unregister(jobName string, buildID string) (kvmName string) {
	kvmName, err := dao.GetKVMByJob(jobName)
	if err != nil {
		log.Printf("Failed to get KVM by job: %s, error: %+v\n", jobName, err)
	}
	err = dao.DeleteRelationalJob(jobName, buildID)
	if err != nil {
		log.Printf("Failed to delete relational job: %s, build ID: %s, error: %+v\n", jobName, buildID, err)
	}
	output()
	return
}

func releaseRegistryWithJob(kvmName string, jobName string) string {
	err := dao.DeleteRegistryWithJob(kvmName, jobName)
	if err != nil {
		log.Printf("Failed to delete registry with KVM: %s, job: %s, error: %+v\n", kvmName, jobName, err)
	}
	output()
	return kvmName
}

func renderOutput(response http.ResponseWriter, statusCode int, input interface{}, contentType ...string) error {
	if statusCode != http.StatusOK {
		response.WriteHeader(statusCode)
		return nil
	}
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
		affinity := request.FormValue("affinity")
		if jobName == "" {
			customAbort(response, http.StatusBadRequest, fmt.Errorf("no job name provided"))
			return
		}
		if strings.TrimSpace(affinity) == "" {
			affinity = "golang"
		}
		renderText(response, register(jobName, affinity))
	}
}

func updateWithBuild(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPut {
		kvmName := request.FormValue("kvm_name")
		jobName := request.FormValue("job_name")
		buildID := request.FormValue("build_id")
		if kvmName == "" || jobName == "" || buildID == "" {
			customAbort(response, http.StatusBadRequest, fmt.Errorf("no KVM name, job name or build ID provided"))
			return
		}
		updateRegistry(buildID, kvmName, jobName)
	}
}

func getJob(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		registerJob(response, request)
	}
}

func getAvailableNodes(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		renderJSON(response, allAvailables())
	}
}

func releaseNode(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		jobName := request.FormValue("job_name")
		buildID := request.FormValue("build_id")
		if jobName == "" || buildID == "" {
			customAbort(response, http.StatusBadRequest, fmt.Errorf("no job name or build ID provided"))
			return
		}
		renderText(response, unregister(jobName, buildID))
	}
}

func releaseKVM(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		kvmName := request.FormValue("kvm_name")
		jobName := request.FormValue("job_name")
		if kvmName == "" || jobName == "" {
			customAbort(response, http.StatusBadRequest, fmt.Errorf("no KVM name or job name provided"))
			return
		}
		renderText(response, releaseRegistryWithJob(kvmName, jobName))
	}
}

func triggerScript(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		scriptName := request.FormValue("name")
		arg := request.FormValue("arg")
		args := strings.Split(arg, ",")
		log.Printf("Executing script name: %s, args: %+v", scriptName, args)
		renderText(response, fmt.Sprintf("Triggerring script: %s with args: %+v...\n", scriptName, args))
		defer func() {
			go func() {
				time.Sleep(2 * time.Second)
				_, err := executeScripts(scriptName, args...)
				if err != nil {
					internalError(response, err)
					return
				}
			}()
		}()
	}
}

func executeScripts(scriptName string, args ...string) (string, error) {
	if _, err := os.Stat(scriptName); os.IsNotExist(err) {
		return fmt.Sprintf("Script file %s not found.", scriptName), err
	}
	cmd := exec.Command("chmod", "+x", scriptName)
	totalCommands := []string{scriptName}
	totalCommands = append(totalCommands, args...)
	cmd = exec.Command("sh", totalCommands...)
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
	http.HandleFunc("/update-build", updateWithBuild)
	http.HandleFunc("/get-job", getJob)
	http.HandleFunc("/available-nodes", getAvailableNodes)
	http.HandleFunc("/release-node", releaseNode)
	http.HandleFunc("/release-kvm", releaseKVM)
	http.HandleFunc("/trigger-script", triggerScript)

	l, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(l, nil)
}
