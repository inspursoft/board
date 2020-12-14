package auth_test

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/auth"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

const (
	adminUserID     = 1
	initialPassword = "123456a?"
	adminUsername   = "boardadmin"
	adminPassword   = "123456a?"
)

var (
	sshKeyPath = utils.GetConfig("SSH_KEY_PATH")
)

func updateAdminPassword() {
	salt := utils.GenerateRandomString()
	encryptedPassword := utils.Encrypt(initialPassword, salt)
	user := model.User{ID: adminUserID, Password: encryptedPassword, Salt: salt}
	isSuccess, err := service.UpdateUser(user, "password", "salt")
	if err != nil {
		logs.Error("Failed to update user password: %+v", err)
	}
	if isSuccess {
		logs.Info("Admin password has been updated successfully.")
	} else {
		logs.Info("Failed to update admin initial password.")
	}
}

func TestMain(m *testing.M) {
	utils.InitializeDefaultConfig()
	utils.AddValue("SSH_KEY_PATH", "/tmp/test-keys")
	dao.InitDB()
	updateAdminPassword()
	os.Exit(m.Run())
}

func TestSignIn(t *testing.T) {
	assert := assert.New(t)
	currentAuth, err := auth.GetAuth("db_auth")
	u, err := (*currentAuth).DoAuth(adminUsername, adminPassword)
	assert.Nil(err, "Error occurred while calling SignIn method.")
	assert.NotNil(u, "User is nil.")
	assert.Equal(adminUsername, u.Username, "Signed in failed.")
}

func TestSignInLdap(t *testing.T) {

	utils.SetConfig("LDAP_URL", fmt.Sprintf("ldap://%s", utils.GetStringValue("BOARD_HOST_IP")))
	utils.SetConfig("LDAP_SEARCH_DN", `cn=admin,dc=example,dc=org`)
	utils.SetConfig("LDAP_BASE_DN", "uid=test,dc=example,dc=org")
	utils.SetConfig("LDAP_FILTER", "")
	utils.SetConfig("LDAP_SEARCH_PWD", "admin")
	utils.SetConfig("LDAP_UID", "cn")
	utils.SetConfig("LDAP_SCOPE", "LDAP_SCOPE_SUBTREE")
	utils.SetConfig("LDAP_SCOPE", "")
	utils.SetConfig("LDAP_TIMEOUT", "5")
// 	assert := assert.New(t)
// 	currentAuth, err := auth.GetAuth("ldap_auth")
// 	u, err := (*currentAuth).DoAuth(`test`, `123456`)
// 	assert.Nil(err, "Error occurred while calling SignIn method.")
// 	assert.NotNil(u, "User is nil.")
}
