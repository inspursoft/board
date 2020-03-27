package service

import (
	"git/inspursoft/board/src/adminserver/encryption"
	"git/inspursoft/board/src/adminserver/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"strings"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/alyu/configparser"
	
	"encoding/base64"
)

//GetAllCfg returns the original data read from cfg file.
func GetAllCfg(which string) (a *models.Configuration, b string) {
	//Cfg refers to an instance of configuration file.
	var Cfg *models.Configuration
	var statusMessage string = "OK"
	var cfgPath string

	configparser.Delimiter = "="

	if which == "" {
		cfgPath = path.Join("/go", "/cfgfile/board.cfg")
	} else {
		cfgPath = path.Join("/go", "/cfgfile/board.cfg.tmp")
	}

	//use configparser to read indicated cfg file.
	config, err0 := configparser.Read(cfgPath)
	if err0 != nil {
		log.Print(err0)
		statusMessage = "BadRequest"
	}
	//section sensitive, global refers to all sections.
	section, err1 := config.Section("global")
	if err1 != nil {
		log.Print(err1)
		statusMessage = "BadRequest"
	}

	//assigning values for each properties.
	Cfgi := models.GetConfiguration(section)

	Cfgi.Other.BoardAdminPassword = ""
	Cfgi.Other.DBPassword = ""
	Cfgi.Jenkinsserver.NodePassword = ""
	Cfgi.Email.Password = ""


	backupPath := path.Join("/go", "/cfgfile/board.cfg.bak1")
	Cfgi.FirstTimePost = !encryption.CheckFileIsExist(backupPath)

	tmpPath := path.Join("/go", "/cfgfile/board.cfg.tmp")
	Cfgi.TmpExist = encryption.CheckFileIsExist(tmpPath)

	if which == "" {
		Cfgi.Current = "cfg"
	} else {
		Cfgi.Current = "tmp"
	}

	//getting the address of the struct and return it with status message.
	Cfg = &Cfgi

	return Cfg, statusMessage
}

//UpdateCfg returns updated struct of data and set values for the cfg file.
func UpdateCfg(cfg *models.Configuration) string {
	var statusMessage string = "OK"
	configparser.Delimiter = "="
	cfgPath := path.Join("/go", "/cfgfile/board.cfg")
	//use configparser to read indicated cfg file.
	config, err0 := configparser.Read(cfgPath)
	if err0 != nil {
		log.Print(err0)
		statusMessage = "BadRequest"
	}
	//section sensitive, global refers to all sections.
	section, err1 := config.Section("global")
	if err1 != nil {
		log.Print(err1)
		statusMessage = "BadRequest"
	}

	//ENCRYPTION
	//existingPassword := section.ValueOf("board_admin_password")
	if cfg.Other.BoardAdminPassword != "" {
		prvKey, _ := ioutil.ReadFile("./private.pem")
		test, _ := base64.StdEncoding.DecodeString(cfg.Other.BoardAdminPassword)
		cfg.Other.BoardAdminPassword = string(encryption.Decrypt("rsa", test, prvKey))
	} else {
		o := orm.NewOrm()
		o.Using("mysql-db2")
		account := models.Account{Id: 1}
		err := o.Read(&account)
		if err == orm.ErrNoRows {
			fmt.Println("not found")
		} else if err == orm.ErrMissPK {
			fmt.Println("pk missing")
		} 
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
		cfg.Other.BoardAdminPassword = password
	}

	b, err := ioutil.ReadFile(path.Join(models.DBconfigdir, "/env"))
	if err != nil {
		logs.Error("error occurred on get DB env: %+v", err)
		panic(err)
	}
	DBpassword := strings.TrimPrefix(string(b), "DB_PASSWORD=")
	DBpassword = strings.Replace(DBpassword, "\n", "", 1)
	cfg.Other.DBPassword = DBpassword

	//setting value for each properties.
	models.UpdateConfiguration(section, cfg)

	//save the data from cache to file.
	err2 := configparser.Save(config, cfgPath)
	if err2 != nil {
		log.Print(err2)
		statusMessage = "BadRequest"
	}

	err = os.Rename(cfgPath, cfgPath+".tmp")
	if err != nil {
		if !os.IsNotExist(err) { // fine if the file does not exists
			log.Print(err)
			statusMessage = "BadRequest"
		}
	}
	err = os.Rename(cfgPath+".bak", cfgPath)
	if err != nil {
		if !os.IsNotExist(err) { // fine if the file does not exists
			log.Print(err)
			statusMessage = "BadRequest"
		}
	}

	return statusMessage
}

//GetKey generates 2 keys and return the public one.
func GetKey() string {
	_, pubKey := encryption.GenKey("rsa")
	//ciphertext := encryption.Encrypt("rsa", []byte("123456a?"), pubKey)
	//fmt.Println("###ciphertext:", base64.StdEncoding.EncodeToString(ciphertext))
	return string(pubKey)
}
