package service

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCheckDeploymentPath(t *testing.T) {
	assert := assert.New(t)
	err := CheckDeploymentPath("./tmp")
	assert.Nil(err, "Error occurred while testing CheckDeploymentPath.")
	deleteFile("./tmp")
}

func TestServiceExists(t *testing.T) {
	assert := assert.New(t)
	s, _ := ServiceExists("", "")
	assert.False(s, "Error occurred while testing ServiceExists.")
}
