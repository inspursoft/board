package models

import (
	"git/inspursoft/board/src/adminserver/models/nodeModel"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/orm"
)

func RegisterModels() {
	orm.RegisterModel(
		new(nodeModel.NodeLog),
		new(nodeModel.NodeLogDetailInfo),
		new(nodeModel.NodeStatus),
		new(Account),
		new(Token),
		new(model.User))
}
