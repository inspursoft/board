package service

import (
	"git/inspursoft/board/src/adminserver/encryption"
	"git/inspursoft/board/src/adminserver/models"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/alyu/configparser"
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
	f, err := os.Create("acc.txt")
	if err != nil {
		log.Print(err)
		statusMessage = "BadRequest"
	}
	f.WriteString("username = " + acc.Username + "\n")
	f.WriteString("password = " + acc.Password + "\n")
	f.Close()
	return statusMessage
}

//Login allow user to use account information to login adminserver.
func Login(acc *models.Account) (a bool, b string) {
	var statusMessage string = "OK"
	var permission bool
	configparser.Delimiter = "="
	config, err0 := configparser.Read("./acc.txt")
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
	password := section.ValueOf("password")
	if acc.Username == username && acc.Password == password {
		permission = true
	} else {
		permission = false
	}
	return permission, statusMessage
}

//Restart Board without loading cfg.
func Restart(path string) string {
	var statusMessage string = "OK"

	statusMessage = Shutdown(path)

	err := Execute("docker-compose -f " + path + "/docker-compose.yml up -d")
	if err != nil {
		log.Println(err)
		statusMessage = "BadRequest"
	}

	return statusMessage
}

//Applycfg restarts Board with applying of cfg.
func Applycfg(cfgpath string) string {
	var statusMessage string = "OK"

	cfgPath := path.Join("/go", "/cfgfile/board.cfg")
	err := os.Rename(cfgPath, cfgPath+".bak1")
	if err != nil {
		if !os.IsNotExist(err) { // fine if the file does not exists
			log.Print(err)
			statusMessage = "BadRequest"
		}
	}
	err = os.Rename(cfgPath+".tmp", cfgPath)
	if err != nil {
		if !os.IsNotExist(err) { // fine if the file does not exists
			log.Print(err)
			statusMessage = "BadRequest"
		}
	}

	statusMessage = Shutdown(cfgpath)

	err = Execute(cfgpath + "/prepare")
	if err != nil {
		log.Println(err)
		statusMessage = "BadRequest"
	}

	err = Execute("docker-compose -f " + cfgpath + "/docker-compose.yml up -d")
	if err != nil {
		log.Println(err)
		statusMessage = "BadRequest"
	}

	return statusMessage
}

//Shutdown Board.
func Shutdown(path string) string {
	var statusMessage string = "OK"

	err := Execute("docker-compose -f " + path + "/docker-compose.yml down")
	if err != nil {
		log.Println(err)
		statusMessage = "BadRequest"
	}

	return statusMessage
}

//Execute command.
func Execute(command string) error {
	cmd := exec.Command("sh", "-c", command)
	bytes, err := cmd.Output()
	log.Println(string(bytes))
	return err
}

//Install method is called when first open the admin server.
func Install() bool {
	return encryption.CheckFileIsExist("./install")
}
