package service

import (
	"log"
	"os"
	"os/exec"
	"path"

)

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

	cfgPath := path.Join(os.Getenv("GOPATH"), "/cfgfile/board.cfg")
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