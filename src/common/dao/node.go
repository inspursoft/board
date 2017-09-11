package dao

import "github.com/astaxie/beego/orm"

type Node struct {
	NodeName     string `json:"node_name" orm:"column(node_name)"`
	NodeIP       string `json:"node_ip" orm:"column(node_ip)"`
	CreateTime   string `json:"create_time" orm:"column(create_time)"`
	CpuUsage     string `json:"cpu_usage" orm:"column(cpu_usage)"`
	MemoryUsage  string `json:"memory_usage" orm:"column(memory_usage)"`
	MemorySize   string `json:"memory_size" orm:"column(memory_size)"`
	StorageTotal string `json:"storage_total" orm:"column(storage_total)"`
	StorageUse   string `json:"storage_use" orm:"column(storage_usage)"`
}

func GetNode(nodeName string) (node Node, err error) {
	sql := `
	SELECT
  node.node_name AS node_name,
  node.ip        AS node_ip,
  create_time    AS create_time,
  cpu_usage      AS cpu_usage,
  mem_usage      AS memory_usage,
  memory_size    AS memory_size,
  storage_total  AS storage_total,
  storage_use    AS storage_usage
FROM node
WHERE time_list_id = (SELECT max(time_list_id)
                      FROM node)
      AND node_name = ?;
	`
	o := orm.NewOrm()
	_ = o.Raw(sql, nodeName).QueryRow(&node)
	return
}
