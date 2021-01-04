package controller

import (
	"encoding/json"
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

type PVolumeController struct {
	c.BaseController
}

func (n *PVolumeController) GetPVolumeAction() {
	pvID, err := strconv.Atoi(n.Ctx.Input.Param(":id"))
	if err != nil {
		n.InternalError(err)
		return
	}
	pv, err := service.GetPVDB(model.PersistentVolume{ID: int64(pvID)}, "id")
	if err != nil {
		n.InternalError(err)
		return
	}
	if pv == nil {
		logs.Error("Not found this PV %d in DB", pvID)
		n.InternalError(err)
		return
	}

	// To optimize the different types of common code
	switch pv.Type {
	case model.PVNFS:
		// PV NFS
		pvo, err := service.GetPVOptionNFS(model.PersistentVolumeOptionNfs{ID: int64(pvID)}, "id")
		if err != nil {
			n.InternalError(err)
			return
		}
		if pv == nil {
			logs.Error("Not found this PV Option %d in DB", pvID)
			n.InternalError(err)
			return
		}
		pv.Option = pvo

	case model.PVCephRBD:
		// PV CephRBD
		pvo, err := service.GetPVOptionRBD(model.PersistentVolumeOptionCephrbd{ID: int64(pvID)}, "id")
		if err != nil {
			n.InternalError(err)
			return
		}
		if pv == nil {
			logs.Error("Not found this PV Option %d in DB", pvID)
			n.InternalError(err)
			return
		}
		pv.Option = pvo
	default:
		logs.Error("Unknown pv type %d", pv.Type)
		n.CustomAbortAudit(http.StatusBadRequest, "Unknown pv type")
		return
	}

	// sync the state with K8S

	pvk8s, err := service.GetPVK8s(pv.Name)
	if err != nil {
		logs.Error("Fail to get this PV %s in cluster %v", pv.Name, err)
		n.InternalError(err)
		return
	}
	if pvk8s == nil {
		pv.State = model.InvalidPV
	} else {
		pv.State = service.ReverseState(string(pvk8s.Status.Phase))
	}

	n.RenderJSON(pv)
	logs.Debug("Return get pv %v", pv)
}

func (n *PVolumeController) RemovePVolumeAction() {
	pvID, err := strconv.Atoi(n.Ctx.Input.Param(":id"))
	if err != nil {
		n.InternalError(err)
		return
	}
	pv, err := service.GetPVDB(model.PersistentVolume{ID: int64(pvID)}, "id")
	if err != nil {
		n.InternalError(err)
		return
	}
	if pv == nil {
		logs.Debug("Not found this PV %d in DB", pvID)
		return
	}

	err = service.DeletePVK8s(pv.Name)
	if err != nil {
		logs.Info("Delete PV %s from K8s Failed %v", pv.Name, err)
	} else {
		logs.Info("Delete PV %s from K8s Successful %v", pv.Name, err)
	}

	switch pv.Type {
	case model.PVNFS:
		// PV NFS
		_, err = service.DeletePVOptionNFS(int64(pvID))
		if err != nil {
			logs.Error("Failed to delete PV NFS option %d", pvID)
			n.InternalError(err)
			return
		}
	case model.PVCephRBD:
		_, err = service.DeletePVOptionRBD(int64(pvID))
		if err != nil {
			logs.Error("Failed to delete PV RBD option %d", pvID)
			n.InternalError(err)
			return
		}

	default:
		logs.Error("Unknown pv type %d", pv.Type)
	}

	// Delete PV DB
	_, err = service.DeletePVDB(int64(pvID))
	if err != nil {
		logs.Error("Failed to delete PV %d", pvID)
		n.InternalError(err)
		return
	}
}

func (n *PVolumeController) GetPVolumeListAction() {
	res, err := service.GetPVList()
	if err != nil {
		logs.Debug("Failed to get PV List")
		n.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.RenderJSON(res)
}

func (n *PVolumeController) AddPVolumeAction() {
	var err error
	var message json.RawMessage
	pv := model.PersistentVolume{
		Option: &message,
	}
	data, err := ioutil.ReadAll(n.Ctx.Request.Body)
	if err != nil {
		n.CustomAbortAudit(http.StatusBadRequest, "Invalid Request")
		return
	}
	err = json.Unmarshal(data, &pv)
	if err != nil {
		n.CustomAbortAudit(http.StatusBadRequest, "Invalid Json Body")
		return
	}
	logs.Debug("Add pv %v", pv)

	switch pv.Type {
	case model.PVNFS:
		// PV NFS
		var PVOptionNFS model.PersistentVolumeOptionNfs
		err = json.Unmarshal(message, &PVOptionNFS)
		if err != nil {
			logs.Error("Failed to unmarshal nfs %v", pv.Option)
			n.CustomAbortAudit(http.StatusBadRequest, "Invalid PV NFS")
			return
		}

		logs.Debug("Receive pv nfs option %v", PVOptionNFS)
		err = service.AddPVolumeNFS(pv, PVOptionNFS)
		if err != nil {
			logs.Error("Failed to create nfs %v", err)
			n.CustomAbortAudit(http.StatusBadRequest, "Invalid PV NFS")
			return
		}

	case model.PVCephRBD:
		// PV CephRBD
		var PVOptionRBD model.PersistentVolumeOptionCephrbd
		err = json.Unmarshal(message, &PVOptionRBD)
		if err != nil {
			logs.Error("Failed to unmarshal rbd %v", pv.Option)
			n.CustomAbortAudit(http.StatusBadRequest, "Invalid PV RBD")
			return
		}

		logs.Debug("Receive pv rbd option %v", PVOptionRBD)
		logs.Debug("Receive rbd monitors %v", strings.Split(PVOptionRBD.Monitors, ","))
		err = service.AddPVolumeCephRBD(pv, PVOptionRBD)
		if err != nil {
			logs.Error("Failed to create rbc %v", err)
			n.CustomAbortAudit(http.StatusBadRequest, "Invalid PV RBD")
			return
		}
	default:
		logs.Error("Unknown pv type %d", pv.Type)
		n.CustomAbortAudit(http.StatusBadRequest, "Unknown pv type")
		return
	}

}

func (n *PVolumeController) CheckPVolumeNameExistingAction() {
	pvName := n.GetString("pv_name")
	if pvName == "" {
		return
	}
	ispvk8s, err := service.GetPVK8s(pvName)
	if err != nil {
		n.InternalError(err)
		return
	}
	if ispvk8s != nil {
		n.CustomAbortAudit(http.StatusConflict, "This pv name is already existing in cluster.")
		return
	}

	ispvDB, err := service.GetPVDB(model.PersistentVolume{Name: pvName}, "name")
	if err != nil {
		n.InternalError(err)
		return
	}
	if ispvDB != nil {
		n.CustomAbortAudit(http.StatusConflict, "This pv name is already existing in DB.")
		return
	}
	logs.Info("PV name of %s is available", pvName)
}
