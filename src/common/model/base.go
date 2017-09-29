package model

import (
	"git/inspursoft/board/src/common/model/dashboard"

	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(User), new(Project), new(ProjectMember), new(Role),
		new(dashboard.Node), new(dashboard.Pod), new(dashboard.Service),
		new(dashboard.ServiceKvMap), new(dashboard.PodKvMap), new(dashboard.ServiceDashboardSecond),
		new(dashboard.ServiceDashboardMinute), new(dashboard.ServiceDashboardHour),
		new(dashboard.ServiceDashboardDay), new(dashboard.TimeListLog), new(ServiceStatus),
		new(ImageTag), new(Image))
}
