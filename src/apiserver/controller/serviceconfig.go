package controller

import (
	"encoding/json"
	"errors"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

const (
	expircyTimeSpan       time.Duration = 900
	selectProject                       = "SELECT_PROJECT"
	selectImageList                     = "SELECT_IMAGES"
	configContainerList                 = "CONFIG_CONTAINERS"
	configExternalService               = "EXTERNAL_SERVICE"
	configEntireService                 = "ENTIRE_SERVICE"
	maximumPortNum                      = 32765
	minimumPortNum                      = 30000
)

var (
	serverNameDuplicateErr             = errors.New("ERR_DUPLICATE_SERVICE_NAME")
	projectIDInvalidErr                = errors.New("ERR_INVALID_PROJECT_ID")
	imageListInvalidErr                = errors.New("ERR_INVALID_IMAGE_LIST")
	portInvalidErr                     = errors.New("ERR_INVALID_SERVICE_NODEPORT")
	instanceInvalidErr                 = errors.New("ERR_INVALID_SERVICE_INSTANCE")
	emptyServiceNameErr                = errors.New("ERR_EMPTY_SERVICE_NAME")
	emptyVolumeTargetStorageServiceErr = errors.New("ERR_EMPTY_VOLUME_TARGET_STORAGE_SERVICE_ERR")
	phaseInvalidErr                    = errors.New("ERR_INVALID_PHASE")
	serviceConfigNotCreateErr          = errors.New("ERR_NOT_CREATE_SERVICE_CONFIG")
	serviceConfigNotSetProjectErr      = errors.New("ERR_NOT_SET_PROJECT_IN_SERVICE_CONFIG")
	emptyExternalServiceListErr        = errors.New("ERR_EMPTY_EXTERNAL_SERVICE_LIST")
	notFoundErr                        = errors.New("ERR_NOT_FOUND")
)

type ConfigServiceStep model.ConfigServiceStep

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

func DeleteConfigServiceStep(key string) error {
	if memoryCache.IsExist(key) {
		return memoryCache.Delete(key)
	}
	return nil
}

func (s *ConfigServiceStep) SelectProject(projectID int64, projectName string) *ConfigServiceStep {
	s.ProjectID = projectID
	s.ProjectName = projectName
	return s
}

func (s *ConfigServiceStep) GetSelectedProject() interface{} {
	return struct {
		ProjectID   int64  `json:"project_id"`
		ProjectName string `json:"project_name"`
	}{
		ProjectID:   s.ProjectID,
		ProjectName: s.ProjectName,
	}
}

func (s *ConfigServiceStep) SelectImageList(imageList []model.ImageIndex) *ConfigServiceStep {
	s.ImageList = imageList
	return s
}

func (s *ConfigServiceStep) GetSelectedImageList() interface{} {
	return struct {
		ProjectID   int64              `json:"project_id"`
		ProjectName string             `json:"project_name"`
		ImageList   []model.ImageIndex `json:"image_list"`
	}{
		ProjectID:   s.ProjectID,
		ProjectName: s.ProjectName,
		ImageList:   s.ImageList,
	}
}

func (s *ConfigServiceStep) ConfigContainerList(containerList []model.Container) *ConfigServiceStep {
	s.ContainerList = containerList
	return s
}

func (s *ConfigServiceStep) GetConfigContainerList() interface{} {
	if len(s.ContainerList) < 1 {
		for _, image := range s.ImageList {
			fromIndex := strings.LastIndex(image.ImageName, "/")
			image.ProjectName = image.ImageName[:fromIndex]
			s.ContainerList = append(s.ContainerList, model.Container{Name: image.ImageName[fromIndex+1:], Image: image})
		}
	} else {
		containerList := make([]model.Container, 0)
		for _, image := range s.ImageList {
			hasChanged := false
			for _, container := range s.ContainerList {
				if image.ImageName == container.Image.ImageName && image.ImageTag == container.Image.ImageTag {
					hasChanged = true
					containerList = append(containerList, container)
					break
				}
			}
			if hasChanged == false {
				fromIndex := strings.LastIndex(image.ImageName, "/")
				image.ProjectName = image.ImageName[:fromIndex]
				containerList = append(containerList, model.Container{Name: image.ImageName[fromIndex+1:], Image: image})
			}
		}
		s.ContainerList = containerList
	}

	return struct {
		ContainerList []model.Container `json:"container_list"`
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
		ProjectName         string                  `json:"project_name"`
		ServiceName         string                  `json:"service_name"`
		Instance            int                     `json:"instance"`
		Public              int                     `json:"service_public"`
		ExternalServiceList []model.ExternalService `json:"external_service_list"`
	}{
		ProjectName:         s.ProjectName,
		ServiceName:         s.ServiceName,
		Instance:            s.Instance,
		Public:              s.Public,
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

func (sc *ServiceConfigController) getKey() string {
	return strconv.Itoa(int(sc.currentUser.ID))
}

func (sc *ServiceConfigController) GetConfigServiceStepAction() {
	key := sc.getKey()
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
	case configContainerList:
		result = configServiceStep.GetConfigContainerList()
	case configExternalService:
		result = configServiceStep.GetConfigExternalService()
	case configEntireService:
		result = configServiceStep
	default:
		sc.serveStatus(http.StatusBadRequest, phaseInvalidErr.Error())
		return
	}

	if err, ok := result.(error); ok {
		if err == projectIDInvalidErr {
			sc.serveStatus(http.StatusBadRequest, err.Error())
			return
		}
		sc.internalError(err)
	}

	sc.Data["json"] = result
	sc.ServeJSON()
}

func (sc *ServiceConfigController) SetConfigServiceStepAction() {
	phase := sc.GetString("phase")
	key := sc.getKey()
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
	case configContainerList:
		sc.configContainerList(key, configServiceStep, reqData)
	case configExternalService:
		sc.configExternalService(key, configServiceStep, reqData)
	case configEntireService:
		sc.configEntireService(key, configServiceStep, reqData)
	default:
		sc.serveStatus(http.StatusBadRequest, phaseInvalidErr.Error())
		return
	}
	logs.Info("set configService after phase %s is :%+v", phase, NewConfigServiceStep(key))
}

func (sc *ServiceConfigController) DeleteServiceStepAction() {
	key := sc.getKey()
	err := DeleteConfigServiceStep(key)
	if err != nil {
		sc.internalError(err)
		return
	}
}
func (sc *ServiceConfigController) selectProject(key string, configServiceStep *ConfigServiceStep) {
	projectID, err := sc.GetInt64("project_id")
	if err != nil {
		sc.internalError(err)
		return
	}

	project, err := service.GetProject(model.Project{ID: projectID}, "id")
	if err != nil {
		sc.internalError(err)
		return
	}
	if project == nil {
		sc.serveStatus(http.StatusBadRequest, projectIDInvalidErr.Error())
		return
	}

	SetConfigServiceStep(key, configServiceStep.SelectProject(projectID, project.Name))
}

func (sc *ServiceConfigController) selectImageList(key string, configServiceStep *ConfigServiceStep, reqData []byte) {
	var imageList []model.ImageIndex
	err := json.Unmarshal(reqData, &imageList)
	if err != nil {
		sc.internalError(err)
		return
	}

	if len(imageList) < 0 {
		sc.serveStatus(http.StatusBadRequest, imageListInvalidErr.Error())
		return
	}
	for _, image := range imageList {
		if strings.Index(image.ImageName, "/") == -1 || len(strings.TrimSpace(image.ImageTag)) == 0 {
			sc.serveStatus(http.StatusBadRequest, imageListInvalidErr.Error())
			return
		}
	}

	SetConfigServiceStep(key, configServiceStep.SelectImageList(imageList))
}

func (sc *ServiceConfigController) configContainerList(key string, configServiceStep *ConfigServiceStep, reqData []byte) {
	var containerList []model.Container
	err := json.Unmarshal(reqData, &containerList)
	if err != nil {
		sc.internalError(err)
		return
	}

	for index, container := range containerList {
		if container.VolumeMounts.TargetPath != "" && container.VolumeMounts.TargetStorageService == "" {
			sc.serveStatus(http.StatusBadRequest, emptyVolumeTargetStorageServiceErr.Error())
			return
		}
		containerList[index].VolumeMounts.VolumeName = strings.ToLower(container.VolumeMounts.VolumeName)
		containerList[index].Name = strings.ToLower(container.Name)
	}

	SetConfigServiceStep(key, configServiceStep.ConfigContainerList(containerList))
}

func (sc *ServiceConfigController) configExternalService(key string, configServiceStep *ConfigServiceStep, reqData []byte) {
	serviceName := strings.ToLower(sc.GetString("service_name"))
	if serviceName == "" {
		sc.serveStatus(http.StatusBadRequest, emptyServiceNameErr.Error())
		return
	}

	isDuplicate, err := sc.checkServiceDuplicateName(serviceName)
	if err != nil {
		sc.serveStatus(http.StatusBadRequest, err.Error())
		return
	}
	if isDuplicate == true {
		sc.serveStatus(http.StatusBadRequest, serverNameDuplicateErr.Error())
		return
	}

	instance, err := sc.GetInt("instance")
	if err != nil {
		sc.internalError(err)
		return
	}
	if instance < 1 {
		sc.serveStatus(http.StatusBadRequest, instanceInvalidErr.Error())
		return
	}

	var externalServiceList []model.ExternalService
	err = json.Unmarshal(reqData, &externalServiceList)
	if err != nil {
		sc.internalError(err)
		return
	}

	for _, external := range externalServiceList {
		if external.NodeConfig.NodePort > maximumPortNum || external.NodeConfig.NodePort < minimumPortNum {
			sc.serveStatus(http.StatusBadRequest, portInvalidErr.Error())
			return
		}
	}

	SetConfigServiceStep(key, configServiceStep.ConfigExternalService(serviceName, instance, externalServiceList))
}

func (sc *ServiceConfigController) checkServiceDuplicateName(serviceName string) (bool, error) {
	key := sc.getKey()
	configServiceStep := GetConfigServiceStep(key)
	if configServiceStep == nil {
		return false, serviceConfigNotCreateErr
	}

	project, err := service.GetProject(model.Project{ID: configServiceStep.ProjectID}, "id")
	if err != nil {
		sc.internalError(err)
	}
	if project == nil {
		return false, serviceConfigNotSetProjectErr
	}

	isServiceDuplicated, err := service.ServiceExists(serviceName, project.Name)
	if err != nil {
		sc.internalError(err)
	}
	return isServiceDuplicated, nil

}

func (sc *ServiceConfigController) checkEntireServiceConfig(entireService *ConfigServiceStep) error {
	project, err := service.GetProject(model.Project{ID: entireService.ProjectID}, "id")
	if err != nil {
		sc.internalError(err)
	}
	if project == nil {
		return projectIDInvalidErr
	}

	serviceName := strings.ToLower(entireService.ServiceName)
	if serviceName == "" {
		return emptyServiceNameErr
	}
	isDuplicate, err := service.ServiceExists(serviceName, project.Name)
	if err != nil {
		sc.internalError(err)
	}
	if isDuplicate == true {
		return serverNameDuplicateErr
	}

	if entireService.Instance < 1 {
		return instanceInvalidErr
	}

	for key, container := range entireService.ContainerList {
		entireService.ContainerList[key].VolumeMounts.VolumeName = strings.ToLower(container.VolumeMounts.VolumeName)
		entireService.ContainerList[key].Name = strings.ToLower(container.Name)
	}

	if len(entireService.ExternalServiceList) < 1 {
		return emptyExternalServiceListErr
	}
	for _, external := range entireService.ExternalServiceList {
		if external.NodeConfig.NodePort > 32765 || external.NodeConfig.NodePort < 30000 {
			return portInvalidErr
		}
	}

	return nil
}

func (sc *ServiceConfigController) configEntireService(key string, configServiceStep *ConfigServiceStep, reqData []byte) {
	var entireService ConfigServiceStep
	err := json.Unmarshal(reqData, &entireService)
	if err != nil {
		sc.internalError(err)
		return
	}

	if err = sc.checkEntireServiceConfig(&entireService); err != nil {
		sc.serveStatus(http.StatusBadRequest, err.Error())
		return
	}

	SetConfigServiceStep(key, &entireService)
}

func (sc *ServiceConfigController) GetConfigServiceFromDBAction() {
	key := sc.getKey()
	configServiceStep := NewConfigServiceStep(key)
	serviceName := strings.ToLower(sc.GetString("service_name"))
	projectName := strings.ToLower(sc.GetString("project_name"))
	serviceData, err := service.GetService(model.ServiceStatus{Name: serviceName, ProjectName: projectName}, "name", "project_name")
	if err != nil {
		sc.internalError(err)
		return
	}
	if serviceData == nil || serviceData.ServiceConfig == "" {
		sc.serveStatus(http.StatusNotFound, notFoundErr.Error())
		return
	}

	logs.Info("service config form DB is %+v\n", serviceData)

	err = json.Unmarshal([]byte(serviceData.ServiceConfig), configServiceStep)
	if err != nil {
		sc.internalError(err)
		return
	}
	SetConfigServiceStep(key, configServiceStep)
}
