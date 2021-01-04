package service

import (
	"errors"
	//"fmt"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/k8sassist"
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"
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

	//Sync state with cluster system

	for _, pvc := range pvcs {
		pvck8s, err := GetPVCK8s(pvc.Name, pvc.ProjectName)
		if err != nil {
			logs.Error("Fail to get this PVC %s in cluster %v", pvc.Name, err)
			// continue to work for other pvcs
			pvc.State = model.InvalidPVC
		}
		if pvck8s == nil {
			pvc.State = model.InvalidPVC
		} else {
			pvc.State = ReverseStatePVC(string(pvck8s.Status.Phase))
		}
	}

	return pvcs, err
}

func QueryPVCByID(pvcID int64) (*model.PersistentVolumeClaimV, error) {
	pvc, err := dao.QueryPVCByID(pvcID)
	if err != nil {
		logs.Debug("Failed to get pvc by id %d %v", pvcID, err)
		return nil, err
	}

	return pvc, err
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

func CreatePVCK8s(pvc *model.PersistentVolumeClaimK8scli) (*model.PersistentVolumeClaimK8scli, error) {
	// add the pvc to k8s
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	var err error
	newpvc, err := k8sclient.AppV1().PersistentVolumeClaim(pvc.Namespace).Create(pvc)
	if err != nil {
		logs.Debug("Failed to add PVC to K8s %v %v", pvc, err)
		return nil, err
	}
	logs.Debug("Added PVC to K8s")
	return newpvc, nil
}

func CreatePVC(pvc model.PersistentVolumeClaimM) (int64, error) {

	project, err := GetProjectByID(pvc.ProjectID)
	if err != nil {
		return 0, err
	}
	if project == nil {
		logs.Debug("Failed to find project %d", pvc.ProjectID)
	}

	// create pvc on cluster
	var pvcK8s model.PersistentVolumeClaimK8scli
	pvcK8s.Name = pvc.Name
	if pvc.PVName != "" {
		pvcK8s.Spec.VolumeName = pvc.PVName
	}
	pvcK8s.Namespace = project.Name
	pvcK8s.Spec.Resources.Requests = make(model.ResourceList)
	pvcK8s.Spec.Resources.Requests["storage"] = model.QuantityStr(pvc.Capacity)
	pvcK8s.Spec.AccessModes = append(pvcK8s.Spec.AccessModes, (model.PersistentVolumeAccessMode)(pvc.Accessmode))

	newpvc, err := CreatePVCK8s(&pvcK8s)
	if err != nil {
		logs.Debug("Failed to create pvc in cluster", pvcK8s)
		return 0, err
	}
	logs.Debug("Create PVC in cluster %v", newpvc)

	pvcID, err := CreatePVCDB(pvc)
	if err != nil {
		return 0, err
	}
	return pvcID, nil
}

func DeletePVCK8s(pvcname string, projectname string) error {

	// delete the hpa from k8s
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	err := k8sclient.AppV1().PersistentVolumeClaim(projectname).Delete(pvcname)
	if err != nil {
		if types.IsNotFoundError(err) {
			logs.Debug("Not found PVC %s", pvcname)
		} else {
			return err
		}
	}
	return nil
}

func GetPVCK8s(pvcname string, projectname string) (*model.PersistentVolumeClaimK8scli, error) {

	// delete the hpa from k8s
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	pvck8s, err := k8sclient.AppV1().PersistentVolumeClaim(projectname).Get(pvcname)
	if err != nil {
		if types.IsNotFoundError(err) {
			logs.Debug("Not found PV %s", pvcname)
			return nil, nil
		} else {
			return nil, err
		}
	}
	return pvck8s, nil
}

func ReverseStatePVC(state string) int {
	var ret = model.UnknownPVC
	switch state {
	case "Pending":
		ret = model.PendingPVC
	case "Bound":
		ret = model.BoundPVC
	case "Lost":
		ret = model.LostPVC
	}
	return ret
}

// check exsiting pvc name list in k8s by project, true is existed
func QueryPVCNameExisting(projectname string, pvcname string) (bool, error) {
	if projectname == "" {
		return false, errors.New("Project name is empty.")
	}

	if pvcname == "" {
		return false, errors.New("PVC name is empty.")
	}

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	pvclist, err := k8sclient.AppV1().PersistentVolumeClaim(projectname).List()
	if err != nil {
		logs.Error("Failed to get pvc %s", projectname)
		return false, err
	}
	for _, pvc := range pvclist.Items {
		if pvcname == pvc.Name {
			logs.Debug("pvc %s existing in %s", pvc.Name, projectname)
			return true, nil
		}
	}
	return false, nil
}
