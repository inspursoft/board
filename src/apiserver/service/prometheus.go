package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
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

	podInfo, err := para.GetServiceInfo(ctx, v1api, timeRange, timeStampArray)
	if err != nil {
		return DashboardInfo{}, err
	}

	containerInfo, err := para.CountPod(ctx, v1api, timeRange, timeStampArray, podInfo)
	if err != nil {
		return DashboardInfo{}, err
	}

	para.CountContainer(timeStampArray, containerInfo)

	for i := 0; i < len(para.ServiceListData); i++ {
		if para.ServiceListData[i].Name != servicename {
			para.ServiceListData[i].ServiceLogsData = []ServiceLogs{}
		}
	}

	storageCapQuery := `kube_node_status_capacity{resource="ephemeral_storage"}`
	storageCapResult, _, err := v1api.QueryRange(ctx, storageCapQuery, timeRange)
	if err != nil {
		return DashboardInfo{}, err
	}
	lineOfStorageCap := strings.Split(storageCapResult.String(), "kube_node_status_capacity")
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

	storageUsedQuery := `kube_node_status_capacity{resource="ephemeral_storage"} - kube_node_status_allocatable{resource="ephemeral_storage"}`
	err = para.GetData(storageUsedQuery, "storageUsed", v1api, ctx, timeRange)
	if err != nil {
		return DashboardInfo{}, err
	}

	memoUsageQuery := `(1 - kube_node_status_allocatable_memory_bytes / kube_node_status_capacity_memory_bytes) * 100`
	err = para.GetData(memoUsageQuery, "memory", v1api, ctx, timeRange)
	if err != nil {
		return DashboardInfo{}, err
	}

	cpuUsageQuery := `100 * (1 - sum by (instance)(node_cpu_seconds_total{mode="idle"}) / sum by (instance)(node_cpu_seconds_total))`
	err = para.GetData(cpuUsageQuery, "CPU", v1api, ctx, timeRange)
	if err != nil {
		return DashboardInfo{}, err
	}

	for i := 0; i < len(para.NodeListData); i++ {
		if para.NodeListData[i].Name != nodename {
			para.NodeListData[i].NodeLogsData = []NodeLogs{}
		}
	}

	if nodename == "average" {
		cpuUsageAvg := fmt.Sprintf("avg(%s)", cpuUsageQuery)
		err = para.GetAvgData(cpuUsageAvg, "CPU", v1api, ctx, timeRange, timeStampArray)
		if err != nil {
			return DashboardInfo{}, err
		}
		memoUsageAvg := fmt.Sprintf("avg(%s)", memoUsageQuery)
		err = para.GetAvgData(memoUsageAvg, "memory", v1api, ctx, timeRange, timeStampArray)
		if err != nil {
			return DashboardInfo{}, err
		}
		storageUsedAvg := fmt.Sprintf("avg(%s)", storageUsedQuery)
		err = para.GetAvgData(storageUsedAvg, "storageUsed", v1api, ctx, timeRange, timeStampArray)
		if err != nil {
			return DashboardInfo{}, err
		}
		storageCapAvg := fmt.Sprintf("avg(%s)", storageCapQuery)
		err = para.GetAvgData(storageCapAvg, "storageCap", v1api, ctx, timeRange, timeStampArray)
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

func (d *DashboardInfo) GetData(query, which string, v1api v1.API, ctx context.Context, timeRange v1.Range) error {
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
			case "storageUsed":
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
			case "storageUsed":
				d.NodeListData[i+1].NodeLogsData[j].StorageUsed = digitsInt
			}
		}
	}
	return nil
}

func (d *DashboardInfo) GetServiceInfo(ctx context.Context, v1api v1.API, timeRange v1.Range, timeStampArray []int64) ([][]string, error) {
	result, warnings, err := v1api.QueryRange(ctx, "kube_service_spec_selector{service!=\"kubernetes\"}", timeRange)
	if err != nil {
		return [][]string{}, err
	}
	if len(warnings) > 0 {
		logs.Info("Warnings: %v\n", warnings)
	}
	lines := grepString(result.String(), "{[^}]*}")
	d.ServiceListData = make([]ServiceList, len(lines)+1)
	d.ServiceListData[0].Name = "total"
	d.ServiceListData[0].ServiceLogsData = make([]ServiceLogs, len(timeStampArray))
	for s, time := range timeStampArray {
		d.ServiceListData[0].ServiceLogsData[s].TimeStamp = time
	}
	podResults := make([][]string, len(lines))
	for i, line := range lines {
		d.ServiceCount++
		serviceName := grepString(line, "service=[^,}]+")
		d.ServiceListData[i+1].Name = strings.Trim(grepString(serviceName[0], `"[^"]+"`)[0], "\"")
		serviceSelectorLabels := grepString(line, "label_[^,}]+")
		if len(serviceSelectorLabels) == 0 {
			continue
		}
		d.ServiceListData[i+1].ServiceLogsData = make([]ServiceLogs, len(timeStampArray))
		for s, time := range timeStampArray {
			d.ServiceListData[i+1].ServiceLogsData[s].TimeStamp = time
		}
		grepName := "kube_pod_labels{" + sliceToPrometheusMetrics(serviceSelectorLabels) + "}"
		pods, _, err := v1api.QueryRange(ctx, grepName, timeRange)
		if err != nil {
			return [][]string{}, err
		}
		podlines := strings.Split(pods.String(), "kube_pod_labels")
		podResults[i] = podlines[1:]
	}
	return podResults, nil
}

func (d *DashboardInfo) CountPod(ctx context.Context, v1api v1.API, timeRange v1.Range, timeStampArray []int64, podResults [][]string) ([][]string, error) {
	containerResults := make([][]string, len(podResults))
	for i, svc := range podResults {
		if len(svc) == 0 {
			continue
		}
		for _, podline := range svc {
			for s, time := range timeStampArray {
				if strings.Contains(podline, strconv.FormatInt(time, 10)) {
					d.ServiceListData[0].ServiceLogsData[s].PodNumber++
					d.ServiceListData[i+1].ServiceLogsData[s].PodNumber++
				}
			}
			podName := grepString(podline, "pod=[^,}]+")
			containerGrepName := "kube_pod_container_info{" + sliceToPrometheusMetrics(podName) + "}"
			containers, _, err := v1api.QueryRange(ctx, containerGrepName, timeRange)
			if err != nil {
				return [][]string{}, err
			}
			containerlines := strings.Split(containers.String(), "kube_pod_container_info")
			containerResults[i] = append(containerResults[i], containerlines[1:]...)
		}
	}
	return containerResults, nil
}

func (d *DashboardInfo) CountContainer(timeStampArray []int64, containerResults [][]string) {
	for i, svc := range containerResults {
		if len(svc) == 0 {
			continue
		}
		for _, containerline := range svc {
			for s, time := range timeStampArray {
				if strings.Contains(containerline, strconv.FormatInt(time, 10)) {
					d.ServiceListData[0].ServiceLogsData[s].ContainerNumber++
					d.ServiceListData[i+1].ServiceLogsData[s].ContainerNumber++
				}
			}
		}
	}
}
