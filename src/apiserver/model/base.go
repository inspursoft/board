package model

import (
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(model.User))
}
