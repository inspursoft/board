package controller

import (
	//"fmt"
	//"git/inspursoft/board/src/apiserver/service"
	"encoding/json"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

type PVolumeController struct {
	BaseController
}

func (n *PVolumeController) GetPVolumeAction() {
	//TODO
}

func (n *PVolumeController) GetPVolumeListAction() {
	//TODO
}

func (n *PVolumeController) AddPVolumeAction() {
	var reqPVolume model.PersistentVolume
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
		//TODO service.AddPVolumeNFS(reqPVolume, PVOptionNFS)

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
		//TODO service.AddPVolumeCephRBD(reqPVolume, PVOptionRBD)
	default:
		logs.Error("Unknown pv type %d", reqPVolume.Type)
		n.customAbort(http.StatusBadRequest, "Unknown pv type")
		return
	}

}
