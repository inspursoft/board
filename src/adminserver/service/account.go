package service

import (
	"encoding/base64"
	"errors"
	"git/inspursoft/board/src/adminserver/encryption"
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/common/utils"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/alyu/configparser"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
)

const (
	defaultInitialPassword = "123456a?"
	adminUserID            = 1
)

var configStorage map[string]interface{}

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

	if err = NextStep(models.InitStatusThird, models.InitStatusComplete); err != nil {
		return err
	}

	return nil
}

//Login allow user to use account information to login adminserver.
func Login(acc *models.Account) (bool, error, string) {
	var token string = ""
	o := orm.NewOrm()
	o.Using("mysql-db2")

	user := models.User{Username: acc.Username, SystemAdmin: 1, Deleted: 0}
	err := o.Read(&user, "username", "system_admin", "deleted")
	if err != nil {
		if err == orm.ErrNoRows {
			return false, errors.New("Forbidden"), token
		}
		return false, err, token
	}

	query := models.User{Username: acc.Username, Password: acc.Password}
	query.Password = utils.Encrypt(query.Password, user.Salt)
	err = o.Read(&query, "username", "password")

	if err != nil {
		return false, err, token
	}

	token = uuid.NewV4().String()
	newtoken := models.Token{Id: 1, Token: token, Time: time.Now().Unix()}
	o2 := orm.NewOrm()
	o2.Using("default")
	if o2.Read(&models.Token{Id: 1}) == orm.ErrNoRows {
		if _, err := o2.Insert(&newtoken); err != nil {
			return false, err, token
		}
	} else {
		if _, err := o2.Update(&newtoken); err != nil {
			return false, err, token
		}
	}

	return true, nil, token
}

//Install method is called when first open the admin server.
func Install() (models.InitStatus, error) {
	o := orm.NewOrm()
	status := models.InitStatusInfo{Id: 1}
	err := o.Read(&status)
	if err != nil {
		return -1, err
	}

	return status.Status, nil
}

//CreateUUID creates a file with an UUID in it.
func CreateUUID() error {
	u := uuid.NewV4()

	folderPath := path.Join("/go", "/secrets")
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, os.ModePerm)
		os.Chmod(folderPath, os.ModePerm)
	}

	uuidPath := path.Join("/go", "/secrets/initialAdminPassword")
	if _, err := os.Stat(uuidPath); os.IsNotExist(err) {
		f, err := os.Create(uuidPath)
		if err != nil {
			return err
		}
		f.WriteString(u.String())
		defer f.Close()
	}

	return nil
}

//ValidateUUID compares input with the UUID stored in the specified file.
func ValidateUUID(input string) (bool, error) {
	uuidPath := path.Join("/go", "/secrets/initialAdminPassword")
	f, err := ioutil.ReadFile(uuidPath)
	if err != nil {
		return false, err
	}

	result := (input == string(f))
	if result {
		os.Remove(uuidPath)
		if err = NextStep(models.InitStatusStart, models.InitStatusFirst); err != nil {
			return false, err
		}
	}
	return result, nil
}

func VerifyToken(input string) (bool, error) {
	o := orm.NewOrm()
	token := models.Token{Id: 1}
	err := o.Read(&token)
	if err != nil {
		logs.Error(err)
		return false, err
	}

	return (input == token.Token && (time.Now().Unix()-token.Time) <= 1800), nil
}
