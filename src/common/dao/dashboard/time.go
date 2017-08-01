package dashboard

import (
	"strconv"
	"github.com/astaxie/beego/orm"
)

func QueryTimeList(model interface{}, id int64) error {
	var sql string
	sId := strconv.Itoa(int(id))
	sql = "select * from time_list_log WHERE id = " + sId
	o := orm.NewOrm()
	err := o.Raw(sql).QueryRow(model)
	return err
}

