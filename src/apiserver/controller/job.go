package controller

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"github.com/astaxie/beego/logs"
)

const (
	jobFilename = "job.yaml"
)

var (
	jobNameDuplicateErr = errors.New("ERR_DUPLICATE_JOB_NAME")
)

type JobController struct {
	BaseController
}

//func (p *JobController) generateDeploymentTravis(serviceName, deploymentURL, serviceURL string) error {
//	userID := p.currentUser.ID
//	var travisCommand travis.TravisCommand
//	travisCommand.Script.Commands = []string{}
//	items := []string{
//		fmt.Sprintf("curl \"%s/jenkins-job/%d/$BUILD_NUMBER\"", boardAPIBaseURL(), userID),
//	}
//	if deploymentURL != "" {
//		items = append(items, fmt.Sprintf("#curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/deployment.yaml %s", serviceName, deploymentURL))
//	}
//	if serviceURL != "" {
//		items = append(items, fmt.Sprintf("#curl -X POST -H 'Content-Type: application/yaml' --data-binary @%s/service.yaml %s", serviceName, serviceURL))
//	}
//	travisCommand.Script.Commands = items
//	return travisCommand.GenerateCustomTravis(p.repoPath)
//}
//
//func (p *JobController) getKey() string {
//	return strconv.Itoa(int(p.currentUser.ID))
//}

func (p *JobController) resolveJobInfo() (j *model.JobStatusMO, err error) {
	jobID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.internalError(err)
		return
	}
	// Get the project info of this service
	j, err = service.GetJobByID(int64(jobID))
	if err != nil {
		p.internalError(err)
		return
	}
	if j == nil {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Invalid job ID: %d", jobID))
		return
	}
	return
}

func (p *JobController) DeployJobAction() {
	var config model.JobConfig
	err := p.resolveBody(&config)
	if err != nil {
		return
	}
	//Judge authority
	project := p.resolveUserPrivilegeByID(config.ProjectID)

	var newjob model.JobStatusMO
	newjob.Name = config.Name
	newjob.ProjectID = config.ProjectID
	newjob.Status = preparing // 0: preparing 1: running 2: suspending
	newjob.OwnerID = p.currentUser.ID
	newjob.OwnerName = p.currentUser.Username
	newjob.ProjectName = project.Name

	jobInfo, err := service.CreateJob(newjob)
	if err != nil {
		p.internalError(err)
		return
	}

	jobDeployInfo, err := service.DeployJob(&config, registryBaseURI())
	if err != nil {
		p.parseError(err, parsePostK8sError)
		return
	}

	p.resolveRepoServicePath(project.Name, newjob.Name)
	err = service.SaveJobDeployYamlFiles(jobDeployInfo, p.repoServicePath, jobFilename)
	if err != nil {
		p.internalError(err)
		return
	}

	jobFile := filepath.Join(newjob.Name, jobFilename)
	p.pushItemsToRepo(jobFile)

	updateJob := model.JobStatusMO{ID: jobInfo.ID, Status: uncompleted, Yaml: string(jobDeployInfo.JobFileInfo)}
	_, err = service.UpdateJob(updateJob, "status", "yaml")
	if err != nil {
		p.internalError(err)
		return
	}

	config.ID = jobInfo.ID
	p.renderJSON(config)
}

////
//func syncK8sStatus(serviceList []*model.ServiceStatusMO) error {
//	var err error
//	// synchronize service status with the cluster system
//	for _, serviceStatusMO := range serviceList {
//		// Get serviceStatus from serviceStatusMO to adapt for updating services
//		serviceStatus := &serviceStatusMO.ServiceStatus
//		if (*serviceStatus).Status == stopped {
//			continue
//		}
//		// Check the deployment status
//		deployment, _, err := service.GetDeployment((*serviceStatus).ProjectName, (*serviceStatus).Name)
//		if deployment == nil && serviceStatus.Name != k8sServices {
//			logs.Info("Failed to get deployment", err)
//			var reason = "The deployment is not established in cluster system"
//			(*serviceStatus).Status = uncompleted
//			// TODO create a new field in serviceStatus for reason
//			(*serviceStatus).Comment = "Reason: " + reason
//			_, err = service.UpdateService(*serviceStatus, "status", "Comment")
//			if err != nil {
//				logs.Error("Failed to update deployment.")
//				break
//			}
//			continue
//		} else {
//			if deployment.Status.Replicas > deployment.Status.AvailableReplicas {
//				logs.Debug("The desired replicas number is not available",
//					deployment.Status.Replicas, deployment.Status.AvailableReplicas)
//				(*serviceStatus).Status = uncompleted
//				reason := "The desired replicas number is not available"
//				(*serviceStatus).Comment = "Reason: " + reason
//				_, err = service.UpdateService(*serviceStatus, "status", "Comment")
//				if err != nil {
//					logs.Error("Failed to update deployment replicas.")
//					break
//				}
//				continue
//			}
//		}
//
//		// Check the service in k8s cluster status
//		serviceK8s, err := service.GetK8sService((*serviceStatus).ProjectName, (*serviceStatus).Name)
//		if serviceK8s == nil {
//			logs.Info("Failed to get service in cluster", err)
//			var reason = "The service is not established in cluster system"
//			(*serviceStatus).Status = uncompleted
//			(*serviceStatus).Comment = "Reason: " + reason
//			_, err = service.UpdateService(*serviceStatus, "status", "Comment")
//			if err != nil {
//				logs.Error("Failed to update service in cluster.")
//				break
//			}
//			continue
//		}
//
//		if serviceStatus.Status == uncompleted {
//			logs.Info("The service is restored to running")
//			(*serviceStatus).Status = running
//			(*serviceStatus).Comment = ""
//			_, err = service.UpdateService(*serviceStatus, "status", "Comment")
//			if err != nil {
//				logs.Error("Failed to update service status.")
//				break
//			}
//			continue
//		}
//	}
//	return err
//}

//get service list
func (p *JobController) GetJobListAction() {
	jobName := p.GetString("job_name")
	pageIndex, _ := p.GetInt("page_index", 0)
	pageSize, _ := p.GetInt("page_size", 0)
	orderField := p.GetString("order_field", "creation_time")
	orderAsc, _ := p.GetInt("order_asc", 0)
	if pageIndex == 0 && pageSize == 0 {
		jobStatus, err := service.GetJobList(jobName, p.currentUser.ID)
		if err != nil {
			p.internalError(err)
			return
		}
		//		err = syncK8sStatus(serviceStatus)
		//		if err != nil {
		//			p.internalError(err)
		//			return
		//		}
		p.renderJSON(jobStatus)
	} else {
		paginatedJobStatus, err := service.GetPaginatedJobList(jobName, p.currentUser.ID, pageIndex, pageSize, orderField, orderAsc)
		if err != nil {
			p.internalError(err)
			return
		}
		//		err = syncK8sStatus(paginatedServiceStatus.ServiceStatusList)
		//		if err != nil {
		//			p.internalError(err)
		//			return
		//		}
		p.renderJSON(paginatedJobStatus)
	}
}

//// API to create service config
//func (p *JobController) CreateJobAction() {
//	var reqServiceProject model.ServiceProject
//	var err error
//
//	err = p.resolveBody(&reqServiceProject)
//	if err != nil {
//		return
//	}
//
//	//Judge authority
//	p.resolveUserPrivilegeByID(reqServiceProject.ProjectID)
//
//	//Assign and return Service ID with mysql
//	var newservice model.ServiceStatus
//	newservice.ProjectID = reqServiceProject.ProjectID
//	newservice.ProjectName = reqServiceProject.ProjectName
//	newservice.Status = preparing // 0: preparing 1: running 2: suspending
//	newservice.OwnerID = p.currentUser.ID
//
//	serviceInfo, err := service.CreateServiceConfig(newservice)
//	if err != nil {
//		p.internalError(err)
//		return
//	}
//	p.renderJSON(serviceInfo.ID)
//}

func (p *JobController) DeleteJobAction() {
	var err error
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(j.ProjectID)

	//TODO: where is the deletion of kubernetes job object?, write it here or in service method? do we need another state to reference it?
	isSuccess, err := service.DeleteJob(j.ID)
	if err != nil {
		p.internalError(err)
		return
	}
	if !isSuccess {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Failed to delete job with ID: %d", j.ID))
		return
	}

	//delete repo files of the job
	p.resolveRepoServicePath(j.ProjectName, j.Name)
	p.removeItemsToRepo(filepath.Join(j.Name, jobFilename))

}

// API to deploy service
func (p *JobController) ToggleJobAction() {
	var err error
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	var reqJobToggle model.JobToggle
	err = p.resolveBody(&reqJobToggle)
	if err != nil {
		return
	}

	//Judge authority
	p.resolveUserPrivilegeByID(j.ProjectID)

	if j.Status == stopped && reqJobToggle.Toggle == 0 {
		p.customAbort(http.StatusBadRequest, "Service already stopped.")
		return
	}

	if j.Status == running && reqJobToggle.Toggle == 1 {
		p.customAbort(http.StatusBadRequest, "Service already running.")
		return
	}

	p.resolveRepoServicePath(j.ProjectName, j.Name)
	if _, err := os.Stat(p.repoServicePath); os.IsNotExist(err) {
		p.customAbort(http.StatusPreconditionFailed, "Job restored from initialization, cannot be switched.")
		return
	}
	if reqJobToggle.Toggle == 0 {
		// stop service
		err = service.StopJobK8s(j)
		if err != nil {
			p.internalError(err)
			return
		}
		// Update job status DB
		_, err = service.UpdateJobStatus(j.ID, stopped)
		if err != nil {
			p.internalError(err)
			return
		}
	} else {
		// start service
		err := service.DeployJobByYaml(j.ProjectName, filepath.Join(p.repoServicePath, jobFilename))
		if err != nil {
			p.parseError(err, parsePostK8sError)
			return
		}
		// Push job to Git repo
		p.pushItemsToRepo(filepath.Join(j.Name, jobFilename))
		p.collaborateWithPullRequest("master", "master", filepath.Join(j.Name, jobFilename))

		// Update job status DB
		_, err = service.UpdateJobStatus(j.ID, running)
		if err != nil {
			p.internalError(err)
			return
		}
	}
}

func (p *JobController) GetJobStatusAction() {
	var err error
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(j.ProjectID)
	jobStatus, err := service.GetK8sJobByK8sassist(j.ProjectName, j.Name)
	if err != nil {
		p.parseError(err, parseGetK8sError)
		return
	}
	p.renderJSON(jobStatus)
}

//func (p *JobController) DeleteDeploymentAction() {
//	var err error
//	s, err := p.resolveJobInfo()
//	if err != nil {
//		return
//	}
//	// Get the path of the service config files
//	p.resolveUserPrivilege(s.ProjectName)
//	p.resolveRepoServicePath(s.ProjectName, s.Name)
//	logs.Debug("Service config path: %s", p.repoServicePath)
//
//	// TODO clear kube-master, even if the service is not deployed successfully
//	p.removeItemsToRepo(filepath.Join(s.Name, deploymentFilename))
//
//	// Delete yaml files
//	err = service.DeleteServiceConfigYaml(p.repoServicePath)
//	if err != nil {
//		logs.Info("Failed to delete service yaml under path: %s", p.repoServicePath)
//		p.internalError(err)
//		return
//	}
//
//	// For terminated service config, actually delete it in DB
//	_, err = service.DeleteServiceByID(s.ID)
//	if err != nil {
//		p.internalError(err)
//		return
//	}
//}

func (p *JobController) JobExists() {
	projectName := p.GetString("project_name")
	p.resolveProjectMember(projectName)
	jobName := p.GetString("job_name")
	isJobExists, err := service.JobExists(jobName, projectName)
	if err != nil {
		p.internalError(err)
		logs.Error("Check job name failed, error: %+v", err.Error())
		return
	}
	if isJobExists {
		p.customAbort(http.StatusConflict, jobNameDuplicateErr.Error())
		return
	}
}

func (f *JobController) resolveUploadedYamlFile(uploadedFileName string) (func(fileName string, jobInfo *model.JobStatusMO) error, io.Reader, error) {
	uploadedFile, _, err := f.GetFile(uploadedFileName)
	if err != nil {
		if err.Error() == "http: no such file" {
			f.customAbort(http.StatusBadRequest, "Missing file: "+uploadedFileName)
			return nil, nil, err
		}
		f.internalError(err)
		return nil, nil, err
	}

	return func(fileName string, jobInfo *model.JobStatusMO) error {
		f.resolveRepoServicePath(jobInfo.ProjectName, jobInfo.Name)
		err = utils.CheckFilePath(f.repoServicePath)
		if err != nil {
			f.internalError(err)
			return nil
		}
		return f.SaveToFile(uploadedFileName, filepath.Join(f.repoServicePath, fileName))
	}, uploadedFile, nil
}

func (f *JobController) UploadYamlFileAction() {
	projectName := f.GetString("project_name")
	f.resolveProjectMember(projectName)

	fhJob, jobFile, err := f.resolveUploadedYamlFile("file")
	if err != nil {
		return
	}
	k8sjobInfo, err := service.CheckJobYamlConfig(jobFile, projectName)
	if err != nil {
		f.customAbort(http.StatusBadRequest, err.Error())
		return
	}

	jobName := k8sjobInfo.Name
	job, err := service.GetJobByProject(jobName, projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if job != nil {
		f.customAbort(http.StatusBadRequest, "Job name has been used.")
		return
	}
	jobInfo, err := service.CreateJob(model.JobStatusMO{
		Name:        jobName,
		ProjectName: projectName,
		Status:      preparing, // 0: preparing 1: running 2: suspending
		OwnerID:     f.currentUser.ID,
		OwnerName:   f.currentUser.Username,
	})
	if err != nil {
		f.internalError(err)
		return
	}
	err = fhJob(jobFilename, jobInfo)
	if err != nil {
		f.internalError(err)
		return
	}
	f.renderJSON(jobInfo)
}

func (f *JobController) DownloadYamlFileAction() {
	projectName := f.GetString("project_name")
	jobName := f.GetString("job_name")
	jobInfo, err := service.GetJobByProject(jobName, projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if jobInfo == nil {
		f.customAbort(http.StatusBadRequest, "Job name is invalid.")
		return
	}
	f.resolveRepoServicePath(projectName, jobName)
	f.resolveDownloadYaml(jobInfo, jobFilename, service.GenerateJobYamlFileFromK8s)

}

func (f *JobController) resolveDownloadYaml(jobConfig *model.JobStatusMO, fileName string, generator func(*model.JobStatusMO, string, string) error) {
	logs.Debug("Current download yaml file: %s", fileName)
	//checkout the path of download
	err := utils.CheckFilePath(f.repoServicePath)
	if err != nil {
		f.internalError(err)
		return
	}
	absFileName := filepath.Join(f.repoServicePath, fileName)
	err = generator(jobConfig, f.repoServicePath, jobFilename)
	if err != nil {
		f.parseError(err, parseGetK8sError)
		return
	}
	logs.Info("User: %s downloaded %s YAML file.", f.currentUser.Username, fileName)
	f.Ctx.Output.Download(absFileName, fileName)
}

//func (p *JobController) DeleteDeployAction() {
//	var err error
//
//	//Judge authority
//	p.resolveUserPrivilegeByID(configService.ProjectID)
//
//	// Clean deployment and service
//
//	s := model.ServiceStatus{Name: configService.ServiceName,
//		ProjectName: configService.ProjectName,
//	}
//
//	err = service.StopServiceK8s(&s)
//	if err != nil {
//		logs.Error("Failed to clean service %s", s.Name)
//		p.internalError(err)
//		return
//	}
//
//	//Clean data DB if existing
//	serviceData, err := service.GetService(s, "name", "project_name")
//	if serviceData != nil {
//		isSuccess, err := service.DeleteService(serviceData.ID)
//		if err != nil {
//			p.internalError(err)
//			return
//		}
//		if !isSuccess {
//			p.customAbort(http.StatusBadRequest, fmt.Sprintf("Failed to delete service with ID: %d", s.ID))
//			return
//		}
//	}
//
//	//delete repo files of the service
//	p.resolveRepoServicePath(s.ProjectName, s.Name)
//	p.removeItemsToRepo(filepath.Join(s.Name, serviceFilename), filepath.Join(s.Name, deploymentFilename))
//
//	//clean the config step
//	err = DeleteConfigServiceStep(key)
//	if err != nil {
//		logs.Debug("Failed to clean the config steps")
//		p.internalError(err)
//		return
//	}
//}

//
////import cluster services
//func (p *JobController) ImportServicesAction() {
//
//	if p.isSysAdmin == false {
//		p.customAbort(http.StatusForbidden, "Insufficient privileges to import services.")
//		return
//	}
//
//	projectList, err := service.GetProjectsByUser(model.Project{}, p.currentUser.ID)
//	if err != nil {
//		logs.Error("Failed to get projects.")
//		p.internalError(err)
//		return
//	}
//
//	for _, project := range projectList {
//		err := service.SyncServiceWithK8s(project.Name)
//		if err != nil {
//			logs.Error("Failed to sync service for project %s.", project.Name)
//			p.internalError(err)
//			return
//		}
//	}
//	logs.Debug("imported services from cluster successfully")
//}
