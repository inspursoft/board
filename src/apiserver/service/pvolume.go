package service

import (
	"errors"
	//"fmt"

	"git/inspursoft/board/src/common/dao"
	//"git/inspursoft/board/src/common/k8sassist"
	//"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
)

func AddPVolumeNFS(pv model.PersistentVolume, pvo model.PersistentVolumeOptionNfs) error {

	//TODO k8s PV process
	pvID, err := CreatePVDB(pv)
	if err != nil {
		return err
	}
	logs.Debug("Create PV %d %s", pvID, pv.Name)
	pvo.ID = pvID
	pvID, err = CreatePVOptionNFS(pvo)
	if err != nil {
		return err
	}
	logs.Debug("Create PV Option %d", pvo.ID)
	return nil
}

func AddPVolumeCephRBD(pv model.PersistentVolume, pvo model.PersistentVolumeOptionCephrbd) error {

	//TODO k8s PV process
	pvID, err := CreatePVDB(pv)
	if err != nil {
		return err
	}
	logs.Debug("Create PV %d %s", pvID, pv.Name)
	pvo.ID = pvID
	pvID, err = CreatePVOptionRBD(pvo)
	if err != nil {
		return err
	}
	logs.Debug("Create PV Option %d", pvo.ID)
	return nil
}

func GetPVList() ([]model.PersistentVolume, error) {
	return dao.GetPVList()
}

func GetPVDB(pv model.PersistentVolume, selectedFields ...string) (*model.PersistentVolume, error) {
	p, err := dao.GetPV(pv, selectedFields...)
	if err != nil {
		return nil, err
	}
	return p, nil
}

//func CheckAutoScaleExist(svc *model.ServiceStatus, hpaname string) (bool, error) {
//	// get the hpaname from storage
//	ass, err := dao.GetAutoScalesByService(model.ServiceAutoScale{}, svc.ID)
//	if err != nil {
//		return false, err
//	}
//	for i := range ass {
//		if ass[i].HPAName == hpaname {
//			return true, nil
//		}
//	}
//	return false, nil
//}

// AutoScale in database
func CreatePVDB(pv model.PersistentVolume) (int64, error) {
	pvID, err := dao.AddPV(pv)
	if err != nil {
		return 0, err
	}
	return pvID, nil
}

func DeletePVDB(pvID int64) (bool, error) {
	s := model.PersistentVolume{ID: pvID}
	_, err := dao.DeletePV(s)
	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdatePVDB(pv model.PersistentVolume, fieldNames ...string) (bool, error) {
	if pv.ID == 0 {
		return false, errors.New("no AutoScale ID provided")
	}
	_, err := dao.UpdatePV(pv, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}

// TODO: code duplicate, need to optimize PV options

func CreatePVOptionNFS(pv model.PersistentVolumeOptionNfs) (int64, error) {
	pvID, err := dao.AddPVOptionNFS(pv)
	if err != nil {
		return 0, err
	}
	return pvID, nil
}

func DeletePVOptionNFS(pvID int64) (bool, error) {
	s := model.PersistentVolumeOptionNfs{ID: pvID}
	_, err := dao.DeletePVOptionNFS(s)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetPVOptionNFS(pv model.PersistentVolumeOptionNfs, selectedFields ...string) (*model.PersistentVolumeOptionNfs, error) {
	n, err := dao.GetPVOptionNFS(pv, selectedFields...)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func CreatePVOptionRBD(pv model.PersistentVolumeOptionCephrbd) (int64, error) {
	pvID, err := dao.AddPVOptionRBD(pv)
	if err != nil {
		return 0, err
	}
	return pvID, nil
}

func DeletePVOptionRBD(pvID int64) (bool, error) {
	s := model.PersistentVolumeOptionCephrbd{ID: pvID}
	_, err := dao.DeletePVOptionRBD(s)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetPVOptionRBD(pv model.PersistentVolumeOptionCephrbd, selectedFields ...string) (*model.PersistentVolumeOptionCephrbd, error) {
	n, err := dao.GetPVOptionRBD(pv, selectedFields...)
	if err != nil {
		return nil, err
	}
	return n, nil
}
