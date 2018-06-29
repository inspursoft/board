package dao_test

import (
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/astaxie/beego/logs"
)

var operationT1 = model.Operation{
	UserID:1, ProjectID:1, ProjectName: "library", ObjectType: "service", 
	ObjectName: "demoshow", Action:"delete", Status: "success"
}

func TestAddOperation(t *testing.T) {
	logs.Info("Start TestAddOperation")
	assert := assert.New(t)

	_, err := dao.AddOperation(operationT1)
	assert.Nil(err, "Should has no errors while executing config adding.")

	//c, _ := dao.GetConfig("auth_mode")
	//assert.NotNil(c, "Should not nil with finding this key: auth_mode")
	//assert.Equal(c.Value, "db_auth", "Should get value db_auth.")
	logs.Info("TestAddOperation Success")
}

func TestUpdateOperation(t *testing.T) {
    logs.Info("Start TestUpdateOperation")
	logs.Info("TestUpdateOperation Success")
}

/*
func TestGetOperation(t *testing.T) {
	assert := assert.New(t)
	config := model.Config{Name: "auth_mode", Value: "ldap_auth"}
	_, err := dao.AddOrUpdateConfig(config)
	assert.Nil(err, "Should has no errors while executing config updating.")

	c, _ := dao.GetConfig("auth_mode")
	assert.NotNil(c, "Should not nil with finding this key: auth_mode")
	assert.Equal(c.Value, "ldap_auth", "Should get value ldap_auth.")
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