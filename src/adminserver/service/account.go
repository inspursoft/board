package service

import (
	"encoding/base64"
	"encoding/hex"
	"git/inspursoft/board/src/adminserver/encryption"
	"git/inspursoft/board/src/adminserver/models"
	"io/ioutil"
	"log"
	"os"
	"path"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"github.com/alyu/configparser"
	uuid "github.com/satori/go.uuid"
	"fmt"
	"time"
)

//VerifyPassword compares the password in cfg with the input one.
func VerifyPassword(passwd *models.Password) (a bool, err string) {
	var statusMessage string = "OK"

	configparser.Delimiter = "="
	cfgPath := path.Join(os.Getenv("GOPATH"), "/cfgfile/board.cfg")
	//use configparser to read indicated cfg file.
	config, _ := configparser.Read(cfgPath)
	//section sensitive, global refers to all sections.
	section, _ := config.Section("global")
	password := section.ValueOf("board_admin_password")

	//ENCRYPTION
	prvKey, err0 := ioutil.ReadFile("./private.pem")
	if err0 != nil {
		log.Print(err0)
		statusMessage = "BadRequest"
	}
	test, err1 := base64.StdEncoding.DecodeString(passwd.Value)
	if err1 != nil {
		log.Print(err1)
		statusMessage = "BadRequest"
	}

	input := string(encryption.Decrypt("rsa", test, prvKey))

	return (input == password), statusMessage
}

//Initialize save the account information into a file.
func Initialize(acc *models.Account) string {
	var statusMessage string = "OK"

	_, pubKey := encryption.GenKey("rsa")
	ciphertext := encryption.Encrypt("rsa", []byte(acc.Password), pubKey)

	err1 := os.Rename("./private.pem", "./private_acc.pem")
	if err1 != nil {
		log.Print(err1)
		statusMessage = "BadRequest"
	}

	accPath := path.Join(os.Getenv("GOPATH"), "/secrets/account-info")
	f, err2 := os.Create(accPath)
	if err2 != nil {
		log.Print(err2)
		statusMessage = "BadRequest"
	}
	f.WriteString("username = " + acc.Username + "\n")
	f.WriteString("password = " + hex.EncodeToString(ciphertext) + "\n")
	defer f.Close()

	return statusMessage
}

//Login allow user to use account information to login adminserver.
func Login(acc *models.Account) (a bool, b string) {
	var statusMessage string = "OK"
	var permission bool
	configparser.Delimiter = "="
	accPath := path.Join(os.Getenv("GOPATH"), "/secrets/account-info")
	config, err0 := configparser.Read(accPath)
	if err0 != nil {
		log.Print(err0)
		statusMessage = "BadRequest"
	}
	section, err1 := config.Section("global")
	if err1 != nil {
		log.Print(err1)
		statusMessage = "BadRequest"
	}
	username := section.ValueOf("username")
	ciphertext := section.ValueOf("password")

	prvKey, err2 := ioutil.ReadFile("./private_acc.pem")
	if err2 != nil {
		log.Print(err2)
		statusMessage = "BadRequest"
	}
	test, err3 := hex.DecodeString(ciphertext)
	if err3 != nil {
		log.Print(err3)
		statusMessage = "BadRequest"
	}
	password := string(encryption.Decrypt("rsa", test, prvKey))

	if acc.Username == username && acc.Password == password {
		permission = true
	} else {
		permission = false
	}
	return permission, statusMessage
}


//Install method is called when first open the admin server.
func Install() models.InitStatus {
	o := orm.NewOrm()
	status := models.InitStatusInfo{Id: 1}
	err := o.Read(&status)

	if err == orm.ErrNoRows {
    	fmt.Println("not found")
	} else if err == orm.ErrMissPK {
    	fmt.Println("pk missing")
	} 
	result := status.Status

	return result
}

//CreateUUID creates a file with an UUID in it.
func CreateUUID() string {
	var statusMessage string = "OK"

	u, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		statusMessage = "BadRequest"
	}

	folderPath := path.Join(os.Getenv("GOPATH"), "/secrets")
    if _, err := os.Stat(folderPath); os.IsNotExist(err) {
        os.Mkdir(folderPath, os.ModePerm) 
        os.Chmod(folderPath, os.ModePerm)
	}

	uuidPath := path.Join(os.Getenv("GOPATH"), "/secrets/initialAdminPassword")
	if _, err = os.Stat(uuidPath); os.IsNotExist(err) {
		f, err := os.Create(uuidPath)
		if err != nil {
			log.Print(err)
			statusMessage = "BadRequest"
		}
		f.WriteString(u.String())
		defer f.Close()
	}
	
	return statusMessage
}

//ValidateUUID compares input with the UUID stored in the specified file.
func ValidateUUID(input string) (a bool, b string) {
	var statusMessage string = "OK"

	uuidPath := path.Join(os.Getenv("GOPATH"), "/secrets/initialAdminPassword")
	f, err := ioutil.ReadFile(uuidPath)
	if err != nil {
		log.Print(err)
		statusMessage = "BadRequest"
	}

	result := (input == string(f))
	if result == true {
		os.Remove(uuidPath)
		logs.Info("loading Board images...")
		err := Execute(models.LoadImagesFile)
		if err != nil {
			log.Print(err)
			statusMessage = "BadRequest"
		}

		o := orm.NewOrm()
		status := models.InitStatusInfo{Id: 1}
		err = o.Read(&status)
		if err == orm.ErrNoRows {
    		fmt.Println("not found")
		} else if err == orm.ErrMissPK {
    		fmt.Println("pk missing")
		} 
		if status.Status == models.InitStatusTrue {
			status.InstallTime = time.Now().Unix()
			status.Status = models.InitStatusFirst
			o.Update(&status, "InstallTime", "Status")
		}
	}

	return result, statusMessage
}
