package dao

import (
	"git/inspursoft/board/src/collector/model/collect"
	"git/inspursoft/board/src/collector/model/collect/dashboard"

	"github.com/astaxie/beego/orm"
)

func InsertDb(model interface{}) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(model)
	return id, err
}

func QuerDb(models interface{}, TableName string, filter bool, filter_tag string,
	filter_value interface{}, selectedFields ...string) (err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(TableName)
	if filter == false {
		_, err = qs.All(models, selectedFields...)
	} else {
		_, err = qs.Filter(filter_tag, filter_value).
			All(models, selectedFields...)
	}
	if err != nil {
		return err
	}
	return nil
}

func QuerDbMax(model interface{}) error {
	var sql string
	sql = "select * from time_list_log where id=(select max(id) from time_list_log)"
	o := orm.NewOrm()
	err := o.Raw(sql).QueryRow(model)
	return err
}

func nRow(timeUnit string) string {
	switch timeUnit {
	case "minute":
		return "12"
	case "hour":
		return "60"
	case "day":
		return "24"
	}
	return "none"
}
func QuerTimeListID(timeUnit string) (TimeList []struct {
	Temp int64 `json:"pod_name" orm:"column(time_list_id)"`
}) {
	sql := "SELECT DISTINCT time_list_id FROM node ORDER BY time_list_id DESC LIMIT " + nRow(timeUnit)
	temp := TimeList
	o := orm.NewOrm()
	o.Raw(sql).QueryRows(&temp)
	return temp
}
func querTable(timeUnit string) string {
	switch timeUnit {
	case "minute":
		return "node"
	case "hour":
		return "node_dashboard_minute"
	case "day":
		return "node_dashboard_hour"
	}
	return "none"
}
func nodesSlice(timeUnit string) (res interface{}) {
	switch timeUnit {
	case "minute":
		res = []collect.Node{}
		temp := res.([]collect.Node)
		o := orm.NewOrm()
		qs := o.QueryTable(querTable(timeUnit))
		qs = qs.Filter("time_list_id__in", func() []int64 {
			var i []int64
			for _, v := range QuerTimeListID(timeUnit) {
				i = append(i, v.Temp)
			}
			return i
		}())
		qs.All(&temp)
		return temp

	case "hour":
		res = []dashboard.NodeDashboardMinute{}
		temp := res.([]dashboard.NodeDashboardMinute)
		o := orm.NewOrm()
		qs := o.QueryTable(querTable(timeUnit))
		qs = qs.Filter("time_list_id__in", func() []int64 {
			var i []int64
			for _, v := range QuerTimeListID(timeUnit) {
				i = append(i, v.Temp)
			}
			return i
		}())
		qs.All(&temp)
		return temp
	case "day":
		res = []dashboard.NodeDashboardHour{}
		temp := res.([]dashboard.NodeDashboardHour)
		o := orm.NewOrm()
		qs := o.QueryTable(querTable(timeUnit))
		qs = qs.Filter("time_list_id__in", func() []int64 {
			var i []int64
			for _, v := range QuerTimeListID(timeUnit) {
				i = append(i, v.Temp)
			}
			return i
		}())
		qs.All(&temp)
		return temp
	}

	return nil
}
func assignNode(ori collect.Node, tar collect.Node) collect.Node {
	return collect.Node{
		NodeName:       ori.NodeName,
		CpuUsage:       (ori.CpuUsage + tar.CpuUsage) / 2.0,
		MemUsage:       (ori.MemUsage + tar.MemUsage) / 2.0,
		StorageTotal:   (ori.StorageTotal + tar.StorageTotal) / 2,
		StorageUse:     (ori.StorageUse + tar.StorageUse) / 2,
		NumbersGpuCore: ori.NumbersGpuCore,
		NumbersCpuCore: ori.NumbersCpuCore,
		TimeListId:     ori.TimeListId,
		PodLimit:       ori.PodLimit,
		CreateTime:     ori.CreateTime,
		InternalIp:     ori.InternalIp,
		MemorySize:     ori.MemorySize,
	}
}

func CalcNode(timeUnit string) interface{} {

	switch timeUnit {
	case "minute":
		var res []collect.Node
		for _, v := range nodesSlice(timeUnit).([]collect.Node) {
			if len(res) == 0 {
				res = append(res, v)
			}
			if k := 0; any(res, v.NodeName, &k, timeUnit) {
				res[k] = assignNode(res[k], v)
			} else {
				res = append(res, v)
			}
		}
		return res
	case "hour":
		var res []dashboard.NodeDashboardMinute
		for _, v := range nodesSlice(timeUnit).([]dashboard.NodeDashboardMinute) {
			if len(res) == 0 {
				res = append(res, v)
			}
			if k := 0; any(res, v.NodeName, &k, timeUnit) {
				res[k] = dashboard.NodeDashboardMinute(assignNode(collect.Node(res[k]), collect.Node(v)))
			} else {
				res = append(res, v)
			}
		}
		return res
	case "day":
		var res []dashboard.NodeDashboardHour
		for _, v := range nodesSlice(timeUnit).([]dashboard.NodeDashboardHour) {
			if len(res) == 0 {
				res = append(res, v)
			}
			if k := 0; any(res, v.NodeName, &k, timeUnit) {
				res[k] = dashboard.NodeDashboardHour(assignNode(collect.Node(res[k]), collect.Node(v)))
			} else {
				res = append(res, v)
			}
		}
		return res
	}

	return nil
}

func any(res interface{}, string string, k *int, timeUnit string) bool {
	switch timeUnit {
	case "minute":
		res := res.([]collect.Node)
		for i, v := range res {
			if v.NodeName == string {
				*k = i
				return true
			}
		}
		return false
	case "hour":
		res := res.([]dashboard.NodeDashboardMinute)
		for i, v := range res {
			if v.NodeName == string {
				*k = i
				return true
			}
		}
		return false

	case "day":
		res := res.([]dashboard.NodeDashboardHour)
		for i, v := range res {
			if v.NodeName == string {
				*k = i
				return true
			}
		}
		return false
	}
	return false
}

func DeleteStaleData(tableName string, dateSpan int) (int64, error) {
	o := orm.NewOrm()
	ptmt, err := o.Raw("delete t from " + tableName + " t " +
		" where t.time_list_id in (select id as time_list_id from time_list_log " +
		"     where datediff(now(), date_format(from_unixtime(record_time), '%Y-%m-%d %H:%i:%s')) >= ?);").Prepare()
	if err != nil {
		return 0, err
	}
	rs, err := ptmt.Exec(dateSpan)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return rs.RowsAffected()
}

func DeleteStaleTimeList(dateSpan int) (int64, error) {
	o := orm.NewOrm()
	ptmt, err := o.Raw(`delete t from time_list_log t
		where datediff(now(), date_format(from_unixtime(t.record_time), '%Y-%m-%d %H:%i:%s')) >=?;`).Prepare()
	if err != nil {
		return 0, err
	}
	rs, err := ptmt.Exec(dateSpan)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return rs.RowsAffected()
}
