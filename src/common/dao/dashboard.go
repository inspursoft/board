package dao

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	defaultMaxQueryNum = 500
)

type DashboardNodeDao struct {
	QueryPara
}
type DashboardServiceDao struct {
	QueryPara
}
type QueryPara struct {
	TimeUnit  string
	TimeStamp int
	TimeCount int
	Name      string
	DuraTime  int
}
type ServiceDataLog struct {
	PodNumber       int64 `json:"pod_number" orm:"column(pod_number)"`
	ContainerNumber int64 `json:"container_number" orm:"column(container_number)"`
	Timestamp       int   `json:"timestamp" orm:"column(record_time)"`
}
type NodeDataLogs struct {
	Record_time   int     `json:"timestamp" orm:"column(record_time)"`
	Cpu_usage     float32 `json:"cpu_usage" orm:"column(cpu_usage)"`
	Mem_usage     float32 `json:"memory_usage" orm:"column(mem_usage)"`
	Storage_total float32 `json:"storage_total" orm:"column(storage_total)"`
	Storage_use   float32 `json:"storage_use" orm:"column(storage_use)"`
}
type NodeListDataLogs struct {
	NodeName  string `json:"node_name" orm:"column(node_name)"`
	Timestamp int64  `json:"timestamp" orm:"column(record_time)"`
}
type ServiceListDataLogs struct {
	NodeName  string `json:"service_name" orm:"column(service_name)"`
	Timestamp int64  `json:"timestamp" orm:"column(record_time)"`
}
type LimitTime struct {
	FieldName string `orm:"column(field_name)"`
	MinTime   int    `orm:"column(min_time)"`
	MaxTime   int    `orm:"column(max_time)"`
}

func (d *DashboardServiceDao) GetLimitTime() (LimitTime, error) {
	var lt LimitTime
	tableName, err := d.getServiceDataTableName()
	if err != nil {
		return lt, err
	}
	if d.Name != "" {
		sql := fmt.Sprintf(`
SELECT
  sd.service_name      AS field_name,
  min(tll.record_time) AS min_time,
  max(tll.record_time) AS max_time
FROM %s sd
  JOIN time_list_log tll ON time_list_id = tll.id
WHERE 1 = 1
      AND sd.service_name = ?
GROUP BY sd.service_name;
	`, tableName)
		o := orm.NewOrm()
		err = o.Raw(sql, d.Name).QueryRow(&lt)
		return lt, err
	} else {
		sql := fmt.Sprintf(`
SELECT
  min(tll.record_time) AS min_time,
  max(tll.record_time) AS max_time
FROM %s sd
  JOIN time_list_log tll ON time_list_id = tll.id
WHERE 1 = 1;`, tableName)
		o := orm.NewOrm()
		err = o.Raw(sql).QueryRow(&lt)
		return lt, err
	}

}

func (d *DashboardNodeDao) GetLimitTime() (LimitTime, error) {
	var lt LimitTime
	tableName, err := d.getNodeDataTableName()
	if err != nil {
		return lt, err
	}
	if d.Name != "" {
		sql := fmt.Sprintf(`
SELECT
  sd.node_name     AS field_name,
  min(tll.record_time) AS min_time,
  max(tll.record_time) AS max_time
FROM %s sd
  JOIN time_list_log tll ON time_list_id = tll.id
WHERE 1 = 1
      AND sd.node_name = ?
GROUP BY sd.node_name;;
	`, tableName)
		o := orm.NewOrm()
		err = o.Raw(sql, d.Name).QueryRow(&lt)
		return lt, err
	} else {
		sql := fmt.Sprintf(`
SELECT
  min(tll.record_time) AS min_time,
  max(tll.record_time) AS max_time
FROM %s sd
  JOIN time_list_log tll ON time_list_id = tll.id
WHERE 1 = 1 ;
`, tableName)
		o := orm.NewOrm()
		err = o.Raw(sql).QueryRow(&lt)
		return lt, err
	}

}

func (d *DashboardNodeDao) getNodeDataTableName() (string, error) {
	switch d.TimeUnit {
	case "second":
		return "node", nil
	case "minute":
		return "node_dashboard_minute", nil
	case "hour":
		return "node_dashboard_hour", nil
	case "day":
		return "node_dashboard_day", nil
	}
	return "", errors.New("wrong")
}

func (d *DashboardNodeDao) getDurationTime() (last int, prev int, err error) {
	if d.TimeStamp == 0 {
		return 0, 0, errors.New("no time stamp")
	}
	if d.DuraTime == 0 {
		switch d.QueryPara.TimeUnit {
		case "second":
			t := d.TimeCount * 5
			return d.TimeStamp, d.TimeStamp - t, nil
		case "minute":
			t := d.TimeCount * 60
			return d.TimeStamp, d.TimeStamp - t, nil
		case "hour":
			t := d.TimeCount * 60 * 60
			return d.TimeStamp, d.TimeStamp - t, nil
		case "day":
			t := d.TimeCount * 60 * 60 * 24
			return d.TimeStamp, d.TimeStamp - t, nil

		}
	} else {
		beego.Debug("given DuraTime", d.TimeStamp, d.TimeStamp-d.DuraTime)
		return d.TimeStamp, d.TimeStamp - d.DuraTime, nil
	}
	return
}
func (d *DashboardNodeDao) GetTotalNodeData() (count int, nodeItems []NodeDataLogs, err error) {
	if d.TimeCount > defaultMaxQueryNum {
		return count, []NodeDataLogs{}, errors.New("time count must < defaultMaxQueryNum")
	}

	tableName, err := d.getNodeDataTableName()
	if err != nil {
		return count, []NodeDataLogs{}, err
	}
	last, prev, err := d.getDurationTime()
	sql := d.genNodeDataTotalSqlString(tableName)
	o := orm.NewOrm()
	i, err := o.Raw(sql, prev, last, d.TimeCount).QueryRows(&nodeItems)
	count = int(i)
	return count, nodeItems, err

}
func (d *DashboardNodeDao) genNodeDataTotalSqlString(tableName string) string {
	sql := fmt.Sprintf(`
		(SELECT
	  nt.time_list_id,
	  nt.record_time,
	  avg(nt.cpu_usage)     AS cpu_usage,
	  avg(nt.mem_usage)     AS mem_usage,
	  avg(nt.memory_size)   AS memory_size,
	  avg(nt.storage_total) AS storage_total,
	  avg(nt.storage_use)   AS storage_use
	FROM (SELECT
	        n.memory_size   AS memory_size,
	        n.storage_total AS storage_total,
	        n.storage_use   AS storage_use,
	        n.cpu_usage     AS cpu_usage,
	        n.mem_usage     AS mem_usage,
	        t.record_time,
	        n.time_list_id  AS time_list_id
	      FROM  %s n
	        LEFT JOIN time_list_log t ON n.time_list_id = t.id
	      WHERE t.record_time >= ?
	            AND t.record_time <= ?
	      ORDER BY n.time_list_id DESC) AS nt
	GROUP BY nt.time_list_id, nt.record_time
	ORDER BY nt.time_list_id DESC
	LIMIT ?)
 	ORDER BY record_time ASC;
		`, tableName)
	return sql
}
func (d *DashboardNodeDao) getNodeListDataSqlString(tableName string) (sql string) {
	sql = fmt.Sprintf(`
	(SELECT distinct
   nd.node_name      AS node_name
 FROM %s nd
   LEFT JOIN time_list_log tll ON nd.time_list_id = tll.id
 WHERE tll.record_time >= ?
       AND tll.record_time <= ?)
 ORDER BY record_time ASC;
	`, tableName)

	return sql
}

func (d *DashboardServiceDao) getServiceListDataSqlString(tableName string) (sql string) {
	sql = fmt.Sprintf(`
(SELECT distinct
   sd.service_name      AS service_name
 FROM %s sd
   LEFT JOIN time_list_log tll ON sd.time_list_id = tll.id
 WHERE tll.record_time >= ?
       AND tll.record_time <= ?)
 ORDER BY record_time ASC;`, tableName)

	return sql
}

func (d *DashboardNodeDao) GetNodeListData() (count int, nodelistItems []NodeListDataLogs, err error) {
	tableName, err := d.getNodeDataTableName()
	if err != nil {
		return count, []NodeListDataLogs{}, err
	}
	last, prev, err := d.getDurationTime()
	sql := d.getNodeListDataSqlString(tableName)
	o := orm.NewOrm()
	i, err := o.Raw(sql, prev, last).QueryRows(&nodelistItems)
	count = int(i)
	fmt.Println(count)
	return count, nodelistItems, err
}

func (d *DashboardServiceDao) GetServiceListData() (count int, serviceListItems []ServiceListDataLogs, err error) {
	tableName, err := d.getServiceDataTableName()
	if err != nil {
		return count, []ServiceListDataLogs{}, err
	}
	last, prev, err := d.getDurationTime()
	sql := d.getServiceListDataSqlString(tableName)
	o := orm.NewOrm()
	i, err := o.Raw(sql, prev, last).QueryRows(&serviceListItems)
	count = int(i)
	fmt.Println(count)
	return count, serviceListItems, err
}

func (d *DashboardNodeDao) genNodeDataSqlString(tableName string) (sql string) {
	sql = fmt.Sprintf(`
	(SELECT
  nd.node_name     AS node_name,
  nd.time_list_id  AS time_list_id,
  tll.record_time  AS record_time,
  nd.cpu_usage     AS cpu_usage,
  nd.mem_usage     AS mem_usage,
  nd.memory_size   AS memory_size,
  nd.storage_total AS storage_total,
  nd.storage_use   AS storage_use
FROM %s nd
  LEFT JOIN time_list_log tll ON nd.time_list_id = tll.id
WHERE tll.record_time >= ?
      AND tll.record_time <= ?
      AND nd.node_name = ?
ORDER BY tll.record_time DESC
LIMIT ?)
 ORDER BY record_time ASC;`, tableName)
	return sql
}

func (d *DashboardNodeDao) GetNodeData() (count int, nodeItems []NodeDataLogs, err error) {
	if d.TimeCount > defaultMaxQueryNum {
		return count, []NodeDataLogs{}, errors.New("time count must < defaultMaxQueryNum")
	}
	tableName, err := d.getNodeDataTableName()
	if err != nil {
		return count, []NodeDataLogs{}, err
	}
	last, prev, err := d.getDurationTime()
	sql := d.genNodeDataSqlString(tableName)
	o := orm.NewOrm()
	i, err := o.Raw(sql, prev, last, d.Name, d.TimeCount).QueryRows(&nodeItems)
	count = int(i)

	return count, nodeItems, err

}
func (d *DashboardServiceDao) genServiceDataTotalSqlString(tableName string) string {
	sql := fmt.Sprintf(`
(SELECT
  sum(nt.pod_number)       AS pod_number,
  sum(nt.container_number) AS container_number,
  nt.record_time
FROM (SELECT
        n.time_list_id     AS time_list_id,
        n.pod_number       AS pod_number,
        n.container_number AS container_number,
        t.record_time
      FROM %s n
        LEFT JOIN time_list_log t ON n.time_list_id = t.id
      WHERE t.record_time >= ?
      AND t.record_time <= ? AND n.service_name in (select name from service_status)
      ORDER BY n.time_list_id DESC) AS nt
GROUP BY nt.record_time
ORDER BY nt.record_time DESC
LIMIT ?)
 ORDER BY record_time ASC;

		`, tableName)
	return sql
}

func (d *DashboardServiceDao) getServiceDataTableName() (string, error) {
	switch d.TimeUnit {
	case "second":
		return "service_dashboard_second", nil
	case "minute":
		return "service_dashboard_minute", nil
	case "hour":
		return "service_dashboard_hour", nil
	case "day":
		return "service_dashboard_day", nil
	}
	return "", errors.New("wrong")
}

func (d *DashboardServiceDao) getDurationTime() (last int, prev int, err error) {
	if d.TimeStamp == 0 {
		return 0, 0, errors.New("no time stamp")
	}
	if d.DuraTime == 0 {
		switch d.QueryPara.TimeUnit {
		case "second":
			t := d.TimeCount * 5
			return d.TimeStamp, d.TimeStamp - t, nil
		case "minute":
			t := d.TimeCount * 60
			return d.TimeStamp, d.TimeStamp - t, nil
		case "hour":
			t := d.TimeCount * 60 * 60
			return d.TimeStamp, d.TimeStamp - t, nil
		case "day":
			t := d.TimeCount * 60 * 60 * 24
			return d.TimeStamp, d.TimeStamp - t, nil

		}
	} else {
		fmt.Println("given DuraTime", d.TimeStamp, d.TimeStamp-d.DuraTime)
		return d.TimeStamp, d.TimeStamp - d.DuraTime, nil
	}
	return
}
func (d *DashboardServiceDao) GetTotalServiceData() (count int, serviceItems []ServiceDataLog, err error) {
	tableName, err := d.getServiceDataTableName()
	if err != nil {
		return count, []ServiceDataLog{}, err
	}
	last, prev, err := d.getDurationTime()
	sql := d.genServiceDataTotalSqlString(tableName)
	o := orm.NewOrm()
	i, err := o.Raw(sql, prev, last, d.TimeCount).QueryRows(&serviceItems)
	count = int(i)
	if err != nil {
		return count, []ServiceDataLog{}, err
	}

	if d.TimeCount > 500 {
		err = errors.New("time count must < 500")
	}
	return count, serviceItems, err

}
func (d *DashboardServiceDao) genServiceDataSqlString(tableName string) (sql string) {
	sql = fmt.Sprintf(`
(SELECT
   sd.service_name     AS service_name,
   sd.pod_number       AS pod_number,
   sd.container_number AS container_number,
   tll.record_time     AS record_time
 FROM %s sd
   JOIN time_list_log tll ON time_list_id = tll.id
 WHERE tll.record_time >= ?
       AND tll.record_time <= ?
       AND sd.service_name = ?
 ORDER BY tll.record_time DESC
 LIMIT ?)
 ORDER BY record_time ASC;
	`, tableName)
	return sql
}

/*func (d *DashboardServiceDao) newGenServiceDataSqlString() (sql string) {
	sql = fmt.Sprintf(`
(SELECT
   ds.service_name     AS service_name,
   ds.pod_number       AS pod_number,
   ds.container_number AS container_number,
   ds.time_list_id     AS record_time
 FROM %s ds
 WHERE ds.time_list_id >=  ? AND
       ds.time_list_id <= ? AND
       ds.service_name = ?
 ORDER BY ds.time_list_id DESC
 LIMIT ?)
 ORDER BY record_time ASC;
	`, fmt.Sprintf(`dashboard_service_%s`, time.Now().Format("2006_01_02")))
	return sql
}*/

func (d *DashboardServiceDao) GetServiceData() (count int, serviceItems []ServiceDataLog, err error) {
	tableName, err := d.getServiceDataTableName()
	if err != nil {
		return count, []ServiceDataLog{}, err
	}
	last, prev, err := d.getDurationTime()
	sql := d.genServiceDataSqlString(tableName)
	o := orm.NewOrm()
	i, err := o.Raw(sql, prev, last, d.Name, d.TimeCount).QueryRows(&serviceItems)

	count = int(i)
	if err != nil {
		return count, []ServiceDataLog{}, err
	}

	if d.TimeCount > 500 {
		err = errors.New("time count must < 500")
	}
	return count, serviceItems, err

}

/*func (d *DashboardServiceDao) GetServiceData() (count int, serviceItems []ServiceDataLog, err error) {
	tableName, err := d.getServiceDataTableName()
	if err != nil {
		return count, []ServiceDataLog{}, err
	}
	last, prev, err := d.getDurationTime()
	sql := d.genServiceDataSqlString(tableName)
	o := orm.NewOrm()
	i, err := o.Raw(sql, prev, last, d.Name, d.TimeCount).QueryRows(&serviceItems)

	count = int(i)
	if err != nil {
		return count, []ServiceDataLog{}, err
	}

	if d.TimeCount > 500 {
		err = errors.New("time count must < 500")
	}
	return count, serviceItems, err

}*/
