package nodeDao

import (
	"git/inspursoft/board/src/adminserver/models/nodeModel"
	"github.com/astaxie/beego/orm"
)

func InsertNodeLog(nodeLog *nodeModel.NodeLog) (int64, error) {
	o := orm.NewOrm()
	var id int64
	var err error
	if id, err = o.Insert(nodeLog); err != nil {
		return 0, err
	}
	return id, nil
}

func GetNodeLog(id int) (*nodeModel.NodeLog, error) {
	o := orm.NewOrm()
	log := &nodeModel.NodeLog{Id: id}
	if err := o.Read(log, "id"); err != nil {
		return nil, err
	}
	return log, nil
}

func UpdateNodeLog(nodeLog *nodeModel.NodeLog) (error) {
	o := orm.NewOrm()
	if _, err := o.Update(nodeLog, "completed", "success"); err != nil {
		return err
	}
	return nil
}

func InsertNodeStatus(nodeStatus *nodeModel.NodeStatus) error {
	o := orm.NewOrm()
	if _, err := o.Insert(nodeStatus); err != nil {
		return err
	}
	return nil
}

func DeleteNodeStatus(nodeStatus *nodeModel.NodeStatus) error {
	o := orm.NewOrm()
	if _, err := o.Delete(nodeStatus, "ip"); err != nil {
		return err
	}
	return nil
}

func GetNodeLogList(nodeLogList *[]nodeModel.NodeLog, count int, offset int) error {
	o := orm.NewOrm()
	if _, err := o.QueryTable(&nodeModel.NodeLog{}).OrderBy("-id").Limit(count, offset).All(nodeLogList); err != nil {
		return err
	}
	return nil
}

func GetNodeStatusList(nodeStatusList *[]nodeModel.NodeStatus) error {
	o := orm.NewOrm()
	if _, err := o.QueryTable(&nodeModel.NodeStatus{}).All(nodeStatusList); err != nil {
		return err
	}
	return nil
}

func InsertNodeLogDetail(detail *nodeModel.NodeLogDetailInfo) (int64, error) {
	o := orm.NewOrm()
	var id int64
	var err error
	if id, err = o.Insert(detail); err != nil {
		return 0, err
	}
	return id, nil
}

func GetNodeLogDetail(logTimestamp int64) (*nodeModel.NodeLogDetailInfo, error) {
	o := orm.NewOrm()
	detail := &nodeModel.NodeLogDetailInfo{CreationTime: logTimestamp}
	if err := o.Read(detail, "creation_time"); err != nil {
		return nil, err
	}
	return detail, nil
}

func GetLogTotalRecordCount() (int64, error) {
	o := orm.NewOrm()
	var count int64
	err := o.Raw(`select count(*) from node_log`).QueryRow(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
