package dashboard

import (
	"fmt"
	model "git/inspursoft/board/src/common/model/dashboard"

	"github.com/astaxie/beego/orm"
)

func timeListFilter(table string, filter interface{}) orm.QuerySeter {
	return orm.NewOrm().QueryTable(table).Filter("time_list_id__in", filter)
}
func PreQueryTotal(sql *string, sqlChan chan string, count chan string,
	timestampChan chan string, doneSql chan bool) {
	tableName := <-sqlChan
	i := <-count
	timeStamp := <-timestampChan
	*sql = fmt.Sprintf(`
SELECT
	DISTINCT time_list_id
FROM
	%s
WHERE
	time_list_id <= (
		SELECT
			time_list_log.id FROM time_list_log
		WHERE
			record_time <= %s
		ORDER BY
			id DESC LIMIT 1   )
ORDER BY
	    	time_list_id DESC
	    LIMIT %s ;`, tableName, timeStamp, i)
	doneSql <- true
}
func QueryTime(sql string) (timeId []int64) {
	type TimeList struct {
		TimeId int64 `json:"pod_name" orm:"column(time_list_id)"`
	}
	var temp []TimeList
	o := orm.NewOrm()
	o.Raw(sql).QueryRows(&temp)
	for _, v := range temp {
		timeId = append(timeId, v.TimeId)
	}
	return
}

func QueryTotal(timeUnit string, sqlChan chan string, timeIdChan chan []int64,
	myServerList *interface{}, doneS chan bool) {
	switch timeUnit {
	case "hour":
		go func() { sqlChan <- "service_dashboard_hour" }()
		go func(myServerList *interface{}) {
			var s []model.ServiceDashboardHour
			i := <-timeIdChan
			qs := timeListFilter("service_dashboard_hour", i).Limit(10000)
			qs.All(&s)
			*myServerList = &s
			doneS <- true
		}(myServerList)
	case "day":
		go func() { sqlChan <- "service_dashboard_day" }()
		go func(myServerList *interface{}) {
			var s []model.ServiceDashboardDay
			i := <-timeIdChan
			qs := timeListFilter("service_dashboard_day", i).Limit(10000)
			qs.All(&s)
			*myServerList = &s
			doneS <- true
		}(myServerList)
	case "minute":
		go func() { sqlChan <- "service_dashboard_minute" }()
		go func(myServerList *interface{}) {
			var s []model.ServiceDashboardMinute
			i := <-timeIdChan
			qs := timeListFilter("service_dashboard_minute", i).Limit(10000)
			qs.All(&s)
			*myServerList = &s
			doneS <- true
		}(myServerList)
	case "second":
		go func() { sqlChan <- "service_dashboard_second" }()
		go func(myServerList *interface{}) {
			var s []model.ServiceDashboardSecond
			i := <-timeIdChan
			qs := timeListFilter("service_dashboard_second", i).Limit(10000)
			qs.All(&s)
			*myServerList = &s
			doneS <- true

		}(myServerList)

	}
}
func GetDashboardServiceList() (service []model.ServiceDashboardSecond) {
	var timeID model.TimeListLog
	QuerDbMax(&timeID, "time_list_log", "id")
	maxTimeId := timeID.Id
	QuerySet(&service, "service_dashboard_second", true, "time_list_id",
		maxTimeId)
	return
}

func QuerDbMax(model interface{}, table string, maxInt string) error {
	var sql string
	sql = "select * from " + table + " where id=(select max(" + maxInt + ") from " + table + ")"
	o := orm.NewOrm()
	err := o.Raw(sql).QueryRow(model)
	return err
}
func QuerySet(models interface{}, TableName string, filter bool, filter_tag string,
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
