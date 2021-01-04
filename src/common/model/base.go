package model

import (
	"github.com/inspursoft/board/src/common/model/dashboard"

	"github.com/astaxie/beego/orm"
)

func InitModelDB() {
	orm.RegisterModel(new(User), new(Project), new(ProjectMember), new(Role),
		new(dashboard.Node), new(dashboard.Pod), new(dashboard.Service),
		new(dashboard.ServiceKvMap), new(dashboard.PodKvMap), new(dashboard.ServiceDashboardSecond),
		new(dashboard.ServiceDashboardMinute), new(dashboard.ServiceDashboardHour),
		new(dashboard.ServiceDashboardDay), new(dashboard.TimeListLog), new(ServiceStatus),
		new(ImageTag), new(Image), new(NodeGroup), new(Config), new(Operation), new(ServiceAutoScale),
		new(PersistentVolume), new(PersistentVolumeOptionNfs), new(PersistentVolumeOptionCephrbd),
		new(PersistentVolumeClaimM), new(PersistentVolumeClaimV), new(HelmRepository), new(ReleaseModel), new(JobStatusMO))
}
