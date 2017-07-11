package model

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(User), new(Project), new(ProjectMember), new(Role))
}
