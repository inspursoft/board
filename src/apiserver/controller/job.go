package controller

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
)

var (
	jobNameDuplicateErr = errors.New("ERR_DUPLICATE_JOB_NAME")
)

type JobController struct {
	BaseController
}

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

	updateJob := model.JobStatusMO{ID: jobInfo.ID, Status: uncompleted, Yaml: string(jobDeployInfo.JobFileInfo)}
	_, err = service.UpdateJob(updateJob, "status", "yaml")
	if err != nil {
		p.internalError(err)
		return
	}

	config.ID = jobInfo.ID
	p.renderJSON(config)
}

//get job list
func (p *JobController) GetJobListAction() {
	jobName := p.GetString("job_name")
	pageIndex, _ := p.GetInt("page_index", 0)
	pageSize, _ := p.GetInt("page_size", 0)
	orderField := p.GetString("order_field", "creation_time")
	orderAsc, _ := p.GetInt("order_asc", 0)

	orderFieldValue, err := service.ParseOrderField("job", orderField)
	if err != nil {
		p.customAbort(http.StatusBadRequest, err.Error())
		return
	}

	if pageIndex == 0 && pageSize == 0 {
		jobStatus, err := service.GetJobList(jobName, p.currentUser.ID)
		if err != nil {
			p.internalError(err)
			return
		}
		err = service.SyncJobK8sStatus(jobStatus)
		if err != nil {
			p.internalError(err)
			return
		}
		p.renderJSON(jobStatus)
	} else {
		paginatedJobStatus, err := service.GetPaginatedJobList(jobName, p.currentUser.ID, pageIndex, pageSize, orderFieldValue, orderAsc)
		if err != nil {
			p.internalError(err)
			return
		}
		err = service.SyncJobK8sStatus(paginatedJobStatus.JobStatusList)
		if err != nil {
			p.internalError(err)
			return
		}
		p.renderJSON(paginatedJobStatus)
	}
}

//get job
func (p *JobController) GetJobAction() {
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(j.ProjectID)

	jobStatus := []*model.JobStatusMO{j}
	err = service.SyncJobK8sStatus(jobStatus)
	if err != nil {
		p.internalError(err)
		return
	}
	p.renderJSON(jobStatus[0])
}

func (p *JobController) DeleteJobAction() {
	var err error
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(j.ProjectID)

	// stop service
	err = service.StopJobK8s(j)
	if err != nil {
		p.internalError(err)
		return
	}
	//TODO: where is the deletion of kubernetes job object?, write it here or in service method? do we need another state to reference it?
	isSuccess, err := service.DeleteJob(j.ID)
	if err != nil {
		p.internalError(err)
		return
	}
	if !isSuccess {
		p.customAbort(http.StatusBadRequest, fmt.Sprintf("Failed to delete job with ID: %d", j.ID))
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
	}
}

func (p *JobController) GetJobPodsAction() {
	var err error
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(j.ProjectID)
	pods, err := service.GetK8sJobPods(j)
	if err != nil {
		p.parseError(err, parseGetK8sError)
		return
	}
	p.renderJSON(pods)
}

func (p *JobController) GetJobLogsAction() {
	var err error
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.resolveUserPrivilegeByID(j.ProjectID)
	podName := p.Ctx.Input.Param(":podname")
	readCloser, err := service.GetK8sJobLogs(j, podName, p.generatePodLogOptions())
	if err != nil {
		p.parseError(err, parseGetK8sError)
		return
	}
	defer readCloser.Close()
	_, err = io.Copy(&utils.FlushResponseWriter{p.Ctx.Output.Context.ResponseWriter}, readCloser)
	if err != nil {
		logs.Error("get job logs error:%+v", err)
	}
}

func (p *JobController) generatePodLogOptions() *model.PodLogOptions {
	var err error
	opt := &model.PodLogOptions{}
	opt.Container = p.GetString("container")
	opt.Follow, err = p.GetBool("follow", false)
	if err != nil {
		logs.Warn("Follow parameter %s is invalid: %+v", p.GetString("follow"), err)
	}
	opt.Previous, err = p.GetBool("previous", false)
	if err != nil {
		logs.Warn("Privious parameter %s is invalid: %+v", p.GetString("privious"), err)
	}
	opt.Timestamps, err = p.GetBool("timestamps", false)
	if err != nil {
		logs.Warn("Timestamps parameter %s is invalid: %+v", p.GetString("timestamps"), err)
	}

	if p.GetString("since_seconds") != "" {
		since, err := p.GetInt64("since_seconds")
		if err != nil {
			logs.Warn("SinceSeconds parameter %s is invalid: %+v", p.GetString("since_seconds"), err)
		} else {
			opt.SinceSeconds = &since
		}
	}

	since := p.GetString("since_time")
	if since != "" {
		sinceTime, err := time.Parse(time.RFC3339, since)
		if err != nil {
			logs.Warn("since_time parameter %s is invalid: %+v", since, err)
		} else {
			opt.SinceTime = &sinceTime
		}
	}

	tail, err := p.GetInt64("tail_lines", -1)
	if err != nil {
		logs.Warn("tail_lines parameter %s is invalid: %+v", p.GetString("tail_lines"), err)
	} else if tail != -1 {
		opt.TailLines = &tail
	}

	limit, err := p.GetInt64("limit_bytes", -1)
	if err != nil {
		logs.Warn("limit_bytes parameter %s is invalid: %+v", p.GetString("limit_bytes"), err)
	} else if limit != -1 {
		opt.LimitBytes = &limit
	}

	return opt
}
func (p *JobController) GetSelectableJobsAction() {
	projectName := p.GetString("project_name")
	p.resolveProjectMember(projectName)
	logs.Info("Get selectable job list for", projectName)
	jobList, err := service.GetJobsByProjectName(projectName)
	if err != nil {
		logs.Error("Failed to get selectable jobs.")
		p.internalError(err)
		return
	}
	p.renderJSON(jobList)
}
