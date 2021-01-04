package dao

import (
	"github.com/inspursoft/board/src/common/model"

	"time"

	//"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func AddNodeGroup(nodeGroup model.NodeGroup) (int64, error) {
	o := orm.NewOrm()

	nodeGroup.CreationTime = time.Now()
	nodeGroup.UpdateTime = nodeGroup.CreationTime

	nodeGroupID, err := o.Insert(&nodeGroup)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return nodeGroupID, err
}

func GetNodeGroup(nodeGroup model.NodeGroup, fieldNames ...string) (*model.NodeGroup, error) {
	o := orm.NewOrm()
	err := o.Read(&nodeGroup, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &nodeGroup, err
}

func UpdateNodeGroup(nodeGroup model.NodeGroup, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()
	nodeGroup.UpdateTime = time.Now()
	nodeGroupID, err := o.Update(&nodeGroup, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return nodeGroupID, err
}

func DeleteNodeGroup(nodeGroup model.NodeGroup) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&nodeGroup)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return num, err
}

func GetNodeGroups() ([]model.NodeGroup, error) {
	var nodeGroupList []model.NodeGroup //TODO new pointer make
	o := orm.NewOrm()
	_, err := o.QueryTable("node_group").All(&nodeGroupList)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return nodeGroupList, err
}
