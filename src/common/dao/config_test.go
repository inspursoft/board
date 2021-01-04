package dao_test

import (
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/model"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestAddConfig(t *testing.T) {
	assert := assert.New(t)
	config := model.Config{Name: "auth_mode", Value: "db_auth"}
	_, err := dao.AddOrUpdateConfig(config)
	assert.Nil(err, "Should has no errors while executing config adding.")

	c, _ := dao.GetConfig("auth_mode")
	assert.NotNil(c, "Should not nil with finding this key: auth_mode")
	assert.Equal(c.Value, "db_auth", "Should get value db_auth.")
}

func TestUpdateConfig(t *testing.T) {
	assert := assert.New(t)
	config := model.Config{Name: "auth_mode", Value: "ldap_auth"}
	_, err := dao.AddOrUpdateConfig(config)
	assert.Nil(err, "Should has no errors while executing config updating.")

	c, _ := dao.GetConfig("auth_mode")
	assert.NotNil(c, "Should not nil with finding this key: auth_mode")
	assert.Equal(c.Value, "ldap_auth", "Should get value ldap_auth.")
}

func TestDeleteConfig(t *testing.T) {
	assert := assert.New(t)
	key := "auth_mode"
	_, err := dao.DeleteConfig(key)
	assert.Nil(err, "Should has no errors while executing config deleting.")

	c, _ := dao.GetConfig("auth_mode")
	assert.Equal(c.Name, "", "Should nil with finding this key: auth_mode")
}
