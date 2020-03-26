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
	"github.com/alyu/configparser"
	uuid "github.com/satori/go.uuid"
	"fmt"
	"time"
)

//VerifyPassword compares the password in cfg with the input one.
func VerifyPassword(passwd *models.Password) (a bool, err string) {
	var statusMessage string = "OK"

	configparser.Delimiter = "="
	cfgPath := path.Join("/go", "/cfgfile/board.cfg")
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

	o := orm.NewOrm()
	o.Using("mysql-db2")
	account := models.Account{Username: acc.Username, Password: hex.EncodeToString(ciphertext)}

	if o.Read(&models.Account{Id: 1}) == orm.ErrNoRows {
		if _, err := o.Insert(&account); err != nil {
			log.Print(err)
			statusMessage = "BadRequest"
		}	
	} else {
		if _, err := o.Update(&account); err != nil {
			log.Print(err)
			statusMessage = "BadRequest"
		}	
	}

	if statusMessage == "OK" {
		o2 := orm.NewOrm()
		o2.Using("default")
		status := models.InitStatusInfo{Id: 1}
		err := o2.Read(&status)
		if err == orm.ErrNoRows {
			fmt.Println("not found")
		} else if err == orm.ErrMissPK {
			fmt.Println("pk missing")
		} 
		if status.Status == models.InitStatusThird {
			status.InstallTime = time.Now().Unix()
			status.Status = models.InitStatusFalse
			o2.Update(&status, "InstallTime", "Status")
		}
	}

	return statusMessage
}

//Login allow user to use account information to login adminserver.
func Login(acc *models.Account) (bool, string, string) {
	var statusMessage string = "OK"
	var permission bool

	o := orm.NewOrm()
	o.Using("mysql-db2")
	account := models.Account{Id: 1}
	err := o.Read(&account)
	if err == orm.ErrNoRows {
		fmt.Println("not found")
	} else if err == orm.ErrMissPK {
		fmt.Println("pk missing")
	} 
	username := account.Username
	ciphertext := account.Password

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

	var token string = ""
	if permission == true {
		u, err := uuid.NewV4()
		if err != nil {
			log.Println(err)
			statusMessage = "BadRequest"
		}
		token = u.String()

		newtoken := models.Token{Id: 1, Token: token, Time: time.Now().Unix()}
		o2 := orm.NewOrm()
		o2.Using("default")
		if o2.Read(&models.Token{Id: 1}) == orm.ErrNoRows {
			if _, err := o2.Insert(&newtoken); err != nil {
				log.Print(err)
				statusMessage = "BadRequest"
			}	
		} else {
			if _, err := o2.Update(&newtoken); err != nil {
				log.Print(err)
				statusMessage = "BadRequest"
			}	
		}
	}
	return permission, statusMessage, token
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

	return status.Status
}

//CreateUUID creates a file with an UUID in it.
func CreateUUID() string {
	var statusMessage string = "OK"

	u, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		statusMessage = "BadRequest"
	}

	folderPath := path.Join("/go", "/secrets")
    if _, err := os.Stat(folderPath); os.IsNotExist(err) {
        os.Mkdir(folderPath, os.ModePerm) 
        os.Chmod(folderPath, os.ModePerm)
	}

	uuidPath := path.Join("/go", "/secrets/initialAdminPassword")
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

	uuidPath := path.Join("/go", "/secrets/initialAdminPassword")
	f, err := ioutil.ReadFile(uuidPath)
	if err != nil {
		log.Print(err)
		statusMessage = "BadRequest"
	}

	result := (input == string(f))
	if result == true {
		os.Remove(uuidPath)

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


func VerifyToken(input string) bool {
	o := orm.NewOrm()
	token := models.Token{Id: 1}
	err := o.Read(&token)
	if err == orm.ErrNoRows {
		fmt.Println("not found")
		return false
	} else if err == orm.ErrMissPK {
		fmt.Println("pk missing")
		return false
	} 

	if input == token.Token && (time.Now().Unix()-token.Time)<=1800 {
		return true
	} else {
		return false
	}
}
