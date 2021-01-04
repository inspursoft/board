package controller

import (
	"encoding/json"
	"errors"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

const (
	expircyTimeSpan         time.Duration = 900
	selectProject                         = "SELECT_PROJECT"
	configContainerList                   = "CONFIG_CONTAINERS"
	configInitContainerList               = "CONFIG_INIT_CONTAINERS"
	configExternalService                 = "EXTERNAL_SERVICE"
	configEntireService                   = "ENTIRE_SERVICE"
	maximumPortNum                        = 32765
	minimumPortNum                        = 30000
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
	nodeOrNodeGroupNameNotFound        = errors.New("ERR_NODE_SELECTOR_NAME_NOT_FOUND")
	resourcerequestErr                 = errors.New("ERR_INVALID_RESOURCE_REQUEST")
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
	c.MemoryCache.Put(key, s, time.Second*expircyTimeSpan)
}

func GetConfigServiceStep(key string) *ConfigServiceStep {
	if s, ok := c.MemoryCache.Get(key).(*ConfigServiceStep); ok {
		return s
	}
	return nil
}

func DeleteConfigServiceStep(key string) error {
	if c.MemoryCache.IsExist(key) {
		return c.MemoryCache.Delete(key)
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

func (s *ConfigServiceStep) ConfigContainerList(containerList []model.Container) *ConfigServiceStep {
	s.ContainerList = containerList
	return s
}

func (s *ConfigServiceStep) ConfigInitContainerList(containerList []model.Container) *ConfigServiceStep {
	s.InitContainerList = containerList
	return s
}

func (s *ConfigServiceStep) GetConfigContainerList() interface{} {
	return struct {
		ProjectID     int64             `json:"project_id"`
		ProjectName   string            `json:"project_name"`
		ContainerList []model.Container `json:"container_list"`
	}{
		ProjectID:     s.ProjectID,
		ProjectName:   s.ProjectName,
		ContainerList: s.ContainerList,
	}
}

func (s *ConfigServiceStep) GetConfigInitContainerList() interface{} {
	return struct {
		ProjectID         int64             `json:"project_id"`
		ProjectName       string            `json:"project_name"`
		InitContainerList []model.Container `json:"initcontainer_list"`
	}{
		ProjectID:         s.ProjectID,
		ProjectName:       s.ProjectName,
		InitContainerList: s.InitContainerList,
	}
}

func (s *ConfigServiceStep) configExternalService(serviceName string, serviceType int, clusterIP string, instance int, public int, nodeOrNodeGroupName string, externalServiceList []model.ExternalService, sessionAffinityFlag int, sessionAffinityTime int) *ConfigServiceStep {
	s.ServiceName = serviceName
	s.ServiceType = serviceType
	s.Instance = instance
	s.Public = public
	s.ClusterIP = clusterIP
	s.NodeSelector = nodeOrNodeGroupName
	s.ExternalServiceList = externalServiceList
	s.SessionAffinityFlag = sessionAffinityFlag
	s.SessionAffinityTime = sessionAffinityTime
	return s
}

func (s *ConfigServiceStep) configAffinity(affinityList []model.Affinity) *ConfigServiceStep {
	if affinityList != nil {
		s.AffinityList = affinityList
	}
	return s
}

func (s *ConfigServiceStep) GetConfigExternalService() interface{} {
	return struct {
		ProjectName         string                  `json:"project_name"`
		ServiceName         string                  `json:"service_name"`
		ServiceType         int                     `json:"service_type"`
		Instance            int                     `json:"instance"`
		Public              int                     `json:"service_public"`
		ClusterIP           string                  `json:"cluster_ip"`
		NodeSelector        string                  `json:"node_selector"`
		ExternalServiceList []model.ExternalService `json:"external_service_list"`
		AffinityList        []model.Affinity        `json:"affinity_list"`
		SessionAffinityFlag int                     `json:"session_affinity_flag"`
		SessionAffinityTime int                     `json:"session_affinity_time"`
	}{
		ProjectName:         s.ProjectName,
		ServiceName:         s.ServiceName,
		ServiceType:         s.ServiceType,
		Instance:            s.Instance,
		Public:              s.Public,
		ClusterIP:           s.ClusterIP,
		NodeSelector:        s.NodeSelector,
		ExternalServiceList: s.ExternalServiceList,
		AffinityList:        s.AffinityList,
		SessionAffinityFlag: s.SessionAffinityFlag,
		SessionAffinityTime: s.SessionAffinityTime,
	}
}

type ServiceConfigController struct {
	c.BaseController
}

func (sc *ServiceConfigController) getKey() string {
	return strconv.Itoa(int(sc.CurrentUser.ID))
}

func (sc *ServiceConfigController) GetConfigServiceStepAction() {
	key := sc.getKey()
	configServiceStep := GetConfigServiceStep(key)
	if configServiceStep == nil {
		sc.CustomAbortAudit(http.StatusNotFound, "Config service step has not been created yet.")
		return
	}
	phase := sc.GetString("phase")
	var result interface{}
	switch phase {
	case selectProject:
		result = configServiceStep.GetSelectedProject()
	case configContainerList:
		result = configServiceStep.GetConfigContainerList()
	case configInitContainerList:
		result = configServiceStep.GetConfigInitContainerList()
	case configExternalService:
		result = configServiceStep.GetConfigExternalService()
	case configEntireService:
		result = configServiceStep
	default:
		sc.ServeStatus(http.StatusBadRequest, phaseInvalidErr.Error())
		return
	}

	if err, ok := result.(error); ok {
		if err == projectIDInvalidErr {
			sc.ServeStatus(http.StatusBadRequest, err.Error())
			return
		}
		sc.InternalError(err)
		return
	}
	sc.RenderJSON(result)
}

func (sc *ServiceConfigController) SetConfigServiceStepAction() {
	phase := sc.GetString("phase")
	key := sc.getKey()
	configServiceStep := NewConfigServiceStep(key)
	reqData, err := ioutil.ReadAll(sc.Ctx.Request.Body)
	if err != nil {
		sc.InternalError(err)
		return
	}
	switch phase {
	case selectProject:
		sc.selectProject(key, configServiceStep)
	case configContainerList:
		sc.configContainerList(key, configServiceStep, configServiceStep.ConfigContainerList, reqData)
	case configInitContainerList:
		sc.configContainerList(key, configServiceStep, configServiceStep.ConfigInitContainerList, reqData)
	case configExternalService:
		sc.configExternalService(key, configServiceStep, reqData)
	case configEntireService:
		sc.configEntireService(key, configServiceStep, reqData)
	default:
		sc.ServeStatus(http.StatusBadRequest, phaseInvalidErr.Error())
		return
	}
	logs.Info("set configService after phase %s is :%+v", phase, NewConfigServiceStep(key))
}

func (sc *ServiceConfigController) DeleteServiceStepAction() {
	key := sc.getKey()
	err := DeleteConfigServiceStep(key)
	if err != nil {
		sc.InternalError(err)
		return
	}
}
func (sc *ServiceConfigController) selectProject(key string, configServiceStep *ConfigServiceStep) {
	projectID, err := sc.GetInt64("project_id")
	if err != nil {
		sc.InternalError(err)
		return
	}

	project := sc.ResolveUserPrivilegeByID(int64(projectID))
	SetConfigServiceStep(key, configServiceStep.SelectProject(projectID, project.Name))
}

func (sc *ServiceConfigController) configContainerList(key string, configServiceStep *ConfigServiceStep, f func([]model.Container) *ConfigServiceStep, reqData []byte) {
	var containerList []model.Container
	err := json.Unmarshal(reqData, &containerList)
	if err != nil {
		sc.InternalError(err)
		return
	}

	//Check CPU Mem request and limit
	for _, container := range containerList {
		if container.CPURequest != "" && container.CPULimit != "" {
			cpurequest, err := strconv.Atoi(strings.TrimRight(container.CPURequest, "m"))
			if err != nil {
				sc.ServeStatus(http.StatusBadRequest, resourcerequestErr.Error())
				return
			}
			cpulimit, err := strconv.Atoi(strings.TrimRight(container.CPULimit, "m"))
			if err != nil {
				sc.ServeStatus(http.StatusBadRequest, resourcerequestErr.Error())
				return
			}
			if cpurequest > cpulimit {
				sc.ServeStatus(http.StatusBadRequest, resourcerequestErr.Error())
				return
			}

		}
		if container.MemRequest != "" && container.MemLimit != "" {
			memrequest, err := strconv.Atoi(strings.TrimRight(container.MemRequest, "Mi"))
			if err != nil {
				sc.ServeStatus(http.StatusBadRequest, resourcerequestErr.Error())
				return
			}
			memlimit, err := strconv.Atoi(strings.TrimRight(container.MemLimit, "Mi"))
			if err != nil {
				sc.ServeStatus(http.StatusBadRequest, resourcerequestErr.Error())
				return
			}
			if memrequest > memlimit {
				sc.ServeStatus(http.StatusBadRequest, resourcerequestErr.Error())
				return
			}

		}
	}

	configServiceStep.ExternalServiceList = service.CheckServiceConfigPortMap(configServiceStep.ExternalServiceList, containerList)
	SetConfigServiceStep(key, f(containerList))
}

func (sc *ServiceConfigController) configExternalService(key string, configServiceStep *ConfigServiceStep, reqData []byte) {
	serviceName := strings.ToLower(sc.GetString("service_name"))
	if !utils.ValidateWithLengthRange(serviceName, 1, 63) {
		sc.ServeStatus(http.StatusBadRequest, "Service Name must be not empty and no more than 63 characters ")
		return
	}

	isDuplicate, err := sc.checkServiceDuplicateName(serviceName)
	if err != nil {
		sc.ServeStatus(http.StatusBadRequest, err.Error())
		return
	}
	if isDuplicate == true {
		sc.ServeStatus(http.StatusBadRequest, serverNameDuplicateErr.Error())
		return
	}

	instance, err := sc.GetInt("instance")
	if err != nil {
		sc.InternalError(err)
		return
	}
	if instance < 1 {
		sc.ServeStatus(http.StatusBadRequest, instanceInvalidErr.Error())
		return
	}

	public, err := sc.GetInt("service_public")
	if err != nil {
		sc.InternalError(err)
		return
	}

	serviceType, err := sc.GetInt("service_type")
	if err != nil {
		sc.InternalError(err)
		return
	}

	clusterIP := sc.GetString("cluster_ip")
	// TODO check valid cluster IP

	// A tmp fix for headless service
	if serviceType == model.ServiceTypeStatefulSet {
		clusterIP = "None"
	}

	sessionAffinityFlag, err := sc.GetInt("session_affinity_flag", 0)
	if err != nil {
		sc.InternalError(err)
		return
	}

	sessionAffinityTime, err := sc.GetInt("session_affinity_time", 0)
	if err != nil {
		sc.InternalError(err)
		return
	}

	nodeOrNodeGroupName := strings.ToLower(sc.GetString("node_selector"))
	if nodeOrNodeGroupName != "" {
		isExists, err := service.NodeOrNodeGroupExists(nodeOrNodeGroupName)
		if err != nil {
			sc.InternalError(err)
			return
		}
		if !isExists {
			sc.ServeStatus(http.StatusBadRequest, nodeOrNodeGroupNameNotFound.Error())
			return
		}
	}

	var serviceConfig model.ConfigServiceStep
	err = json.Unmarshal(reqData, &serviceConfig)
	if err != nil {
		sc.InternalError(err)
		return
	}

	for _, external := range serviceConfig.ExternalServiceList {
		if external.NodeConfig.NodePort != 0 && (external.NodeConfig.NodePort > maximumPortNum || external.NodeConfig.NodePort < minimumPortNum) {
			sc.ServeStatus(http.StatusBadRequest, portInvalidErr.Error())
			return
		}
	}
	configServiceStep.configExternalService(serviceName, serviceType, clusterIP, instance, public, nodeOrNodeGroupName, serviceConfig.ExternalServiceList, sessionAffinityFlag, sessionAffinityTime)
	SetConfigServiceStep(key, configServiceStep.configAffinity(serviceConfig.AffinityList))
}

func (sc *ServiceConfigController) checkServiceDuplicateName(serviceName string) (bool, error) {
	key := sc.getKey()
	configServiceStep := GetConfigServiceStep(key)
	if configServiceStep == nil {
		return false, serviceConfigNotCreateErr
	}

	project, err := service.GetProject(model.Project{ID: configServiceStep.ProjectID}, "id")
	if err != nil {
		sc.InternalError(err)
	}
	if project == nil {
		return false, serviceConfigNotSetProjectErr
	}

	isServiceDuplicated, err := service.ServiceExists(serviceName, project.Name)
	if err != nil {
		sc.InternalError(err)
		return false, err
	}
	return isServiceDuplicated, nil

}

func (sc *ServiceConfigController) checkEntireServiceConfig(entireService *ConfigServiceStep) error {
	project, err := service.GetProject(model.Project{ID: entireService.ProjectID}, "id")
	if err != nil {
		sc.InternalError(err)
		return err
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
		sc.InternalError(err)
	}
	if isDuplicate == true {
		return serverNameDuplicateErr
	}
	entireService.ProjectName = project.Name

	if entireService.Instance < 1 {
		return instanceInvalidErr
	}

	//TODO: Skip check and transfer to Lower case
	//	for key, container := range entireService.ContainerList {
	//		entireService.ContainerList[key].VolumeMounts.VolumeName = strings.ToLower(container.VolumeMounts.VolumeName)
	//		entireService.ContainerList[key].Name = strings.ToLower(container.Name)
	//	}

	if len(entireService.ExternalServiceList) < 1 {
		return emptyExternalServiceListErr
	}
	for _, external := range entireService.ExternalServiceList {
		if external.NodeConfig.NodePort != 0 && (external.NodeConfig.NodePort > maximumPortNum || external.NodeConfig.NodePort < minimumPortNum) {
			return portInvalidErr
		}
	}

	return nil
}

func (sc *ServiceConfigController) configEntireService(key string, configServiceStep *ConfigServiceStep, reqData []byte) {
	var entireService ConfigServiceStep
	err := json.Unmarshal(reqData, &entireService)
	if err != nil {
		sc.InternalError(err)
		return
	}

	if err = sc.checkEntireServiceConfig(&entireService); err != nil {
		sc.ServeStatus(http.StatusBadRequest, err.Error())
		return
	}

	SetConfigServiceStep(key, &entireService)
}
