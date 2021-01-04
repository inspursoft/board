package controller

import (
	//"encoding/json"
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"

	//"io/ioutil"
	"net/http"
	"strconv"

	//"strings"

	"github.com/astaxie/beego/logs"
)

type PVClaimController struct {
	c.BaseController
}

func (n *PVClaimController) GetPVClaimAction() {
	pvcID, err := strconv.Atoi(n.Ctx.Input.Param(":id"))
	if err != nil {
		n.InternalError(err)
		return
	}
	pvc, err := service.QueryPVCByID(int64(pvcID))
	if err != nil {
		n.InternalError(err)
		return
	}
	if pvc == nil {
		logs.Error("Not found this PVC %d in DB", pvcID)
		n.InternalError(err)
		return
	}

	var pvcDetail model.PersistentVolumeClaimDetail
	pvcDetail.PVClaim = *pvc

	//sync the state with K8S

	pvck8s, err := service.GetPVCK8s(pvc.Name, pvc.ProjectName)
	if err != nil {
		logs.Error("Fail to get this PVC %s in cluster %v", pvc.Name, err)
		n.InternalError(err)
		return
	}
	if pvck8s == nil {
		pvcDetail.PVClaim.State = model.InvalidPVC // TODO duplicate
		pvcDetail.State = model.InvalidPVC
	} else {
		pvcDetail.State = service.ReverseStatePVC(string(pvck8s.Status.Phase))
		pvcDetail.PVClaim.State = pvcDetail.State // TODO duplicate
	}

	n.RenderJSON(pvcDetail)
	logs.Debug("Return get pvc %v", pvcDetail)
}

func (n *PVClaimController) RemovePVClaimAction() {
	pvcID, err := strconv.Atoi(n.Ctx.Input.Param(":id"))
	if err != nil {
		n.InternalError(err)
		return
	}
	//check pvc existing DB
	pvc, err := service.GetPVCDB(model.PersistentVolumeClaimM{ID: int64(pvcID)}, "id")
	if err != nil {
		n.InternalError(err)
		return
	}
	if pvc == nil {
		logs.Debug("Not found this PVC %d in DB", pvcID)
		return
	}

	// Get PVC view mode
	pvcview, err := service.QueryPVCByID(int64(pvcID))
	if err != nil {
		n.InternalError(err)
		return
	}

	err = service.DeletePVCK8s(pvcview.Name, pvcview.ProjectName)
	if err != nil {
		logs.Info("Delete PVC %s from K8s Failed %v", pvcview.Name, err)
	} else {
		logs.Info("Delete PVC %s from K8s Successful %v", pvcview.Name, err)
	}

	// Delete PV DB
	_, err = service.DeletePVCDB(int64(pvcID))
	if err != nil {
		logs.Error("Failed to delete PVC %d", pvcID)
		n.InternalError(err)
		return
	}
}

func (n *PVClaimController) GetPVClaimListAction() {
	res, err := service.QueryPVCsByUser(n.CurrentUser.ID)
	if err != nil {
		logs.Debug("Failed to get PV List")
		n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.RenderJSON(res)
}

func (n *PVClaimController) AddPVClaimAction() {
	var reqPVC model.PersistentVolumeClaimM
	var err error
	err = n.ResolveBody(&reqPVC)
	if err != nil {
		return
	}

	if reqPVC.Name == "" || reqPVC.ProjectID == 0 {
		n.CustomAbortAudit(http.StatusBadRequest, "PVC Name and project ID should not null")
		return
	}

	//	pvcExists, err := service.PVCExists(reqPVC.ProjectID, reqPVC.Name)
	//	if err != nil {
	//		n.InternalError(err)
	//		return
	//	}
	//	if pvcExists {
	//		n.CustomAbortAudit(http.StatusConflict, "Node Group name already exists.")
	//		return
	//	}

	pvcID, err := service.CreatePVC(reqPVC)
	if err != nil {
		logs.Debug("Failed to add pvc %v", reqPVC)
		n.InternalError(err)
		return
	}
	logs.Info("Added PVC %s %d", reqPVC.Name, pvcID)
}

// Get the exsiting PVC name for checking
func (n *PVClaimController) GetPVCNameExisting() {
	projectName := n.GetString("project_name")
	if projectName == "" {
		logs.Debug("Failed to get Project name")
		n.CustomAbortAudit(http.StatusBadRequest, "No project name")
		return
	}

	pvcName := n.GetString("pvc_name")
	if projectName == "" {
		logs.Debug("Failed to get PVC name")
		n.CustomAbortAudit(http.StatusBadRequest, "No pvc name")
		return
	}

	res, err := service.QueryPVCNameExisting(projectName, pvcName)
	if err != nil {
		logs.Debug("Failed to check PVC name")
		n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	if res == true {
		n.CustomAbortAudit(http.StatusConflict, "The pvc name existing")
		return
	}
}
