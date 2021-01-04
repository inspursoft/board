package dao

import (
	"github.com/inspursoft/board/src/common/model"

	//"time"

	//"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func AddAutoScale(autoscale model.ServiceAutoScale) (int64, error) {
	o := orm.NewOrm()

	autoscaleID, err := o.Insert(&autoscale)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return autoscaleID, err
}

func GetAutoScale(autoscale model.ServiceAutoScale, fieldNames ...string) (*model.ServiceAutoScale, error) {
	o := orm.NewOrm()
	err := o.Read(&autoscale, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &autoscale, err
}

func UpdateAutoScale(autoscale model.ServiceAutoScale, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	//autoscale.UpdateTime = time.Now()
	//fieldNames = append(fieldNames, "update_time")
	autoscaleID, err := o.Update(&autoscale, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return autoscaleID, err
}

func DeleteAutoScale(autoscale model.ServiceAutoScale) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&autoscale)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func generateAutoScalesBySvcSQL(query model.ServiceAutoScale, svcID int64) (string, []interface{}) {
	var autoscaleBySvcSQL = `select  *
	from service_auto_scale
	where service_auto_scale.service_id = ?`

	params := make([]interface{}, 0)
	params = append(params, svcID)

	//	if query.ServiceCount != 0 {
	//		projectByUserSQL += ` and p.service_count = ?`
	//		params = append(params, query.ServiceCount)
	//	}

	return autoscaleBySvcSQL, params
}

func queryAutoScales(autoscalesBySvcSQL string, params []interface{}) ([]*model.ServiceAutoScale, error) {
	autoscales := make([]*model.ServiceAutoScale, 0)
	_, err := orm.NewOrm().Raw(autoscalesBySvcSQL, params).QueryRows(&autoscales)
	if err != nil {
		return nil, err
	}
	return autoscales, nil
}
func GetAutoScalesByService(query model.ServiceAutoScale, svcID int64) ([]*model.ServiceAutoScale, error) {
	autoscalesBySvcSQL, params := generateAutoScalesBySvcSQL(query, svcID)
	return queryAutoScales(autoscalesBySvcSQL, params)
}

// Sync autoscale from k8s to DB
func SyncAutoScaleData(autoscale model.ServiceAutoScale) (int64, error) {

	var asquery model.ServiceAutoScale
	asquery.HPAName = autoscale.HPAName
	o := orm.NewOrm()
	err := o.Read(&asquery, "name")
	if err != orm.ErrNoRows {
		return 0, nil
	}

	autoscaleID, err := o.Insert(&autoscale)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return autoscaleID, err
}
