package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"time"
)

const (
	expircyTimeSpan       time.Duration = 900
	selectProject                       = "SELECT_PROJECT"
	selectImageList                     = "SELECT_IMAGES"
	createImage                         = "CREATE_IMAGE"
	configContainerList                 = "CONFIG_CONTAINERS"
	configExternalService               = "EXTERNAL_SERVICE"
)

type ConfigServiceStep struct {
	ProjectID           int64
	ServiceID           int64
	ImageList           []model.ImageIndex
	Dockerfile          model.Dockerfile
	ServiceName         string
	Instance            int
	ContainerList       []model.Container
	ExternalServiceList []model.ExternalService
}

func NewConfigServiceStep(key string) *ConfigServiceStep {
	configServiceStep := GetConfigServiceStep(key)
	if configServiceStep == nil {
		return &ConfigServiceStep{}
	}
	return configServiceStep
}

func SetConfigServiceStep(key string, s *ConfigServiceStep) {
	memoryCache.Put(key, s, time.Second*expircyTimeSpan)
}

func GetConfigServiceStep(key string) *ConfigServiceStep {
	if s, ok := memoryCache.Get(key).(*ConfigServiceStep); ok {
		return s
	}
	return nil
}

func (s *ConfigServiceStep) SelectProject(projectID int64) *ConfigServiceStep {
	s.ProjectID = projectID
	return s
}

func (s *ConfigServiceStep) GetSelectedProject() interface{} {
	return struct {
		ProjectID int64
	}{
		ProjectID: s.ProjectID,
	}
}

func (s *ConfigServiceStep) SelectImageList(imageList []model.ImageIndex) *ConfigServiceStep {
	s.ImageList = imageList
	return s
}

func (s *ConfigServiceStep) GetSelectedImageList() interface{} {
	return struct {
		ImageList []model.ImageIndex
	}{
		ImageList: s.ImageList,
	}
}

func (s *ConfigServiceStep) CreateDockerImage(dockerfile model.Dockerfile) *ConfigServiceStep {
	s.Dockerfile = dockerfile
	return s
}

func (s *ConfigServiceStep) GetDockerImage() interface{} {
	return struct {
		Dockerfile model.Dockerfile
	}{
		Dockerfile: s.Dockerfile,
	}
}

func (s *ConfigServiceStep) ConfigContainerList(containerList []model.Container) *ConfigServiceStep {
	s.ContainerList = containerList
	return s
}

func (s *ConfigServiceStep) GetConfigContainerList() interface{} {
	return struct {
		ContainerList []model.Container
	}{
		ContainerList: s.ContainerList,
	}
}

func (s *ConfigServiceStep) ConfigExternalService(serviceName string, instance int, externalServiceList []model.ExternalService) *ConfigServiceStep {
	s.ServiceName = serviceName
	s.Instance = instance
	s.ExternalServiceList = externalServiceList
	return s
}

func (s *ConfigServiceStep) GetConfigExternalService() interface{} {
	return struct {
		ServiceName         string
		Instance            int
		ExternalServiceList []model.ExternalService
	}{
		ServiceName:         s.ServiceName,
		Instance:            s.Instance,
		ExternalServiceList: s.ExternalServiceList,
	}
}

type ServiceConfigController struct {
	baseController
}

func (sc *ServiceConfigController) Prepare() {
	user := sc.getCurrentUser()
	if user == nil {
		sc.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	sc.currentUser = user
	sc.isSysAdmin = (user.SystemAdmin == 1)
}

func (sc *ServiceConfigController) GetConfigServiceStepAction() {
	key := sc.token
	configServiceStep := GetConfigServiceStep(key)
	if configServiceStep == nil {
		sc.customAbort(http.StatusNotFound, "Config service step has not been created yet.")
		return
	}
	phase := sc.GetString("phase")
	var result interface{}
	switch phase {
	case selectProject:
		result = configServiceStep.GetSelectedProject()
	case selectImageList:
		result = configServiceStep.GetSelectedImageList()
	case createImage:
		result = configServiceStep.GetConfigContainerList()
	case configContainerList:
		result = configServiceStep.GetConfigContainerList()
	case configExternalService:
		result = configServiceStep.GetConfigExternalService()
	}
	sc.Data["json"] = result
	sc.ServeJSON()
}

func (sc *ServiceConfigController) SetConfigServiceStepAction() {
	phase := sc.GetString("phase")
	key := sc.token
	configServiceStep := NewConfigServiceStep(key)
	reqData, err := sc.resolveBody()
	if err != nil {
		sc.internalError(err)
		return
	}
	switch phase {
	case selectProject:
		sc.selectProject(key, configServiceStep)
	case selectImageList:
		sc.selectImageList(key, configServiceStep, reqData)
	case createImage:
		sc.createImage(key, configServiceStep, reqData)
	case configContainerList:
		sc.configContainerList(key, configServiceStep, reqData)
	case configExternalService:
		sc.configExternalService(key, configServiceStep, reqData)
	}
}

func (sc *ServiceConfigController) selectProject(key string, configServiceStep *ConfigServiceStep) {
	projectID, err := sc.GetInt64("project_id")
	if err != nil {
		sc.internalError(err)
		return
	}
	SetConfigServiceStep(key, configServiceStep.SelectProject(projectID))
}

func (sc *ServiceConfigController) selectImageList(key string, configServiceStep *ConfigServiceStep, reqData []byte) {
	var imageList []model.ImageIndex
	err := json.Unmarshal(reqData, &imageList)
	if err != nil {
		sc.internalError(err)
		return
	}
	SetConfigServiceStep(key, configServiceStep.SelectImageList(imageList))
}

func (sc *ServiceConfigController) createImage(key string, configServiceStep *ConfigServiceStep, reqData []byte) {
	var dockerfile model.Dockerfile
	err := json.Unmarshal(reqData, &dockerfile)
	if err != nil {
		sc.internalError(err)
		return
	}
	SetConfigServiceStep(key, configServiceStep.CreateDockerImage(dockerfile))
}

func (sc *ServiceConfigController) configContainerList(key string, configServiceStep *ConfigServiceStep, reqData []byte) {
	var containerList []model.Container
	err := json.Unmarshal(reqData, &containerList)
	if err != nil {
		sc.internalError(err)
		return
	}
	SetConfigServiceStep(key, configServiceStep.ConfigContainerList(containerList))
}

func (sc *ServiceConfigController) configExternalService(key string, configServiceStep *ConfigServiceStep, reqData []byte) {
	serviceName := sc.GetString("service_name")
	instance, err := sc.GetInt("instance")
	if err != nil {
		sc.internalError(err)
		return
	}
	var externalServiceList []model.ExternalService
	err = json.Unmarshal(reqData, &externalServiceList)
	if err != nil {
		sc.internalError(err)
		return
	}
	SetConfigServiceStep(key, configServiceStep.ConfigExternalService(serviceName, instance, externalServiceList))
}
