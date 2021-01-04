package dao

import (
	"github.com/inspursoft/board/src/common/model"

	//"time"

	//"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func AddPV(pv model.PersistentVolume) (int64, error) {
	o := orm.NewOrm()

	pvID, err := o.Insert(&pv)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return pvID, err
}

func GetPV(pv model.PersistentVolume, fieldNames ...string) (*model.PersistentVolume, error) {
	o := orm.NewOrm()
	err := o.Read(&pv, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &pv, err
}

func UpdatePV(pv model.PersistentVolume, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	//autoscale.UpdateTime = time.Now()
	//fieldNames = append(fieldNames, "update_time")
	pvID, err := o.Update(&pv, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return pvID, err
}

func DeletePV(pv model.PersistentVolume) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&pv)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func GetPVList() ([]model.PersistentVolume, error) {
	var pvList []model.PersistentVolume //TODO new pointer make
	o := orm.NewOrm()
	_, err := o.QueryTable("persistent_volume").All(&pvList)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return pvList, err
}

// TODO: code duplicate, need to optimize PV options

func AddPVOptionNFS(pv model.PersistentVolumeOptionNfs) (int64, error) {
	o := orm.NewOrm()

	pvID, err := o.Insert(&pv)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return pvID, err
}

func GetPVOptionNFS(pv model.PersistentVolumeOptionNfs, fieldNames ...string) (*model.PersistentVolumeOptionNfs, error) {
	o := orm.NewOrm()
	err := o.Read(&pv, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &pv, err
}

func UpdatePVOptionNFS(pv model.PersistentVolumeOptionNfs, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	//autoscale.UpdateTime = time.Now()
	//fieldNames = append(fieldNames, "update_time")
	pvID, err := o.Update(&pv, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return pvID, err
}

func DeletePVOptionNFS(pv model.PersistentVolumeOptionNfs) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&pv)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func AddPVOptionRBD(pv model.PersistentVolumeOptionCephrbd) (int64, error) {
	o := orm.NewOrm()

	pvID, err := o.Insert(&pv)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return pvID, err
}

func GetPVOptionRBD(pv model.PersistentVolumeOptionCephrbd, fieldNames ...string) (*model.PersistentVolumeOptionCephrbd, error) {
	o := orm.NewOrm()
	err := o.Read(&pv, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &pv, err
}

func UpdatePVOptionRBD(pv model.PersistentVolumeOptionCephrbd, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	//autoscale.UpdateTime = time.Now()
	//fieldNames = append(fieldNames, "update_time")
	pvID, err := o.Update(&pv, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return pvID, err
}

func DeletePVOptionRBD(pv model.PersistentVolumeOptionCephrbd) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&pv)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}
