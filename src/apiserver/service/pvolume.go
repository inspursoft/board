package service

import (
	"errors"
	//"fmt"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/k8sassist"
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"
	"strings"

	"github.com/astaxie/beego/logs"
)

func AddPVolumeNFS(pv model.PersistentVolume, pvo model.PersistentVolumeOptionNfs) error {

	//k8s PV process
	var pvk8s model.PersistentVolumeK8scli
	var pvoption model.NFSVolumeSource

	genPersistentVolumeK8scli(pv, &pvk8s)
	pvoption.Path = pvo.Path
	pvoption.Server = pvo.Server
	pvoption.ReadOnly = pv.Readonly
	pvk8s.Spec.PersistentVolumeSource = model.PersistentVolumeSource{
		NFS: &pvoption,
	}

	newpvk8s, err := CreatePVK8s(&pvk8s)
	if err != nil {
		logs.Error("Failed to add PV to K8s %v %v", pvk8s, err)
		return err
	}
	logs.Info("Created PV in K8s %v", newpvk8s)

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
	var pvk8s model.PersistentVolumeK8scli
	var pvoption model.RBDPersistentVolumeSource

	genPersistentVolumeK8scli(pv, &pvk8s)
	pvoption.FSType = pvo.Fstype
	pvoption.CephMonitors = strings.Split(pvo.Monitors, ",")
	pvoption.ReadOnly = pv.Readonly
	pvoption.Keyring = pvo.Keyring
	pvoption.RBDImage = pvo.Image
	pvoption.RadosUser = pvo.User
	pvoption.RBDPool = pvo.Pool
	// TODO support secret
	// pvoption.SecretRef = pvo.Secretname
	pvk8s.Spec.PersistentVolumeSource = model.PersistentVolumeSource{
		RBD: &pvoption,
	}

	newpvk8s, err := CreatePVK8s(&pvk8s)
	if err != nil {
		logs.Error("Failed to add PV to K8s %v %v", pvk8s, err)
		return err
	}
	logs.Info("Created PV in K8s %v", newpvk8s)

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
	pvList, err := dao.GetPVList()
	if err != nil {
		return nil, err
	}
	for i, pv := range pvList {
		pvk8s, err := GetPVK8s(pv.Name)
		if err != nil {
			logs.Error("Fail to get this PV %s in cluster %v", pv.Name, err)
			pvList[i].State = model.UnknownPV
		}
		if pvk8s == nil {
			pvList[i].State = model.InvalidPV
		} else {
			pvList[i].State = ReverseState(string(pvk8s.Status.Phase))
		}
	}
	return pvList, nil
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

func CreatePVK8s(pv *model.PersistentVolumeK8scli) (*model.PersistentVolumeK8scli, error) {
	// add the pv to k8s
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	var err error
	newpv, err := k8sclient.AppV1().PersistentVolume().Create(pv)
	if err != nil {
		logs.Debug("Failed to add PV to K8s %v %v", pv, err)
		return nil, err
	}
	logs.Debug("Added PV to K8s")
	return newpv, nil
}

func DeletePVK8s(pvname string) error {

	// delete the hpa from k8s
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	err := k8sclient.AppV1().PersistentVolume().Delete(pvname)
	if err != nil {
		if types.IsNotFoundError(err) {
			logs.Debug("Not found PV %s", pvname)
		} else {
			return err
		}
	}
	return nil
}

func GetPVK8s(pvname string) (*model.PersistentVolumeK8scli, error) {

	// delete the hpa from k8s
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	pvk8s, err := k8sclient.AppV1().PersistentVolume().Get(pvname)
	if err != nil {
		if types.IsNotFoundError(err) {
			logs.Debug("Not found PV %s", pvname)
			return nil, nil
		} else {
			return nil, err
		}
	}
	return pvk8s, nil
}

func genPersistentVolumeK8scli(pv model.PersistentVolume, pvK8s *model.PersistentVolumeK8scli) *model.PersistentVolumeK8scli {
	pvK8s.Name = pv.Name
	pvK8s.Labels = make(map[string]string)
	pvK8s.Labels["pvname"] = pv.Name
	pvK8s.Spec.Capacity = make(model.ResourceList)
	pvK8s.Spec.Capacity["storage"] = model.QuantityStr(pv.Capacity)
	pvK8s.Spec.AccessModes = append(pvK8s.Spec.AccessModes, (model.PersistentVolumeAccessMode)(pv.Accessmode))
	pvK8s.Spec.PersistentVolumeReclaimPolicy = (model.PersistentVolumeReclaimPolicy)(pv.Reclaim)
	return pvK8s
}

func ReverseState(state string) int {
	var ret = model.UnknownPV
	switch state {
	case "Pending":
		ret = model.PendingPV
	case "Available":
		ret = model.AvailablePV
	case "Bound":
		ret = model.BoundPV
	case "Released":
		ret = model.ReleasedPV
	case "Failed":
		ret = model.FailedPV
	}
	return ret
}
