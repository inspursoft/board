package service

import (
	"github.com/inspursoft/board/src/adminserver/encryption"
	"github.com/inspursoft/board/src/adminserver/models"
	"os"
	"path"

	"github.com/alyu/configparser"
)

//GetAllCfg returns the original data read from cfg file.
func GetAllCfg(which string, show bool) (*models.Configuration, error) {
	var cfgPath string
	configparser.Delimiter = "="
	if which == "" {
		cfgPath = path.Join("/go", "/cfgfile/board.cfg")
	} else {
		cfgPath = path.Join("/go", "/cfgfile/board.cfg.tmp")
	}

	//use configparser to read indicated cfg file.
	config, err := configparser.Read(cfgPath)
	if err != nil {
		return nil, err
	}
	//section sensitive, global refers to all sections.
	section, err := config.Section("global")
	if err != nil {
		return nil, err
	}

	//assigning values for each properties.
	cfg := models.GetConfiguration(section)

	if !show {
		cfg.Db.BoardAdminPassword = ""
		cfg.Db.Password = ""
		cfg.Jenkins.NodePassword = ""
		cfg.Email.Password = ""
		cfg.Es.Password = ""
		cfg.Gitlab.SSHPassword = ""
	}

	backupPath := path.Join("/go", "/cfgfile/board.cfg.bak1")
	cfg.FirstTimePost = !encryption.CheckFileIsExist(backupPath)
	tmpPath := path.Join("/go", "/cfgfile/board.cfg.tmp")
	cfg.TmpExist = encryption.CheckFileIsExist(tmpPath)
	if which == "" {
		cfg.Current = "cfg"
	} else {
		cfg.Current = "tmp"
	}

	return &cfg, nil
}

//UpdateCfg returns updated struct of data and set values for the cfg file.
func UpdateCfg(cfg *models.Configuration) error {
	configparser.Delimiter = "="
	cfgPath := path.Join("/go", "/cfgfile/board.cfg")
	//use configparser to read indicated cfg file.
	config, err := configparser.Read(cfgPath)
	if err != nil {
		return err
	}
	//section sensitive, global refers to all sections.
	section, err := config.Section("global")
	if err != nil {
		return err
	}

	if cfg.Email.Identity == "" {
		cfg.Email.Identity = " "
	}

	//setting value for each properties.
	models.UpdateConfiguration(section, cfg)

	//save the data from cache to file.
	err = configparser.Save(config, cfgPath)
	if err != nil {
		return err
	}

	err = os.Rename(cfgPath, cfgPath+".tmp")
	if err != nil {
		if !os.IsNotExist(err) { // fine if the file does not exists
			return err
		}
	}
	err = os.Rename(cfgPath+".bak", cfgPath)
	if err != nil {
		if !os.IsNotExist(err) { // fine if the file does not exists
			return err
		}
	}

	return nil
}
