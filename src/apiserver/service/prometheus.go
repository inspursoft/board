package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type DashboardInfo struct {
	IsOverMaxLimit  bool          `json:"is_over_max_limit"`
	IsOverMinLimit  bool          `json:"is_over_min_limit"`
	NodeCount       int           `json:"node_count"`
	NodeListData    []NodeList    `json:"node_list_data"`
	ServiceCount    int           `json:"service_count"`
	ServiceListData []ServiceList `json:"service_list_data"`
	TimeUnit        string        `json:"time_unit"`
}

type NodeList struct {
	Name         string     `json:"name"`
	NodeLogsData []NodeLogs `json:"node_logs_data"`
}

type NodeLogs struct {
	TimeStamp    int64   `json:"timestamp"`
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
	StorageTotal int     `json:"storage_total"`
	StorageUsed  int     `json:"storage_used"`
}

type ServiceList struct {
	Name            string        `json:"name"`
	ServiceLogsData []ServiceLogs `json:"service_logs_data"`
}

type ServiceLogs struct {
	TimeStamp       int64 `json:"timestamp"`
	ContainerNumber int   `json:"container_number"`
	PodNumber       int   `json:"pod_number"`
}

type RequestPayload struct {
	TimeStamp int64  `json:"timestamp"`
	TimeCount int    `json:"time_count"`
	TimeUnit  string `json:"time_unit"`
}

func GetDashBoardData(request RequestPayload, nodename, servicename string) (DashboardInfo, error) {
	var para DashboardInfo
	client, err := api.NewClient(api.Config{
		Address: "http://prometheus:9090/",
	})
	if err != nil {
		return DashboardInfo{}, err
	}
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	para.TimeUnit = request.TimeUnit
	timeStampArray := make([]int64, request.TimeCount)
	var timeRange v1.Range
	switch request.TimeUnit {
	case "second":
		for i := range timeStampArray {
			timeStampArray[i] = request.TimeStamp - int64((len(timeStampArray)-i)*5)
		}
		timeRange.Step = time.Second * 5
	case "minute":
		for i := range timeStampArray {
			timeStampArray[i] = request.TimeStamp - int64((len(timeStampArray)-i)*60)
		}
		timeRange.Step = time.Minute
	case "hour":
		for i := range timeStampArray {
			timeStampArray[i] = request.TimeStamp - int64((len(timeStampArray)-i)*3600)
		}
		timeRange.Step = time.Hour
	case "day":
		for i := range timeStampArray {
			timeStampArray[i] = request.TimeStamp - int64((len(timeStampArray)-i)*3600*24)
		}
		timeRange.Step = time.Hour * 24
	default:
		return DashboardInfo{}, errors.New("wrong time unit")
	}
	timeRange.Start = time.Unix(timeStampArray[0], 0)
	timeRange.End = time.Unix(timeStampArray[request.TimeCount-1], 0)

	// if timestamp < timeStampArray[0] {
	// 	para.IsOverMinLimit = true
	// } else if timestamp > timeStampArray[request.TimeCount-1] {
	// 	para.IsOverMaxLimit = true
	// }

	result, warnings, err := v1api.QueryRange(ctx, "kube_service_spec_selector{service!=\"kubernetes\"}", timeRange)
	if err != nil {
		return DashboardInfo{}, err
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	lines := grepString(result.String(), "{[^}]*}")
	para.ServiceListData = make([]ServiceList, len(lines)+1)
	para.ServiceListData[0].Name = "total"
	para.ServiceListData[0].ServiceLogsData = make([]ServiceLogs, request.TimeCount)
	for i, line := range lines {
		para.ServiceCount++
		serviceName := grepString(line, "service=[^,}]+")
		para.ServiceListData[i+1].Name = strings.Trim(grepString(serviceName[0], `"[^"]+"`)[0], "\"")
		para.ServiceListData[i+1].ServiceLogsData = make([]ServiceLogs, request.TimeCount)
		serviceSelectorLabels := grepString(line, "label_[^,}]+")
		if len(serviceSelectorLabels) == 0 {
			continue
		}
		grepName := "kube_pod_labels{" + sliceToPrometheusMetrics(serviceSelectorLabels) + "}"
		pods, _, err := v1api.QueryRange(ctx, grepName, timeRange)
		if err != nil {
			return DashboardInfo{}, err
		}
		podlines := strings.Split(pods.String(), "kube_pod_labels")
		for _, podline := range podlines[1:] {
			for s, time := range timeStampArray {
				if strings.Contains(podline, strconv.FormatInt(time, 10)) {
					para.ServiceListData[0].ServiceLogsData[s].TimeStamp = time
					para.ServiceListData[0].ServiceLogsData[s].PodNumber++
					para.ServiceListData[i+1].ServiceLogsData[s].TimeStamp = time
					para.ServiceListData[i+1].ServiceLogsData[s].PodNumber++
				}
			}
			podName := grepString(podline, "pod=[^,}]+")
			containerGrepName := "kube_pod_container_info{" + sliceToPrometheusMetrics(podName) + "}"
			containers, _, err := v1api.QueryRange(ctx, containerGrepName, timeRange)
			if err != nil {
				return DashboardInfo{}, err
			}
			containerlines := strings.Split(containers.String(), "kube_pod_container_info")
			for _, containerline := range containerlines[1:] {
				for s, time := range timeStampArray {
					if strings.Contains(containerline, strconv.FormatInt(time, 10)) {
						para.ServiceListData[0].ServiceLogsData[s].ContainerNumber++
						para.ServiceListData[i+1].ServiceLogsData[s].ContainerNumber++
					}
				}
			}
		}
	}

	for i := 0; i < len(para.ServiceListData); i++ {
		if para.ServiceListData[i].Name != servicename {
			para.ServiceListData[i].ServiceLogsData = []ServiceLogs{}
		}
	}

	//------------------------storage-total--------------------------------
	StorageCapQuery := `kube_node_status_capacity{resource="ephemeral_storage"}`
	StorageCapResult, _, err := v1api.QueryRange(ctx, StorageCapQuery, timeRange)
	if err != nil {
		return DashboardInfo{}, err
	}
	lineOfStorageCap := strings.Split(StorageCapResult.String(), "kube_node_status_capacity")
	para.NodeCount = len(lineOfStorageCap) - 1
	para.NodeListData = make([]NodeList, len(lineOfStorageCap))
	para.NodeListData[0].Name = "average"
	para.NodeListData[0].NodeLogsData = make([]NodeLogs, request.TimeCount)

	for i, v := range lineOfStorageCap[1:] {
		nodeName := grepString(v, "node=[^,}]+")
		para.NodeListData[i+1].Name = strings.Trim(grepString(nodeName[0], `"[^"]+"`)[0], "\"")
		para.NodeListData[i+1].NodeLogsData = make([]NodeLogs, request.TimeCount)
		data := grepString(v, "\n[0-9]+")
		for j, w := range data {
			para.NodeListData[i+1].NodeLogsData[j].TimeStamp = timeStampArray[j]
			digits, err := strconv.Atoi(strings.Trim(w, "\n"))
			if err != nil {
				return DashboardInfo{}, err
			}
			para.NodeListData[i+1].NodeLogsData[j].StorageTotal = digits
		}
	}
	//------------------------storage-used---------------------------------
	StorageUsedQuery := `kube_node_status_capacity{resource="ephemeral_storage"} - kube_node_status_allocatable{resource="ephemeral_storage"}`
	err = para.GetData(StorageUsedQuery, "storageUsed", v1api, ctx, timeRange, timeStampArray)
	if err != nil {
		return DashboardInfo{}, err
	}
	//------------------------memory-usage---------------------------------
	MemoUsage := `(1 - kube_node_status_allocatable_memory_bytes / kube_node_status_capacity_memory_bytes) * 100`
	err = para.GetData(MemoUsage, "memory", v1api, ctx, timeRange, timeStampArray)
	if err != nil {
		return DashboardInfo{}, err
	}
	//-------------------------CPU-usage-----------------------------------
	CPUUsage := `100 * (1 - sum by (instance)(node_cpu_seconds_total{mode="idle"}) / sum by (instance)(node_cpu_seconds_total))`
	err = para.GetData(CPUUsage, "CPU", v1api, ctx, timeRange, timeStampArray)
	if err != nil {
		return DashboardInfo{}, err
	}

	for i := 0; i < len(para.NodeListData); i++ {
		if para.NodeListData[i].Name != nodename {
			para.NodeListData[i].NodeLogsData = []NodeLogs{}
		}
	}

	//average:
	if nodename == "average" {
		CPUUsageAvg := fmt.Sprintf("avg(%s)", CPUUsage)
		err = para.GetAvgData(CPUUsageAvg, "CPU", v1api, ctx, timeRange, timeStampArray)
		if err != nil {
			return DashboardInfo{}, err
		}
		MemoUsageAvg := fmt.Sprintf("avg(%s)", MemoUsage)
		err = para.GetAvgData(MemoUsageAvg, "memory", v1api, ctx, timeRange, timeStampArray)
		if err != nil {
			return DashboardInfo{}, err
		}
		StorageUsedAvg := fmt.Sprintf("avg(%s)", StorageUsedQuery)
		err = para.GetAvgData(StorageUsedAvg, "storageUsed", v1api, ctx, timeRange, timeStampArray)
		if err != nil {
			return DashboardInfo{}, err
		}
		StorageCapAvg := fmt.Sprintf("avg(%s)", StorageCapQuery)
		err = para.GetAvgData(StorageCapAvg, "storageCap", v1api, ctx, timeRange, timeStampArray)
		if err != nil {
			return DashboardInfo{}, err
		}
	}

	return para, nil
}

func grepString(src, reg string) []string {
	re, _ := regexp.Compile(reg)
	return re.FindAllString(src, -1)
}

func sliceToPrometheusMetrics(src []string) string {
	return strings.Join(src, ", ")
}

func (d *DashboardInfo) GetAvgData(query, which string, v1api v1.API, ctx context.Context, timeRange v1.Range, timeStampArray []int64) error {
	result, _, err := v1api.QueryRange(ctx, query, timeRange)
	if err != nil {
		return err
	}
	data := grepString(result.String(), "\n[0-9.]+")
	for j, w := range data {
		var digitsInt int
		var digitsFloat float64
		switch which {
		case "storageUsed", "storageCap":
			digitsInt, err = strconv.Atoi(strings.Trim(w, "\n"))
		case "CPU", "memory":
			digitsFloat, err = strconv.ParseFloat(strings.Trim(w, "\n"), 64)
		}
		if err != nil {
			return err
		}
		switch which {
		case "CPU":
			d.NodeListData[0].NodeLogsData[j].TimeStamp = timeStampArray[j]
			d.NodeListData[0].NodeLogsData[j].CPUUsage = digitsFloat
		case "memory":
			d.NodeListData[0].NodeLogsData[j].MemoryUsage = digitsFloat
		case "storageCap":
			d.NodeListData[0].NodeLogsData[j].StorageTotal = digitsInt
		case "storageUsed":
			d.NodeListData[0].NodeLogsData[j].StorageUsed = digitsInt
		}
	}
	return nil
}

func (d *DashboardInfo) GetData(query, which string, v1api v1.API, ctx context.Context, timeRange v1.Range, timeStampArray []int64) error {
	result, _, err := v1api.QueryRange(ctx, query, timeRange)
	if err != nil {
		return err
	}
	lines := strings.Split(result.String(), "=>")
	for i, v := range lines[1:] {
		data := grepString(v, "\n[0-9.]+")
		for j, w := range data {
			var digitsInt int
			var digitsFloat float64
			switch which {
			case "storageUsed", "storageCap":
				digitsInt, err = strconv.Atoi(strings.Trim(w, "\n"))
			case "CPU", "memory":
				digitsFloat, err = strconv.ParseFloat(strings.Trim(w, "\n"), 64)
			}
			if err != nil {
				return err
			}
			switch which {
			case "CPU":
				d.NodeListData[i+1].NodeLogsData[j].CPUUsage = digitsFloat
			case "memory":
				d.NodeListData[i+1].NodeLogsData[j].MemoryUsage = digitsFloat
			case "storageCap":
				d.NodeListData[i+1].NodeLogsData[j].StorageTotal = digitsInt
			case "storageUsed":
				d.NodeListData[i+1].NodeLogsData[j].StorageUsed = digitsInt
			}
		}
	}
	return nil
}
