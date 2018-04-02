package jenkins

import (
	"git/inspursoft/board/src/common/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepare() {

}

func cleanUp() {

}

func TestMain(m *testing.M) {
	os.Exit(func() int {
		prepare()
		m.Run()
		cleanUp()
		return 0
	}())
}

var user = model.User{
	Username: "tester",
	Password: "123456a?",
	Email:    "tester@inspur.com",
}

func TestCreateJob(t *testing.T) {
	err := NewJenkinsHandler(user.Username, "").CreateJob("target10")
	assert := assert.New(t)
	assert.Nilf(err, "Error occurred while creating job: %+v", err)
}
