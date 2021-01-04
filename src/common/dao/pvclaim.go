package dao

import (
	"github.com/inspursoft/board/src/common/model"

	//"time"

	//"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func AddPVC(pvc model.PersistentVolumeClaimM) (int64, error) {
	o := orm.NewOrm()

	pvcID, err := o.Insert(&pvc)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return pvcID, err
}

func GetPVC(pvc model.PersistentVolumeClaimM, fieldNames ...string) (*model.PersistentVolumeClaimM, error) {
	o := orm.NewOrm()
	err := o.Read(&pvc, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &pvc, err
}

func UpdatePVC(pvc model.PersistentVolumeClaimM, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	//autoscale.UpdateTime = time.Now()
	//fieldNames = append(fieldNames, "update_time")
	pvcID, err := o.Update(&pvc, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return pvcID, err
}

func DeletePVC(pvc model.PersistentVolumeClaimM) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&pvc)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func GetPVCList() ([]model.PersistentVolumeClaimM, error) {
	var pvcList []model.PersistentVolumeClaimM //TODO new pointer make
	o := orm.NewOrm()
	_, err := o.QueryTable("persistent_volume_claim_m").All(&pvcList)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return pvcList, err
}

func QueryPVCByProjectID(projectID int64) ([]*model.PersistentVolumeClaimV, error) {

	var pvcByProjectSQL = `select pvc.id, pvc.name, pvc.projectid, p.name as projectname, 
	pvc.capacity, pvc.accessmode, pvc.class, 
	pvc.pvname
	from persistent_volume_claim_m pvc 
	left join project p on pvc.projectid = p.id
	where pvc.projectid = ?`

	pvcs := make([]*model.PersistentVolumeClaimV, 0)
	_, err := orm.NewOrm().Raw(pvcByProjectSQL, projectID).QueryRows(&pvcs)
	if err != nil {
		return nil, err
	}
	return pvcs, nil
}

func QueryPVCByID(pvcID int64) (*model.PersistentVolumeClaimV, error) {

	var pvcByIDSQL = `select pvc.id, pvc.name, pvc.projectid, p.name as projectname, 
	pvc.capacity, pvc.accessmode, pvc.class, 
	pvc.pvname
	from persistent_volume_claim_m pvc 
	left join project p on pvc.projectid = p.id
	where pvc.id = ?`

	var pvc model.PersistentVolumeClaimV
	err := orm.NewOrm().Raw(pvcByIDSQL, pvcID).QueryRow(&pvc)
	if err != nil {
		return nil, err
	}
	return &pvc, nil
}
