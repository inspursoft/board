package dashboard

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"regexp"
)

type Node struct {
	Id             int64    `json:"id" orm:"pk;auto"`
	NodeName       string `json:"pod_name" orm:"column(node_name)"`
	NumbersCpuCore string  `json:"pod_name" orm:"column(numbers_cpu_core)"`
	NumbersGpuCore string  `json:"pod_name" orm:"column(numbers_gpu_core)"`
	MemorySize     string  `json:"pod_name" orm:"column(memory_size)"`
	PodLimit       string  `json:"pod_name" orm:"column(pod_limit)"`
	CreateTime     string `json:"Creat_time" orm:"column(create_time)"`
	InternalIp     string `json:"ip" orm:"column(ip)"`
	CpuUsage       float32 `json:"cpu_usage" orm:"column(cpu_usage)"`
	MemUsage       float32 `json:"mem_usage" orm:"column(mem_usage)"`
	TimeListId     int64 `json:"pod_name" orm:"column(time_list_id)"`
	StorageTotal   int64 `json:"pod_name" orm:"column(storage_total)"`
	StorageUse     int64 `json:"pod_name" orm:"column(storage_use)"`
	TimeStamp      int64 `json:"pod_name" orm:"column(record_time)"`
}

func genTableStr(TimeUnit string) string {
	switch TimeUnit {
	case "second":
		return "node"
	case "minute":
		return "node_dashboard_minute"
	case "hour":
		return "node_dashboard_hour"
	case "day":
		return "node_dashboard_day"
	default:
		return "wrong"
	}
}
func verifyNodeName(nodeName string) bool {
	ma := `^(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|
	[1-9])\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|
	[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|
	[1-9]|0)\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|
	[0-9])$`
	if a, _ := regexp.MatchString(ma, nodeName); a == true {
		return true
	} else {
		return false
	}
}
func verifyInt(input string) bool {
	i := `^\+?[1-9][0-9]*$`
	if a, _ := regexp.MatchString(i, input); a == true {
		return true
	} else {
		return false
	}
}
func genSQLStr(nodeName string, recordTime string, count string, TimeUnit string) (sqlStr string) {
	if genTableStr(TimeUnit) != "wrong" && verifyNodeName(nodeName) && verifyInt(count) && verifyInt(recordTime) {
		sqlStr = fmt.Sprintf(`SELECT DISTINCT *
FROM %s
  RIGHT JOIN time_list_log ON node.time_list_id = time_list_log.id
WHERE 1 = 1
      AND node_name = '%s'
      AND record_time <= %s
ORDER BY time_list_id DESC
LIMIT  %s;`, genTableStr(TimeUnit), nodeName, recordTime, count)
	} else {
		fmt.Println("input wrong")
	}

	return
}
func nodeListSql() string {
	return `
SELECT DISTINCT
  DISTINCT
  node_name,
  time_list_id
FROM node
WHERE time_list_id = (SELECT max(time_list_id)
                      FROM node);`

}
func genNodeModel() (model []Node) {
	return
}
func QueryNode(nodeName string, recordTime string, count string, TimeUnit string) []Node {
	ss := genNodeModel()
	queryBySql(genSQLStr(nodeName, recordTime, count, TimeUnit))(&ss)
	return ss

}
func NodeList() {
	mode := nodeListModel()
	queryBySql(nodeListSql())(&mode)
	fmt.Println(mode)

}
func queryBySql(sqlStr string) func(... interface{}) (int64, error) {
	return orm.NewOrm().Raw(sqlStr).QueryRows
}
func nodeListModel() (list []struct {
	NodeName   string `json:"pod_name" orm:"column(node_name)"`
	TimeListId int64 `json:"pod_name" orm:"column(time_list_id)"`
}) {
	return
}
