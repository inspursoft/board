package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

var sshKeyPath = utils.GetConfig("SSH_KEY_PATH")

const (
	sshPrivateKey = "id_rsa"
	sshPublicKey  = "id_rsa.pub"
)

func ConfigSSHAccess(username string, accessToken string) error {
	sshKeyUserPath := filepath.Join(sshKeyPath(), username)
	err := os.MkdirAll(sshKeyUserPath, 0755)
	if err != nil {
		return err
	}
	sshPrivateKeyPath := filepath.Join(sshKeyUserPath, sshPrivateKey)
	if _, err := os.Stat(sshPrivateKeyPath); os.IsNotExist(err) {
		err = exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", sshPrivateKeyPath, "-q", "-N", "").Run()
		if err != nil {
			return fmt.Errorf("Failed to generate SSH Key pairs: %+v", err)
		}
	}
	data, err := ioutil.ReadFile(filepath.Join(sshKeyUserPath, sshPublicKey))
	if err != nil {
		return err
	}
	publicKey := bytes.NewBuffer(data).String()
	return CurrentDevOps().ConfigSSHAccess(username, accessToken, publicKey)
}

func SignUp(user model.User) (bool, error) {
	err := CurrentDevOps().SignUp(user)
	if err != nil {
		return false, fmt.Errorf("Failed to create account for DevOps Git repository: %+v", err)
	}
	accessToken, err := CurrentDevOps().CreateAccessToken(user.Username, user.Password)
	if err != nil {
		return false, err
	}
	err = ConfigSSHAccess(user.Username, accessToken)
	if err != nil {
		return false, err
	}
	user.RepoToken = accessToken
	user.Salt = utils.GenerateRandomString()
	user.Password = utils.Encrypt(user.Password, user.Salt)
	userID, err := dao.AddUser(user)
	if err != nil {
		return false, err
	}
	return (userID != 0), nil
}

func GetUserByID(userID int64) (*model.User, error) {
	query := model.User{ID: userID, Deleted: 0}
	return dao.GetUser(query, "id", "deleted")
}

func GetUserByName(username string) (*model.User, error) {
	query := model.User{Username: username, Deleted: 0}
	return dao.GetUser(query, "username", "deleted")
}

func GetUserByEmail(email string) (*model.User, error) {
	query := model.User{Email: email, Deleted: 0}
	return dao.GetUser(query, "email", "deleted")
}

func GetUserByResetUUID(resetUuid string) (*model.User, error) {
	query := model.User{ResetUUID: resetUuid, Deleted: 0}
	return dao.GetUser(query, "reset_uuid", "deleted")
}

func GetUsers(field string, value interface{}, selectedFields ...string) ([]*model.User, error) {
	return dao.GetUsers(field, value, selectedFields...)
}

func GetPaginatedUsers(field string, value interface{}, pageIndex int, pageSize int, orderField string, orderAsc int, selectedField ...string) (*model.PaginatedUsers, error) {
	return dao.GetPaginatedUsers(field, value, pageIndex, pageSize, orderField, orderAsc, selectedField...)
}

func UpdateUser(user model.User, selectedFields ...string) (bool, error) {
	if user.ID == 0 {
		return false, errors.New("No user ID provided.")
	}
	_, err := dao.UpdateUser(user, selectedFields...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdateUserUUID(userID int64, resetUUID string) (bool, error) {
	return UpdateUser(model.User{ID: userID, ResetUUID: resetUUID}, "reset_uuid")
}

func ResetUserPassword(user model.User, newPassword string) (bool, error) {
	user.Password = utils.Encrypt(newPassword, user.Salt)
	user.ResetUUID = ""
	return UpdateUser(user, "password", "reset_uuid")
}

func DeleteUser(userID int64) (bool, error) {
	user, err := GetUserByID(userID)
	if err != nil {
		logs.Error("Failed to get user by ID: %d, error: %+v", userID, err)
		return false, nil
	}
	user.Username = "%" + user.Username + "%"
	user.Email = "%" + user.Email + "%"
	user.Deleted = 1
	_, err = dao.UpdateUser(*user, "username", "email", "deleted")
	if err != nil {
		logs.Error("Failed to update user: %v, error: %+v", user, err)
		return false, err
	}
	return true, nil
}

func UserExists(fieldName string, value string, userID int64) (bool, error) {
	query := model.User{ID: userID, Username: value, Email: value}
	user, err := dao.GetUser(query, fieldName)
	if err != nil {
		return false, err
	}
	if userID == 0 {
		return (user != nil && user.ID != 0), nil
	}
	return (user != nil && user.ID != userID), nil
}

func IsSysAdmin(userID int64) (bool, error) {
	query := model.User{ID: userID}
	user, err := dao.GetUser(query, "id")
	if err != nil {
		return false, err
	}
	return (user != nil && user.ID != 0 && user.SystemAdmin == 1), nil
}

//Use the Base64 Encode, to support others later
func DecodeUserPassword(password string) (string, error) {
	pwdBytes, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		logs.Error("Decode failed %s", password)
		return "", err
	}
	return string(pwdBytes), nil
}
