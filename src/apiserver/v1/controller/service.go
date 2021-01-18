package controller

import (
	"fmt"
	c "git/inspursoft/board/src/apiserver/controllers/commons"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/devops/travis"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

const (
	dockerfileName         = "Dockerfile"
	deploymentFilename     = "deployment.yaml"
	statefulsetFilename    = "statefulset.yaml"
	serviceFilename        = "service.yaml"
	rollingUpdateFilename  = "rollingUpdateDeployment.yaml"
	deploymentTestFilename = "testdeployment.yaml"
	serviceTestFilename    = "testservice.yaml"

	apiheader        = "Content-Type: application/yaml"
	deploymentAPI    = "/apis/extensions/v1beta1/namespaces/"
	serviceAPI       = "/api/v1/namespaces/"
	test             = "test"
	serviceNamespace = "default" //TODO create in project post
	k8sServices      = "kubernetes"
	deploymentType   = "deployment"
	serviceType      = "service"
	statefulsetType  = "statefulset"
	startingDuration = 300 * time.Second //300 seconds
)

const (
	preparing = iota
	running
	stopped
	uncompleted
	warning
	deploying
	completed
	failed
	unknown
	autonomousOffline
	partAutonomousOffline
)

var devOpsOpt = utils.GetConfig("DEVOPS_OPT")

type ServiceController struct {
	c.BaseController
}

func (p *ServiceController) generateDeploymentTravis(serviceName, deploymentURL, serviceURL string) error {
	userID := p.CurrentUser.ID
	var travisCommand travis.TravisCommand
	travisCommand.Script.Commands = []string{}
	items := []string{
		fmt.Sprintf("curl \"%s/jenkins-job/%d/$BUILD_NUMBER\"", c.BoardAPIBaseURL(), userID),
	}
	if deploymentURL != "" {
		items = append(items, fmt.Sprintf("#curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/deployment.yaml %s", serviceName, deploymentURL))
	}
	if serviceURL != "" {
		items = append(items, fmt.Sprintf("#curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/service.yaml %s", serviceName, serviceURL))
	}
	travisCommand.Script.Commands = items
	return travisCommand.GenerateCustomTravis(p.RepoPath)
}

func (p *ServiceController) getKey() string {
	return strconv.Itoa(int(p.CurrentUser.ID))
}

func (p *ServiceController) resolveServiceInfo() (s *model.ServiceStatus, err error) {
	serviceID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.InternalError(err)
		return
	}
	// Get the project info of this service
	s, err = service.GetServiceByID(int64(serviceID))
	if err != nil {
		p.InternalError(err)
		return
	}
	if s == nil {
		p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Invalid service ID: %d", serviceID))
		return
	}
	return
}

func (p *ServiceController) DeployServiceAction() {
	key := p.getKey()
	configService := NewConfigServiceStep(key)

	//Judge authority
	project := p.ResolveUserPrivilegeByID(configService.ProjectID)

	var newservice model.ServiceStatus
	newservice.Name = configService.ServiceName
	newservice.ProjectID = configService.ProjectID
	newservice.Status = preparing // 0: preparing 1: running 2: suspending
	newservice.OwnerID = p.CurrentUser.ID
	newservice.OwnerName = p.CurrentUser.Username
	newservice.Public = configService.Public
	newservice.ProjectName = project.Name
	newservice.Type = configService.ServiceType

	serviceInfo, err := service.CreateServiceConfig(newservice)
	if err != nil {
		p.InternalError(err)
		return
	}

	deployInfo, err := service.DeployService((*model.ConfigServiceStep)(configService), c.RegistryBaseURI())
	if err != nil {
		p.ParseError(err, c.ParsePostK8sError)
		_, deleteServiceError := service.DeleteService(serviceInfo.ID)
		if deleteServiceError != nil {
			logs.Error("Failed to delete the service data of %s in database. Error: %s", serviceInfo.Name, deleteServiceError.Error())
		}
		return
	}

	p.ResolveRepoServicePath(project.Name, newservice.Name)
	err = service.GenerateDeployYamlFiles(deployInfo, p.RepoServicePath)
	if err != nil {
		p.InternalError(err)
		return
	}

	deploymentFile := filepath.Join(newservice.Name, deploymentFilename)
	serviceFile := filepath.Join(newservice.Name, serviceFilename)
	items := []string{deploymentFile, serviceFile}
	p.PushItemsToRepo(items...)
	p.CollaborateWithPullRequest("master", "master", items...)
	p.MergeCollaborativePullRequest()

	updateService := model.ServiceStatus{ID: serviceInfo.ID, Status: uncompleted, ServiceYaml: string(deployInfo.ServiceFileInfo),
		DeploymentYaml: string(deployInfo.DeploymentFileInfo)}
	_, err = service.UpdateService(updateService, "status", "service_yaml", "deployment_yaml")
	if err != nil {
		p.InternalError(err)
		return
	}

	err = DeleteConfigServiceStep(key)
	if err != nil {
		p.InternalError(err)
		return
	}
	logs.Info("Service with ID:%d has been deleted in cache.", serviceInfo.ID)

	configService.ServiceID = serviceInfo.ID
	p.RenderJSON(configService)
}

func (p *ServiceController) DeployServiceTestAction() {
	key := p.getKey()
	configService := NewConfigServiceStep(key)
	configService.ServiceName = test + configService.ServiceName
	SetConfigServiceStep(key, configService)
	p.DeployServiceAction()
}

//
func syncK8sStatus(serviceList []*model.ServiceStatusMO) error {
	var err error
	// synchronize service status with the cluster system
	for _, serviceStatusMO := range serviceList {
		// Get serviceStatus from serviceStatusMO to adapt for updating services
		serviceStatus := &serviceStatusMO.ServiceStatus
		if (*serviceStatus).Status == stopped {
			continue
		}

		switch serviceStatus.Type {
		case model.ServiceTypeNormalNodePort, model.ServiceTypeClusterIP:
			// Check the deployment status
			deployment, _, _ := service.GetDeployment((*serviceStatus).ProjectName, (*serviceStatus).Name)
			if deployment == nil && serviceStatus.Name != k8sServices {
				logs.Info("Failed to get deployment", err)
				var reason = "The deployment is not established in cluster system"
				(*serviceStatus).Status = uncompleted
				// TODO create a new field in serviceStatus for reason
				(*serviceStatus).Comment = "Reason: " + reason
				_, err = service.UpdateService(*serviceStatus, "status", "Comment")
				if err != nil {
					logs.Error("Failed to update deployment.")
					break
				}
				continue
			} else {
				if deployment.Status.Replicas > deployment.Status.AvailableReplicas {
					logs.Debug("The desired replicas number is not available",
						deployment.Status.Replicas, deployment.Status.AvailableReplicas)
					(*serviceStatus).Status = uncompleted
					reason := "The desired replicas number is not available"
					(*serviceStatus).Comment = "Reason: " + reason
					_, err = service.UpdateService(*serviceStatus, "status", "Comment")
					if err != nil {
						logs.Error("Failed to update deployment replicas.")
						break
					}
					continue
				}
			}
		case model.ServiceTypeStatefulSet:
			// Check the statefulset status
			statefulset, _, _ := service.GetStatefulSet((*serviceStatus).ProjectName, (*serviceStatus).Name)
			if statefulset == nil || statefulset.Status.Replicas < *statefulset.Spec.Replicas {
				logs.Debug("The statefulset %s is not ready %v", (*serviceStatus).Name, err)
				(*serviceStatus).Status = uncompleted
			}
		case model.ServiceTypeEdgeComputing:
			// Check the deployment status
			deployment, _, _ := service.GetDeployment((*serviceStatus).ProjectName, (*serviceStatus).Name)
			if deployment == nil && serviceStatus.Name != k8sServices {
				logs.Info("Failed to get deployment", err)
				var reason = "The deployment is not established in cluster system"
				(*serviceStatus).Status = uncompleted
				// TODO create a new field in serviceStatus for reason
				(*serviceStatus).Comment = "Reason: " + reason
				_, err = service.UpdateService(*serviceStatus, "status", "Comment")
				if err != nil {
					logs.Error("Failed to update deployment.")
					break
				}
				continue
			} else if deployment.Status.Replicas > deployment.Status.AvailableReplicas {
				logs.Debug("The desired replicas number is not available",
					deployment.Status.Replicas, deployment.Status.AvailableReplicas)
				//Get the nodes which the pod were deployed.
				var status int
				var reason string
				podList, err := service.GetPodsByLabelSelector((*serviceStatus).ProjectName, &model.LabelSelector{
					MatchLabels: deployment.Spec.Selector})
				nodeList := service.GetNodeList()
				if err != nil {
					status = unknown
					reason = "Get the pod info failed."
				} else if nodeList == nil {
					status = unknown
					reason = "Get the node list failed."
				} else {
					AutonomousOfflineNum := 0
					for _, pod := range podList.Items {
						for _, node := range nodeList {
							if node.NodeIP == pod.Status.HostIP && node.Status == service.AutonomousOffline {
								AutonomousOfflineNum++
							}
						}
					}
					if AutonomousOfflineNum == 0 {
						status = uncompleted
						reason = "The desired replicas number is not available"
					} else {
						status = autonomousOffline
						reason = "The nodes offline"
					}
					// } else if AutonomousOfflineNum == len(podList.Items) {
					// 	status = autonomousOffline
					// 	reason = "The nodes offline"
					// } else {
					// 	status = partAutonomousOffline
					// 	reason = "The nodes offline"
					// }
				}

				(*serviceStatus).Status = status
				(*serviceStatus).Comment = "Reason: " + reason
				_, err = service.UpdateService(*serviceStatus, "status", "Comment")
				if err != nil {
					logs.Error("Failed to update deployment replicas.")
					break
				}
				continue
			} else if deployment.Status.Replicas == deployment.Status.AvailableReplicas {
				(*serviceStatus).Status = running
				_, err = service.UpdateService(*serviceStatus, "status", "Comment")
				if err != nil {
					logs.Error("Failed to update service in cluster.")
					break
				}
				continue
			}
		}

		// Check the service in k8s cluster status
		serviceK8s, err := service.GetK8sService((*serviceStatus).ProjectName, (*serviceStatus).Name)
		if serviceK8s == nil {
			logs.Info("Failed to get service in cluster", err)
			var reason = "The service is not established in cluster system"
			(*serviceStatus).Status = uncompleted
			(*serviceStatus).Comment = "Reason: " + reason
			_, err = service.UpdateService(*serviceStatus, "status", "Comment")
			if err != nil {
				logs.Error("Failed to update service in cluster.")
				break
			}
			continue
		}

		if serviceStatus.Status == uncompleted {
			logs.Info("The service is restored to running")
			(*serviceStatus).Status = running
			(*serviceStatus).Comment = ""
			_, err = service.UpdateService(*serviceStatus, "status", "Comment")
			if err != nil {
				logs.Error("Failed to update service status.")
				break
			}
			continue
		}
	}
	return err
}

//get service list
func (p *ServiceController) GetServiceListAction() {
	serviceName := p.GetString("service_name")
	source, _ := p.GetInt("service_source", -1)
	sourceid, _ := p.GetInt("service_sourceid", -1)
	pageIndex, _ := p.GetInt("page_index", 0)
	pageSize, _ := p.GetInt("page_size", 0)
	orderField := p.GetString("order_field", "creation_time")
	orderAsc, _ := p.GetInt("order_asc", 0)

	orderFieldValue, err := service.ParseOrderField("service", orderField)
	if err != nil {
		p.CustomAbortAudit(http.StatusBadRequest, err.Error())
		return
	}

	if pageIndex == 0 && pageSize == 0 {
		var sourcePtr *int
		if source != -1 {
			sourcePtr = &source
		}
		var sourceidPtr *int64
		if sourceid != -1 {
			var sourceid64 int64 = int64(sourceid)
			sourceidPtr = &sourceid64
		}
		serviceStatus, err := service.GetServiceList(serviceName, p.CurrentUser.ID, sourcePtr, sourceidPtr)
		if err != nil {
			p.InternalError(err)
			return
		}
		err = syncK8sStatus(serviceStatus)
		if err != nil {
			p.InternalError(err)
			return
		}
		p.RenderJSON(serviceStatus)
	} else {
		paginatedServiceStatus, err := service.GetPaginatedServiceList(serviceName, p.CurrentUser.ID, pageIndex, pageSize, orderFieldValue, orderAsc)
		if err != nil {
			p.InternalError(err)
			return
		}
		err = syncK8sStatus(paginatedServiceStatus.ServiceStatusList)
		if err != nil {
			p.InternalError(err)
			return
		}
		p.RenderJSON(paginatedServiceStatus)
	}
}

// API to create service config
func (p *ServiceController) CreateServiceConfigAction() {
	var reqServiceProject model.ServiceProject
	var err error

	err = p.ResolveBody(&reqServiceProject)
	if err != nil {
		return
	}

	//Judge authority
	p.ResolveUserPrivilegeByID(reqServiceProject.ProjectID)

	//Assign and return Service ID with mysql
	var newservice model.ServiceStatus
	newservice.ProjectID = reqServiceProject.ProjectID
	newservice.ProjectName = reqServiceProject.ProjectName
	newservice.Status = preparing // 0: preparing 1: running 2: suspending
	newservice.OwnerID = p.CurrentUser.ID

	serviceInfo, err := service.CreateServiceConfig(newservice)
	if err != nil {
		p.InternalError(err)
		return
	}
	p.RenderJSON(serviceInfo.ID)
}

func (p *ServiceController) DeleteServiceAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(s.ProjectID)

	if err = service.CheckServiceDeletable(s); err != nil {
		p.InternalError(err)
		return
	}
	// Call stop service if running
	if s.Status != stopped {
		err = service.StopServiceK8s(s)
		if err != nil {
			p.InternalError(err)
			return
		}
	}

	// 	Delete the service's autoscale rule
	hpas, err := service.ListAutoScales(s)
	if err != nil {
		p.InternalError(err)
		return
	}
	for _, hpa := range hpas {
		err := service.DeleteAutoScale(s, hpa.ID)
		if err != nil {
			p.InternalError(err)
			return
		}
		logs.Debug("Deleted Hpa %d %s", hpa.ID, hpa.HPAName)
	}

	isSuccess, err := service.DeleteService(s.ID)
	if err != nil {
		p.InternalError(err)
		return
	}
	if !isSuccess {
		p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Failed to delete service with ID: %d", s.ID))
		return
	}

	//delete repo files of the service
	p.ResolveRepoServicePath(s.ProjectName, s.Name)
	p.RemoveItemsToRepo(filepath.Join(s.Name, serviceFilename), filepath.Join(s.Name, deploymentFilename))

}

func (p *ServiceController) ResolveServiceOwnerRepoServicePath(projectName string, serviceName string, serviceOwner string) {
	repoName, err := service.ResolveRepoName(projectName, serviceOwner)
	if err != nil {
		p.CustomAbortAudit(http.StatusPreconditionFailed, fmt.Sprintf("Failed to generate repo path: %+v", err))
		return
	}
	p.RepoName = repoName
	p.RepoPath = service.ResolveRepoPath(repoName, serviceOwner)
	p.RepoServicePath = filepath.Join(p.RepoPath, serviceName)
	logs.Debug("Set repo path at file upload: %s and repo name: %s", p.RepoPath, p.RepoName)
}

// API to deploy service
func (p *ServiceController) ToggleServiceAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	var reqServiceToggle model.ServiceToggle
	err = p.ResolveBody(&reqServiceToggle)
	if err != nil {
		return
	}

	//Judge authority
	p.ResolveUserPrivilegeByID(s.ProjectID)

	if s.Status == stopped && reqServiceToggle.Toggle == 0 {
		p.CustomAbortAudit(http.StatusBadRequest, "Service already stopped.")
		return
	}

	if s.Status == running && reqServiceToggle.Toggle == 1 {
		p.CustomAbortAudit(http.StatusBadRequest, "Service already running.")
		return
	}
	p.ResolveServiceOwnerRepoServicePath(s.ProjectName, s.Name, s.OwnerName)
	if devOpsOpt() == "legacy" {
		if _, err := os.Stat(p.RepoServicePath); os.IsNotExist(err) {
			p.CustomAbortAudit(http.StatusPreconditionFailed, "Service restored from initialization, cannot be switched.")
			return
		}
	}
	if reqServiceToggle.Toggle == 0 {
		// stop service
		err = service.StopServiceK8s(s)
		if err != nil {
			p.InternalError(err)
			return
		}
		// Update service status DB
		_, err = service.UpdateServiceStatus(s.ID, stopped)
		if err != nil {
			p.InternalError(err)
			return
		}
	} else {
		// start service
		logs.Debug("Deploy service by YAML with project name: %s", s.ProjectName)
		err := service.DeployServiceByYaml(s.ProjectName, p.RepoServicePath)
		if err != nil {
			p.ParseError(err, c.ParsePostK8sError)
			return
		}
		// Push deployment to Git repo
		// items := []string{filepath.Join(s.Name, deploymentFilename), filepath.Join(s.Name, serviceFilename)}
		// p.PushItemsToRepo(items...)
		// p.CollaborateWithPullRequest("master", "master", items...)
		// p.MergeCollaborativePullRequest()
		// Update service status DB
		_, err = service.UpdateServiceStatus(s.ID, running)
		if err != nil {
			p.InternalError(err)
			return
		}
	}
}

func (p *ServiceController) GetServiceInfoAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	if s.Public != 1 {
		p.ResolveUserPrivilegeByID(s.ProjectID)
	}

	serviceStatus := &model.Service{}
	if s.Type != model.ServiceTypeEdgeComputing {
		serviceStatus, err = service.GetServiceByK8sassist(s.ProjectName, s.Name)
		if err != nil {
			p.ParseError(err, c.ParseGetK8sError)
			return
		}
	} else {
		deployment, _, err := service.GetDeployment(s.ProjectName, s.Name)
		if err != nil {
			p.ParseError(err, c.ParseGetK8sError)
			return
		}
		for _, container := range deployment.Spec.Template.Spec.Containers {
			for _, port := range container.Ports {
				serviceStatus.Ports = append(serviceStatus.Ports, model.ServicePort{
					Name:     port.Name,
					Protocol: string(port.Protocol),
					NodePort: port.HostPort,
				})
			}
		}
	}

	//Get NodeIP
	//endpointUrl format /api/v1/namespaces/default/endpoints/
	nodesStatus, err := service.GetNodesStatus()
	if err != nil {
		p.ParseError(err, c.ParseGetK8sError)
		return
	}
	if len(serviceStatus.Ports) == 0 || len(nodesStatus.Items) == 0 {
		p.RenderJSON("NA")
		return
	}

	var serviceInfo model.ServiceInfoStruct
	for _, ports := range serviceStatus.Ports {
		serviceInfo.NodePort = append(serviceInfo.NodePort, ports.NodePort)
	}
	for _, items := range nodesStatus.Items {
		serviceInfo.NodeName = append(serviceInfo.NodeName, items.Status.Addresses...)
	}
	serviceInfo.ServiceContainers, err = service.GetServiceContainers(s)
	if err != nil {
		p.ParseError(err, c.ParseGetK8sError)
		return
	}
	p.RenderJSON(serviceInfo)
}

func (p *ServiceController) GetServicePodLogsAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	if s.Public != 1 {
		p.ResolveUserPrivilegeByID(s.ProjectID)
	}
	podName := p.Ctx.Input.Param(":podname")
	readCloser, err := service.GetK8sPodLogs(s.ProjectName, podName, p.GeneratePodLogOptions())
	if err != nil {
		p.ParseError(err, c.ParseGetK8sError)
		return
	}
	defer readCloser.Close()
	_, err = io.Copy(&utils.FlushResponseWriter{p.Ctx.Output.Context.ResponseWriter}, readCloser)
	if err != nil {
		logs.Error("get service logs error:%+v", err)
	}
}

func (p *ServiceController) GetServiceStatusAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(s.ProjectID)
	serviceStatus, err := service.GetServiceByK8sassist(s.ProjectName, s.Name)
	if err != nil {
		p.ParseError(err, c.ParseGetK8sError)
		return
	}
	p.RenderJSON(serviceStatus)
}

func (p *ServiceController) ServicePublicityAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(s.ProjectID)
	var reqServiceUpdate model.ServicePublicityUpdate
	err = p.ResolveBody(&reqServiceUpdate)
	if err != nil {
		return
	}
	if s.Public != reqServiceUpdate.Public {
		_, err := service.UpdateServicePublic(s.ID, reqServiceUpdate.Public)
		if err != nil {
			p.InternalError(err)
			return
		}
	}
}

func (p *ServiceController) DeleteServiceConfigAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	// Get the path of the service config files
	p.ResolveUserPrivilege(s.ProjectName)
	p.ResolveRepoServicePath(s.ProjectName, s.Name)
	logs.Debug("Service config path: %s", p.RepoServicePath)

	// Delete yaml files
	// TODO
	err = service.DeleteServiceConfigYaml(p.RepoServicePath)
	if err != nil {
		logs.Info("failed to delete service yaml", p.RepoServicePath)
		p.InternalError(err)
		return
	}

	// For terminated service config, actually delete it in DB
	_, err = service.DeleteServiceByID(s.ID)
	if err != nil {
		p.InternalError(err)
		return
	}
}

func (p *ServiceController) DeleteDeploymentAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	// Get the path of the service config files
	p.ResolveUserPrivilege(s.ProjectName)
	p.ResolveRepoServicePath(s.ProjectName, s.Name)
	logs.Debug("Service config path: %s", p.RepoServicePath)

	// TODO clear kube-master, even if the service is not deployed successfully
	p.RemoveItemsToRepo(filepath.Join(s.Name, deploymentFilename))

	// Delete yaml files
	err = service.DeleteServiceConfigYaml(p.RepoServicePath)
	if err != nil {
		logs.Info("Failed to delete service yaml under path: %s", p.RepoServicePath)
		p.InternalError(err)
		return
	}

	// For terminated service config, actually delete it in DB
	_, err = service.DeleteServiceByID(s.ID)
	if err != nil {
		p.InternalError(err)
		return
	}
}

func (p *ServiceController) StoreServiceRoute() {
	serviceIdentity := p.GetString("service_identity")
	serviceURL := p.GetString("service_url")
	c.MemoryCache.Put(strings.ToLower(serviceIdentity), serviceURL, time.Second*time.Duration(c.TokenCacheExpireSeconds))
	logs.Debug("Service identity: %s, URL: %s", serviceIdentity, serviceURL)
}

func (p *ServiceController) ServiceExists() {
	projectName := p.GetString("project_name")
	p.ResolveProjectMember(projectName)
	serviceName := p.GetString("service_name")
	isServiceExists, err := service.ServiceExists(serviceName, projectName)
	if err != nil {
		p.InternalError(err)
		logs.Error("Check service name failed, error: %+v", err.Error())
		return
	}
	if isServiceExists == true {
		p.CustomAbortAudit(http.StatusConflict, serverNameDuplicateErr.Error())
		return
	}
}

func (p *ServiceController) ScaleServiceAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(s.ProjectID)

	var reqServiceScale model.ServiceScale
	err = p.ResolveBody(&reqServiceScale)
	if err != nil {
		return
	}
	// change the replica number of service
	res, err := service.ScaleReplica(s, reqServiceScale.Replica)
	if res != true {
		logs.Info("Failed to scale service replica", s, reqServiceScale.Replica)
		p.InternalError(err)
		return
	}
	logs.Info("Scale service replica successfully")

	_, deploymentFileInfo, err := service.GetDeployment(s.ProjectName, s.Name)
	if err != nil {
		logs.Error("Failed to get deployment %s", s.Name)
		return
	}
	p.ResolveRepoServicePath(s.ProjectName, s.Name)
	err = utils.GenerateFile(deploymentFileInfo, p.RepoServicePath, deploymentFilename)
	if err != nil {
		logs.Error("Failed to update file of deployment %s", s.Name)
		return
	}
	p.PushItemsToRepo(filepath.Join(s.Name, deploymentFilename))
}

//get selectable service list
func (p *ServiceController) GetSelectableServicesAction() {
	projectName := p.GetString("project_name")
	p.ResolveProjectMember(projectName)
	logs.Info("Get selectable service list for", projectName)
	serviceList, err := service.GetServicesByProjectName(projectName)
	if err != nil {
		logs.Error("Failed to get selectable services.")
		p.InternalError(err)
		return
	}
	p.RenderJSON(serviceList)
}

func (f *ServiceController) resolveUploadedYamlFile(uploadedFileName string) (func(fileName string, serviceInfo *model.ServiceStatus) error, io.Reader, error) {
	uploadedFile, _, err := f.GetFile(uploadedFileName)
	if err != nil {
		if err.Error() == "http: no such file" {
			f.CustomAbortAudit(http.StatusBadRequest, "Missing file: "+uploadedFileName)
			return nil, nil, err
		}
		f.InternalError(err)
		return nil, nil, err
	}

	return func(fileName string, serviceInfo *model.ServiceStatus) error {
		f.ResolveRepoServicePath(serviceInfo.ProjectName, serviceInfo.Name)
		err = utils.CheckFilePath(f.RepoServicePath)
		if err != nil {
			f.InternalError(err)
			return nil
		}
		return f.SaveToFile(uploadedFileName, filepath.Join(f.RepoServicePath, fileName))
	}, uploadedFile, nil
}

func (f *ServiceController) UploadYamlFileAction() {
	projectName := f.GetString("project_name")
	f.ResolveProjectMember(projectName)

	fhDeployment, deploymentFile, err := f.resolveUploadedYamlFile("deployment_file")
	if err != nil {
		return
	}
	fhService, serviceFile, err := f.resolveUploadedYamlFile("service_file")
	if err != nil {
		return
	}
	deployInfo, err := service.CheckDeployYamlConfig(serviceFile, deploymentFile, projectName)
	if err != nil {
		f.CustomAbortAudit(http.StatusBadRequest, err.Error())
		return
	}

	serviceName := deployInfo.Service.ObjectMeta.Name
	serviceInfo, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		f.InternalError(err)
		return
	}
	if serviceInfo != nil {
		f.CustomAbortAudit(http.StatusBadRequest, "Service name has been used.")
		return
	}
	serviceInfo, err = service.CreateServiceConfig(model.ServiceStatus{
		Name:        serviceName,
		ProjectName: projectName,
		Status:      preparing, // 0: preparing 1: running 2: suspending
		OwnerID:     f.CurrentUser.ID,
		OwnerName:   f.CurrentUser.Username,
		Type:        service.GetServiceType(deployInfo.Service.Type),
	})
	if err != nil {
		f.InternalError(err)
		return
	}
	err = fhDeployment(deploymentFilename, serviceInfo)
	if err != nil {
		f.InternalError(err)
		return
	}
	err = fhService(serviceFilename, serviceInfo)
	if err != nil {
		f.InternalError(err)
		return
	}
	f.RenderJSON(serviceInfo)
}

func (f *ServiceController) DownloadDeploymentYamlFileAction() {
	projectName := f.GetString("project_name")
	serviceName := f.GetString("service_name")
	serviceInfo, err := service.GetServiceByProject(serviceName, projectName)
	if err != nil {
		f.InternalError(err)
		return
	}
	if serviceInfo == nil {
		f.CustomAbortAudit(http.StatusBadRequest, "Service name is invalid.")
		return
	}
	f.ResolveRepoServicePath(projectName, serviceName)
	yamlType := f.GetString("yaml_type")

	switch yamlType {
	case "":
		f.CustomAbortAudit(http.StatusBadRequest, "No YAML type found.")
		return
	case deploymentType:
		f.resolveDownloadYaml(serviceInfo, deploymentFilename, service.GenerateDeploymentYamlFileFromK8s)
	case serviceType:
		f.resolveDownloadYaml(serviceInfo, serviceFilename, service.GenerateServiceYamlFileFromK8s)
	case statefulsetType:
		f.resolveDownloadYaml(serviceInfo, statefulsetFilename, service.GenerateStatefulSetYamlFileFromK8s)
	}

}

func (f *ServiceController) resolveDownloadYaml(serviceConfig *model.ServiceStatus, fileName string, generator func(*model.ServiceStatus, string) error) {
	logs.Debug("Current download yaml file: %s", fileName)
	//checkout the path of download
	err := utils.CheckFilePath(f.RepoServicePath)
	if err != nil {
		f.InternalError(err)
		return
	}
	absFileName := filepath.Join(f.RepoServicePath, fileName)
	err = generator(serviceConfig, f.RepoServicePath)
	if err != nil {
		f.ParseError(err, c.ParseGetK8sError)
		return
	}
	logs.Info("User: %s downloaded %s YAML file.", f.CurrentUser.Username, fileName)
	f.Ctx.Output.Download(absFileName, fileName)
}

func (p *ServiceController) GetScaleStatusAction() {
	// Get the current service status
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(s.ProjectID)
	scaleStatus, err := service.GetScaleStatus(s)
	if err != nil {
		logs.Debug("Get scale deployment status failed %s", s.Name)
		p.InternalError(err)
		return
	}
	p.RenderJSON(scaleStatus)
}

func (p *ServiceController) DeleteDeployAction() {
	var err error

	key := p.getKey()
	configService := NewConfigServiceStep(key)

	//Judge authority
	p.ResolveUserPrivilegeByID(configService.ProjectID)

	// Clean deployment and service

	s := model.ServiceStatus{Name: configService.ServiceName,
		ProjectName: configService.ProjectName,
	}

	err = service.StopServiceK8s(&s)
	if err != nil {
		logs.Error("Failed to clean service %s", s.Name)
		p.InternalError(err)
		return
	}

	//Clean data DB if existing
	serviceData, err := service.GetService(s, "name", "project_name")
	if serviceData != nil {
		isSuccess, err := service.DeleteService(serviceData.ID)
		if err != nil {
			p.InternalError(err)
			return
		}
		if !isSuccess {
			p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Failed to delete service with ID: %d", s.ID))
			return
		}
	}

	//delete repo files of the service
	p.ResolveRepoServicePath(s.ProjectName, s.Name)
	p.RemoveItemsToRepo(filepath.Join(s.Name, serviceFilename), filepath.Join(s.Name, deploymentFilename))

	//clean the config step
	err = DeleteConfigServiceStep(key)
	if err != nil {
		logs.Debug("Failed to clean the config steps")
		p.InternalError(err)
		return
	}
}

//get service nodeport list
func (p *ServiceController) GetServiceNodePorts() {
	projectName := p.GetString("project_name")

	if projectName != "" {
		p.ResolveProjectMember(projectName)
		logs.Info("Get service nodeport list for", projectName)
	}

	//Check nodeports in cluster
	nodeportList, err := service.GetNodePortsK8s(projectName)
	if err != nil {
		logs.Error("Failed to get selectable services.")
		p.InternalError(err)
		return
	}

	//Check nodeports in DB
	p.RenderJSON(nodeportList)
}

//import cluster services
func (p *ServiceController) ImportServicesAction() {

	if p.IsSysAdmin == false {
		p.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to import services.")
		return
	}

	projectList, err := service.GetProjectsByUser(model.Project{}, p.CurrentUser.ID)
	if err != nil {
		logs.Error("Failed to get projects.")
		p.InternalError(err)
		return
	}

	for _, project := range projectList {
		err := service.SyncServiceWithK8s(project.Name)
		if err != nil {
			logs.Error("Failed to sync service for project %s.", project.Name)
			p.InternalError(err)
			return
		}
	}
	logs.Debug("imported services from cluster successfully")
}

// DeployStatefulSetAction is an api action to deploy statefulset
func (p *ServiceController) DeployStatefulSetAction() {
	key := p.getKey()
	configService := NewConfigServiceStep(key)

	if configService.ServiceType != model.ServiceTypeStatefulSet {
		p.CustomAbortAudit(http.StatusBadRequest, "Invalid service type.")
		return
	}

	//Judge authority
	project := p.ResolveUserPrivilegeByID(configService.ProjectID)

	var newservice model.ServiceStatus
	newservice.Name = configService.ServiceName
	newservice.ProjectID = configService.ProjectID
	newservice.Status = preparing // 0: preparing 1: running 2: suspending
	newservice.OwnerID = p.CurrentUser.ID
	newservice.OwnerName = p.CurrentUser.Username
	newservice.Public = configService.Public
	newservice.ProjectName = project.Name
	newservice.Type = configService.ServiceType

	serviceInfo, err := service.CreateServiceConfig(newservice)
	if err != nil {
		p.InternalError(err)
		return
	}

	statefulsetInfo, err := service.DeployStatefulSet((*model.ConfigServiceStep)(configService), c.RegistryBaseURI())
	if err != nil {
		p.ParseError(err, c.ParsePostK8sError)
		return
	}

	p.ResolveRepoServicePath(project.Name, newservice.Name)
	err = service.GenerateStatefulSetYamlFiles(statefulsetInfo, p.RepoServicePath)
	if err != nil {
		p.InternalError(err)
		return
	}

	statefulsetFile := filepath.Join(newservice.Name, statefulsetFilename)
	serviceFile := filepath.Join(newservice.Name, serviceFilename)
	items := []string{statefulsetFile, serviceFile}
	p.PushItemsToRepo(items...)

	updateService := model.ServiceStatus{ID: serviceInfo.ID, Status: uncompleted, ServiceYaml: string(statefulsetInfo.ServiceFileInfo),
		DeploymentYaml: string(statefulsetInfo.StatefulSetFileInfo)}
	_, err = service.UpdateService(updateService, "status", "service_yaml", "deployment_yaml")
	if err != nil {
		p.InternalError(err)
		return
	}

	err = DeleteConfigServiceStep(key)
	if err != nil {
		p.InternalError(err)
		return
	}
	logs.Info("Service with ID:%d has been deleted in cache.", serviceInfo.ID)

	configService.ServiceID = serviceInfo.ID
	p.RenderJSON(configService)
}

func (p *ServiceController) DeleteStatefulSetAction() {
	var err error
	s, err := p.resolveServiceInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(s.ProjectID)

	if err = service.CheckServiceDeletable(s); err != nil {
		p.InternalError(err)
		return
	}
	// Call stop service if running
	if s.Status != stopped {
		err = service.StopStatefulSetK8s(s)
		if err != nil {
			p.InternalError(err)
			return
		}
	}

	// 	Delete the service's autoscale rule later

	isSuccess, err := service.DeleteService(s.ID)
	if err != nil {
		p.InternalError(err)
		return
	}
	if !isSuccess {
		p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Failed to delete service with ID: %d", s.ID))
		return
	}

	//delete repo files of the service
	p.ResolveRepoServicePath(s.ProjectName, s.Name)
	p.RemoveItemsToRepo(filepath.Join(s.Name, serviceFilename), filepath.Join(s.Name, statefulsetFilename))

}
