package service

import (
	"fmt"
	"github.com/inspursoft/board/src/adminserver/common"
	"github.com/inspursoft/board/src/adminserver/dao"
	"github.com/inspursoft/board/src/adminserver/models"
	"github.com/inspursoft/board/src/common/model"
	t "github.com/inspursoft/board/src/common/token"
	"github.com/inspursoft/board/src/common/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
)

var TokenServerURL = fmt.Sprintf("http://%s:%s/tokenservice/token", "tokenserver", "4000")
var DefaultCacheDuration time.Duration

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
	var err error
	user := model.User{Username: acc.Username, SystemAdmin: 1, Deleted: 0}
	user, err = dao.LoginCheckAuth(user)
	if err != nil {
		return false, "", err
	}

	query := model.User{Username: acc.Username, Password: acc.Password}
	query.Password = utils.Encrypt(query.Password, user.Salt)
	query, err = dao.LoginCheckPassword(query)
	if err != nil {
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

	err = InitTokenCacheDuration()
	if err != nil {
		return false, "", err
	}
	dao.GlobalCache.Put(query.Username, token.TokenString, DefaultCacheDuration)
	dao.GlobalCache.Put(token.TokenString, payload, DefaultCacheDuration)

	err = UpdateApiserverCache(query.Username, token.TokenString)
	if err != nil {
		return false, "", err
	}
	err = UpdateApiserverCache(token.TokenString, payload)
	if err != nil {
		return false, "", err
	}

	return true, token.TokenString, nil
}

func GetCurrentUser(token string) (*model.User, string) {
	if isTokenExists := dao.GlobalCache.IsExist(token); !isTokenExists {
		logs.Info("Token stored in cache has expired.")
		return nil, ""
	}
	var hasResignedToken bool
	payload, err := t.VerifyToken(TokenServerURL, token)
	if err != nil {
		if err == t.ErrInvalidToken {
			if lastPayload, ok := dao.GlobalCache.Get(token).(map[string]interface{}); ok {
				newToken, err := t.SignToken(TokenServerURL, lastPayload)
				if err != nil {
					logs.Error("failed to sign token: %+v\n", err)
					return nil, ""
				}
				hasResignedToken = true
				token = newToken.TokenString
				payload = lastPayload
				dao.GlobalCache.Put(token, payload, DefaultCacheDuration)
				err = UpdateApiserverCache(token, payload)
				if err != nil {
					logs.Error("failed to update apiserver cache: %+v\n", err)
					return nil, ""
				}
				logs.Info("Token has been re-signed due to timeout.")
			}
		} else {
			logs.Error("failed to verify token: %+v\n", err)
			return nil, ""
		}
	}

	if strID, ok := payload["id"].(string); ok {
		userID, err := strconv.Atoi(strID)
		if err != nil {
			logs.Error("Error occurred on converting userID: %+v\n", err)
			return nil, ""
		}
		user := model.User{ID: int64(userID), Deleted: 0}
		user, err = dao.GetUserByID(user)
		if err != nil {
			logs.Error("Error occurred while getting user by ID: %d\n", err)
			return nil, ""
		}
		if currentToken, ok := dao.GlobalCache.Get(user.Username).(string); ok {
			if !hasResignedToken && currentToken != "" && currentToken != token {
				logs.Info("Another admin user has signed in other place.")
				return nil, ""
			}
			if currentToken != token {
				dao.GlobalCache.Put(user.Username, token, DefaultCacheDuration)
				err = UpdateApiserverCache(user.Username, token)
				if err != nil {
					logs.Error("failed to update apiserver cache: %+v\n", err)
					return nil, ""
				}
			}
		}
		user.Password = ""
		return &user, token
	}
	return nil, ""
}

func UpdateApiserverCache(key string, value interface{}) error {
	cache := make(map[string]interface{})
	cache["key"] = key
	cache["value"] = value
	return utils.RequestHandle(http.MethodPost, "http://apiserver:8088/api/v1/cache-store", nil, cache, nil)
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

func RemoveUUIDTokenCache() {
	UUIDpath := "/go/secrets/initialAdminPassword"
	if _, err := os.Stat(UUIDpath); !os.IsNotExist(err) {
		dao.RemoveUUIDToken()
		os.Remove(UUIDpath)
	}
}

func InitTokenCacheDuration() error {
	TokenCacheSeconds, err := common.ReadCfgItem("token_cache_expire_seconds", "/go/cfgfile/board.cfg")
	if err != nil {
		return err
	}
	TokenCacheSecondsNum, err := strconv.Atoi(TokenCacheSeconds)
	if err != nil {
		return err
	}
	DefaultCacheDuration = time.Second * time.Duration(TokenCacheSecondsNum)
	return err
}
