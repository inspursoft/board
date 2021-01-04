package vm

import (
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"time"

	"github.com/astaxie/beego/logs"
)

type Project struct {
	ID                int64     `json:"project_id" orm:"column(id)"`
	Name              string    `json:"project_name" orm:"column(name)"`
	Comment           string    `json:"project_comment" orm:"column(comment)"`
	CreationTime      time.Time `json:"project_creation_time" orm:"column(creation_time)"`
	UpdateTime        time.Time `json:"project_update_time" orm:"column(update_time)"`
	Deleted           int       `json:"project_deleted" orm:"column(deleted)"`
	OwnerID           int       `json:"project_owner_id" orm:"column(owner_id)"`
	OwnerName         string    `json:"project_owner_name" orm:"column(owner_name)"`
	Public            int       `json:"project_public" orm:"column(public)"`
	Toggleable        bool      `json:"project_toggleable" orm:"column(toggleable)"`
	CurrentUserRoleID int64     `json:"project_current_user_role_id" orm:"column(current_user_role_id)"`
	ServiceCount      int       `json:"project_service_count" orm:"column(service_count)"`
	IstioSupport      bool      `json:"project_istio_support" orm:"column(istio_support)"`
	Deletable         bool      `json:"project_deletable,omitempty" orm:"-"`
}

func (p Project) ToMO() (m model.Project) {
	err := utils.Adapt(p, &m)
	if err != nil {
		logs.Error("Failed to convert VM to MO: %+v", err)
		return
	}
	logs.Debug("Converted to model: %+v", m)
	return
}
