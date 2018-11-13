package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

type PVolumeController struct {
	BaseController
}

func (n *PVolumeController) GetPVolumeAction() {
	pvID, err := strconv.Atoi(n.Ctx.Input.Param(":id"))
	if err != nil {
		n.internalError(err)
		return
	}
	pv, err := service.GetPVDB(model.PersistentVolume{ID: int64(pvID)}, "id")
	if err != nil {
		n.internalError(err)
		return
	}
	if pv == nil {
		logs.Error("Not found this PV %d in DB", pvID)
		n.internalError(err)
		return
	}

	// TODO sync the state with K8S

	// To optimize the different types of common code
	switch pv.Type {
	case model.PVNFS:
		// PV NFS
		pvo, err := service.GetPVOptionNFS(model.PersistentVolumeOptionNfs{ID: int64(pvID)}, "id")
		if err != nil {
			n.internalError(err)
			return
		}
		if pv == nil {
			logs.Error("Not found this PV Option %d in DB", pvID)
			n.internalError(err)
			return
		}
		pv.Option = pvo

	case model.PVCephRBD:
		// PV CephRBD
		pvo, err := service.GetPVOptionRBD(model.PersistentVolumeOptionCephrbd{ID: int64(pvID)}, "id")
		if err != nil {
			n.internalError(err)
			return
		}
		if pv == nil {
			logs.Error("Not found this PV Option %d in DB", pvID)
			n.internalError(err)
			return
		}
		pv.Option = pvo
	default:
		logs.Error("Unknown pv type %d", pv.Type)
		n.customAbort(http.StatusBadRequest, "Unknown pv type")
		return
	}

	n.renderJSON(pv)
	logs.Debug("Return get pv %v", pv)
}

func (n *PVolumeController) RemovePVolumeAction() {
	pvID, err := strconv.Atoi(n.Ctx.Input.Param(":id"))
	if err != nil {
		n.internalError(err)
		return
	}
	pv, err := service.GetPVDB(model.PersistentVolume{ID: int64(pvID)}, "id")
	if err != nil {
		n.internalError(err)
		return
	}
	if pv == nil {
		logs.Debug("Not found this PV %d in DB", pvID)
		return
	}

	//TODO check and delete pv from k8s system

	switch pv.Type {
	case model.PVNFS:
		// PV NFS
		_, err = service.DeletePVOptionNFS(int64(pvID))
		if err != nil {
			logs.Error("Failed to delete PV NFS option %d", pvID)
			n.internalError(err)
			return
		}
	case model.PVCephRBD:
		_, err = service.DeletePVOptionRBD(int64(pvID))
		if err != nil {
			logs.Error("Failed to delete PV RBD option %d", pvID)
			n.internalError(err)
			return
		}

	default:
		logs.Error("Unknown pv type %d", pv.Type)
	}

	// Delete PV DB
	_, err = service.DeletePVDB(int64(pvID))
	if err != nil {
		logs.Error("Failed to delete PV %d", pvID)
		n.internalError(err)
		return
	}
}

func (n *PVolumeController) GetPVolumeListAction() {
	res, err := service.GetPVList()
	if err != nil {
		logs.Debug("Failed to get PV List")
		n.customAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	n.renderJSON(res)
}

func (n *PVolumeController) AddPVolumeAction() {
	var err error
	var message json.RawMessage
	pv := model.PersistentVolume{
		Option: &message,
	}
	data, err := ioutil.ReadAll(n.Ctx.Request.Body)
	if err != nil {
		n.customAbort(http.StatusBadRequest, "Invalid Request")
		return
	}
	err = json.Unmarshal(data, &pv)
	if err != nil {
		n.customAbort(http.StatusBadRequest, "Invalid Json Body")
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
			n.customAbort(http.StatusBadRequest, "Invalid PV NFS")
			return
		}

		logs.Debug("Receive pv nfs option %v", PVOptionNFS)
		err = service.AddPVolumeNFS(pv, PVOptionNFS)
		if err != nil {
			logs.Error("Failed to create nfs %v", err)
			n.customAbort(http.StatusBadRequest, "Invalid PV NFS")
			return
		}

	case model.PVCephRBD:
		// PV CephRBD
		var PVOptionRBD model.PersistentVolumeOptionCephrbd
		err = json.Unmarshal(message, &PVOptionRBD)
		if err != nil {
			logs.Error("Failed to unmarshal rbd %v", pv.Option)
			n.customAbort(http.StatusBadRequest, "Invalid PV RBD")
			return
		}

		logs.Debug("Receive pv rbd option %v", PVOptionRBD)
		logs.Debug("Receive rbd monitors %v", strings.Split(PVOptionRBD.Monitors, ","))
		err = service.AddPVolumeCephRBD(pv, PVOptionRBD)
		if err != nil {
			logs.Error("Failed to create rbc %v", err)
			n.customAbort(http.StatusBadRequest, "Invalid PV RBD")
			return
		}
	default:
		logs.Error("Unknown pv type %d", pv.Type)
		n.customAbort(http.StatusBadRequest, "Unknown pv type")
		return
	}

}
