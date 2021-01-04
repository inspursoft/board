package models

import (
	"github.com/inspursoft/board/src/adminserver/models/nodeModel"
	"github.com/inspursoft/board/src/common/model"

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
