package dao_test

import (
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/model"
	"testing"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var operationT1 = model.Operation{
	UserID: 1, ProjectID: 1, ProjectName: "library", ObjectType: "service",
	ObjectName: "demoshow", Action: "delete", Status: "unknown",
}

var operationT2 = model.Operation{
	UserID: 1, ProjectID: 1, ProjectName: "library", ObjectType: "image",
	ObjectName: "demoshow", Action: "delete", Status: "unknown",
}

var operationT1ID int64

func TestAddOperation(t *testing.T) {
	logs.Info("Start TestAddOperation")
	assert := assert.New(t)

	id, err := dao.AddOperation(operationT1)
	assert.Nil(err, "Should has no errors while executing config adding.")

	//c, _ := dao.GetConfig("auth_mode")
	//assert.NotNil(c, "Should not nil with finding this key: auth_mode")
	//assert.Equal(c.Value, "db_auth", "Should get value db_auth.")
	operationT1ID = id
	logs.Info("TestAddOperation Success id %d", operationT1ID)
}

func TestUpdateOperation(t *testing.T) {
	logs.Info("Start TestUpdateOperation")
	assert := assert.New(t)
	operationT1.Status = "success"
	operationT1.ID = operationT1ID
	number, err := dao.UpdateOperation(operationT1, "status")
	assert.Nil(err, "Failed to update the Operation test found err.")
	//assert.Equal(newOperation.Status, "success", "Failed to update the Operation test.")
	logs.Info(number)
	logs.Info("TestUpdateOperation Success")
}

func TestGetOperation(t *testing.T) {
	assert := assert.New(t)
	operation := model.Operation{ID: operationT1ID}

	c, _ := dao.GetOperation(operation, "id")
	assert.NotNil(c, "Should not nil with finding this key: operation")
	assert.Equal(c.Status, "success", "Should get value success.")
	logs.Info(c)
}

/*
func TestGetOperations(t *testing.T) {
	assert := assert.New(t)

	_, err := dao.AddOperation(operationT2)
	assert.Nil(err, "Should has no errors while executing config adding.")

	logs.Info("Test operation query 1")
	var operationQuery = model.Operation{
		ObjectType: "image",
		Status:     "unknown",
	}

	c, _ := dao.GetOperations(operationQuery, "", "")
	assert.NotNil(c, "Should not nil with finding this key: operation")
	for _, o := range c {
		logs.Info(o)
	}
	logs.Info(c)

	logs.Info("Test operation query 2")
	var operationQuery2 = model.Operation{
		Action: "delete",
	}

	c, _ = dao.GetOperations(operationQuery2, "", "")
	assert.NotNil(c, "Should not nil with finding this key: operation")
	for _, o := range c {
		logs.Info(o)
	}

	logs.Info("Test operation query 3")
	var operationQuery3 model.Operation
	c, _ = dao.GetOperations(operationQuery3, "", time.Now().String())
	assert.NotNil(c, "Should not nil with finding this key: operation")
	for _, o := range c {
		logs.Info(o)
	}

	logs.Info("Test operation query 4")
	c, _ = dao.GetOperations(operationQuery3, time.Now().String(), "")
	assert.NotNil(c, "Should not nil with finding this key: operation")
	for _, o := range c {
		logs.Info(o)
	}

	logs.Info("TestGetOperations PASS")
}

func TestDeleteOperation(t *testing.T) {
	assert := assert.New(t)
	key := "auth_mode"
	_, err := dao.DeleteConfig(key)
	assert.Nil(err, "Should has no errors while executing config deleting.")

	c, _ := dao.GetConfig("auth_mode")
	assert.Equal(c.Name, "", "Should nil with finding this key: auth_mode")
}
*/
