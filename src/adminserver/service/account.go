package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/encryption"
	"git/inspursoft/board/src/adminserver/models"
	t "git/inspursoft/board/src/common/token"
	"git/inspursoft/board/src/common/utils"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/alyu/configparser"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
)

var TokenServerURL = fmt.Sprintf("http://%s:%s/tokenservice/token", "tokenserver", "4000")
var TokenCacheExpireSeconds int

var ErrInvalidToken = errors.New("error for invalid token")

const (
	defaultInitialPassword = "123456a?"
	adminUserID            = 1
)

//VerifyPassword compares the password in cfg with the input one.
func VerifyPassword(passwd *models.Password) (bool, error) {

	configparser.Delimiter = "="
	cfgPath := path.Join("/go", "/cfgfile/board.cfg")
	//use configparser to read indicated cfg file.
	config, _ := configparser.Read(cfgPath)
	//section sensitive, global refers to all sections.
	section, _ := config.Section("global")
	password := section.ValueOf("board_admin_password")

	//ENCRYPTION
	prvKey, err := ioutil.ReadFile("./private.pem")
	if err != nil {
		return false, err
	}
	test, err := base64.StdEncoding.DecodeString(passwd.Value)
	if err != nil {
		return false, err
	}

	input := string(encryption.Decrypt("rsa", test, prvKey))

	return (input == password), nil
}

//Initialize save the account information into a file.
func Initialize(acc *models.Account) error {

	if acc.Password == "" {
		acc.Password = defaultInitialPassword
	}
	salt := utils.GenerateRandomString()
	encryptedPassword := utils.Encrypt(acc.Password, salt)
	user := models.User{ID: adminUserID, Username: acc.Username, Password: encryptedPassword, Salt: salt}
	o := orm.NewOrm()
	o.Using("mysql-db2")
	user.UpdateTime = time.Now()
	_, err := o.Update(&user, "password", "salt")
	if err != nil {
		return err
	}
	logs.Info("Admin password has been updated successfully.")

	o2 := orm.NewOrm()
	o2.Using("default")
	account := models.Account{Username: acc.Username, Password: acc.Password}
	if o2.Read(&models.Account{Id: 1}) == orm.ErrNoRows {
		if _, err := o2.Insert(&account); err != nil {
			return err
		}
	} else {
		if _, err := o2.Update(&account); err != nil {
			return err
		}
	}

	return nil
}

//Login allow user to use account information to login adminserver.
func Login(acc *models.Account) (bool, error, string) {
	o := orm.NewOrm()
	o.Using("mysql-db2")

	user := models.User{Username: acc.Username, SystemAdmin: 1, Deleted: 0}
	err := o.Read(&user, "username", "system_admin", "deleted")
	if err != nil {
		if err == orm.ErrNoRows {
			return false, errors.New("Forbidden"), ""
		}
		return false, err, ""
	}

	query := models.User{Username: acc.Username, Password: acc.Password}
	query.Password = utils.Encrypt(query.Password, user.Salt)
	err = o.Read(&query, "username", "password")
	if err != nil {
		return false, errors.New("Wrong password"), ""
	}

	payload := make(map[string]interface{})
	payload["id"] = strconv.Itoa(int(query.ID))
	payload["username"] = query.Username
	payload["email"] = query.Email
	payload["realname"] = query.Realname
	payload["is_system_admin"] = query.SystemAdmin
	token, err := t.SignToken(TokenServerURL, payload)
	if err != nil {
		return false, err, ""
	}

	TokenCacheExpireSeconds = 1800
	logs.Info("Set token server URL as %s and will expiration time after %d second(s) in cache.", TokenServerURL, TokenCacheExpireSeconds)
	dao.GlobalCache.Put(query.Username, token.TokenString, time.Second*time.Duration(TokenCacheExpireSeconds))
	dao.GlobalCache.Put(token.TokenString, payload, time.Second*time.Duration(TokenCacheExpireSeconds))

	return true, nil, token.TokenString
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
		o := orm.NewOrm()
		o.Using("mysql-db2")
		user := models.User{ID: int64(userID), Deleted: 0}
		err = o.Read(&user, "id", "deleted")
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
	existingToken := models.Token{Id: 1}
	newtoken := models.Token{Id: 1, Token: u, Time: time.Now().Unix()}
	o := orm.NewOrm()
	o.Using("default")
	if o.Read(&existingToken) == orm.ErrNoRows {
		if _, err := o.Insert(&newtoken); err != nil {
			return err
		}
	} else if (newtoken.Time - existingToken.Time) > 1800 {
		if _, err := o.Update(&newtoken); err != nil {
			return err
		}
	} else {
		return errors.New("another admin user has signed in other place")
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
func ValidateUUID(input string) (bool, error) {
	uuidPath := path.Join("/go", "/secrets/initialAdminPassword")
	f, err := ioutil.ReadFile(uuidPath)
	if err != nil {
		return false, err
	}

	return (input == string(f)), nil
}

func VerifyUUIDToken(input string) (bool, error) {
	o := orm.NewOrm()
	token := models.Token{Id: 1}
	err := o.Read(&token)
	if err != nil {
		logs.Error(err)
		return false, err
	}

	return (input == token.Token && (time.Now().Unix()-token.Time) <= 1800), nil
}
