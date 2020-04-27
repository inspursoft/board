package service

import (
	"fmt"
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/models"
	t "git/inspursoft/board/src/common/token"
	"git/inspursoft/board/src/common/utils"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
)

var TokenServerURL = fmt.Sprintf("http://%s:%s/tokenservice/token", "tokenserver", "4000")
var TokenCacheExpireSeconds int

const (
	defaultInitialPassword = "123456a?"
	adminUserID            = 1
)

func Login(acc *models.Account) (bool, string, error) {
	if err := CheckBoard(); err != nil {
		return ValidateUUID(acc.Password)
	} else {
		return LoginWithDB(acc)
	}
}

//LoginWithDB allow user to use account information to login adminserver.
func LoginWithDB(acc *models.Account) (bool, string, error) {

	user := models.User{Username: acc.Username, SystemAdmin: 1, Deleted: 0}
	if err := dao.LoginCheckAuth(user); err != nil {
		return false, "", err
	}

	query := models.User{Username: acc.Username, Password: acc.Password}
	query.Password = utils.Encrypt(query.Password, user.Salt)
	if err := dao.LoginCheckPassword(query); err != nil {
		return false, "", err
	}

	payload := make(map[string]interface{})
	payload["id"] = strconv.Itoa(int(query.ID))
	payload["username"] = query.Username
	payload["email"] = query.Email
	payload["realname"] = query.Realname
	payload["is_system_admin"] = query.SystemAdmin
	token, err := t.SignToken(TokenServerURL, payload)
	if err != nil {
		return false, "", err
	}

	TokenCacheExpireSeconds = 1800
	logs.Info("Set token server URL as %s and will expiration time after %d second(s) in cache.", TokenServerURL, TokenCacheExpireSeconds)
	dao.GlobalCache.Put(query.Username, token.TokenString, time.Second*time.Duration(TokenCacheExpireSeconds))
	dao.GlobalCache.Put(token.TokenString, payload, time.Second*time.Duration(TokenCacheExpireSeconds))

	return true, token.TokenString, nil
}

func GetCurrentUser(token string) *models.User {
	if isTokenExists := dao.GlobalCache.IsExist(token); !isTokenExists {
		logs.Info("Token stored in cache has expired.")
		return nil
	}
	payload, err := t.VerifyToken(TokenServerURL, token)
	if err != nil {
		logs.Error("failed to verify token: %+v\n", err)
		return nil
	}

	if strID, ok := payload["id"].(string); ok {
		userID, err := strconv.Atoi(strID)
		if err != nil {
			logs.Error("Error occurred on converting userID: %+v\n", err)
			return nil
		}
		user := models.User{ID: int64(userID), Deleted: 0}
		user, err = dao.GetUserByID(user)
		if err != nil {
			logs.Error("Error occurred while getting user by ID: %d\n", err)
			return nil
		}
		if currentToken, ok := dao.GlobalCache.Get(user.Username).(string); ok {
			if currentToken != "" && currentToken != token {
				logs.Info("Another admin user has signed in other place.")
				return nil
			}
		}
		user.Password = ""
		return &user
	}
	return nil
}

//CreateUUID creates a file with an UUID in it.
func CreateUUID() error {
	u := uuid.NewV4().String()
	newtoken := models.Token{Id: 1, Token: u, Time: time.Now().Unix()}
	err := dao.UpdateUUIDToken(newtoken)
	if err != nil {
		return err
	}

	folderPath := path.Join("/go", "/secrets")
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, os.ModePerm)
		os.Chmod(folderPath, os.ModePerm)
	}
	uuidPath := path.Join("/go", "/secrets/initialAdminPassword")
	f, err := os.Create(uuidPath)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(u)

	return nil
}

//ValidateUUID compares input with the UUID stored in the specified file.
func ValidateUUID(input string) (bool, string, error) {
	uuidPath := path.Join("/go", "/secrets/initialAdminPassword")
	f, err := ioutil.ReadFile(uuidPath)
	if err != nil {
		return false, "", err
	}

	return (input == string(f)), input, nil
}

func VerifyUUIDToken(input string) (bool, error) {
	token := models.Token{Id: 1}
	token, err := dao.GetUUIDToken(token)
	if err != nil {
		return false, err
	}

	return (input == token.Token && (time.Now().Unix()-token.Time) <= 1800), nil
}
