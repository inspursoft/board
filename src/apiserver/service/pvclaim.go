package service

import (
	"errors"
	//"fmt"
	"git/inspursoft/board/src/common/dao"
	//"git/inspursoft/board/src/common/k8sassist"
	//"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"
	//"strings"

	"github.com/astaxie/beego/logs"
)

func GetPVCList() ([]model.PersistentVolumeClaimM, error) {
	pvcList, err := dao.GetPVCList()
	if err != nil {
		return nil, err
	}
	return pvcList, nil
}

func QueryPVCsByUser(userID int64) ([]*model.PersistentVolumeClaimV, error) {
	pvcs := make([]*model.PersistentVolumeClaimV, 0)
	query := model.Project{}
	projectList, err := GetProjectsByMember(query, userID)
	if err != nil {
		logs.Debug("Failed to get projects %d %v", userID, err)
		return nil, err
	}
	for _, p := range projectList {
		querypvcs, err := dao.QueryPVCByProjectID(p.ID)
		if err != nil {
			logs.Debug("Failed to get pvc by id %d %v", p.ID, err)
			return nil, err
		}
		pvcs = append(pvcs, querypvcs...)
	}
	logs.Debug("Guery PVC list %v", pvcs)
	return pvcs, err
}

func GetPVCDB(pvc model.PersistentVolumeClaimM, selectedFields ...string) (*model.PersistentVolumeClaimM, error) {
	p, err := dao.GetPVC(pvc, selectedFields...)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func CreatePVCDB(pvc model.PersistentVolumeClaimM) (int64, error) {
	pvcID, err := dao.AddPVC(pvc)
	if err != nil {
		return 0, err
	}
	return pvcID, nil
}

func DeletePVCDB(pvcID int64) (bool, error) {
	s := model.PersistentVolumeClaimM{ID: pvcID}
	_, err := dao.DeletePVC(s)
	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdatePVCDB(pvc model.PersistentVolumeClaimM, fieldNames ...string) (bool, error) {
	if pvc.ID == 0 {
		return false, errors.New("no PVC ID provided")
	}
	_, err := dao.UpdatePVC(pvc, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}

//func CreatePVK8s(pv *model.PersistentVolumeK8scli) (*model.PersistentVolumeK8scli, error) {
//	// add the pv to k8s
//	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
//		KubeConfigPath: kubeConfigPath(),
//	})
//	var err error
//	newpv, err := k8sclient.AppV1().PersistentVolume().Create(pv)
//	if err != nil {
//		logs.Debug("Failed to add PV to K8s %v %v", pv, err)
//		return nil, err
//	}
//	logs.Debug("Added PV to K8s")
//	return newpv, nil
//}

//func DeletePVK8s(pvname string) error {

//	// delete the hpa from k8s
//	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
//		KubeConfigPath: kubeConfigPath(),
//	})
//	err := k8sclient.AppV1().PersistentVolume().Delete(pvname)
//	if err != nil {
//		if types.IsNotFoundError(err) {
//			logs.Debug("Not found PV %s", pvname)
//		} else {
//			return err
//		}
//	}
//	return nil
//}

//func GetPVK8s(pvname string) (*model.PersistentVolumeK8scli, error) {

//	// delete the hpa from k8s
//	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
//		KubeConfigPath: kubeConfigPath(),
//	})
//	pvk8s, err := k8sclient.AppV1().PersistentVolume().Get(pvname)
//	if err != nil {
//		if types.IsNotFoundError(err) {
//			logs.Debug("Not found PV %s", pvname)
//			return nil, nil
//		} else {
//			return nil, err
//		}
//	}
//	return pvk8s, nil
//}

//func genPersistentVolumeK8scli(pv model.PersistentVolume, pvK8s *model.PersistentVolumeK8scli) *model.PersistentVolumeK8scli {
//	pvK8s.Name = pv.Name
//	pvK8s.Labels = make(map[string]string)
//	pvK8s.Labels["pvname"] = pv.Name
//	pvK8s.Spec.Capacity = make(model.ResourceList)
//	pvK8s.Spec.Capacity["storage"] = model.QuantityStr(pv.Capacity)
//	pvK8s.Spec.AccessModes = append(pvK8s.Spec.AccessModes, (model.PersistentVolumeAccessMode)(pv.Accessmode))
//	pvK8s.Spec.PersistentVolumeReclaimPolicy = (model.PersistentVolumeReclaimPolicy)(pv.Reclaim)
//	return pvK8s
//}
