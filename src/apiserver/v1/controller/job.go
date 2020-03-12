package controller

import (
	"errors"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	c "git/inspursoft/board/src/common/controller"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
)

var (
	jobNameDuplicateErr = errors.New("ERR_DUPLICATE_JOB_NAME")
)

type JobController struct {
	c.BaseController
}

func (p *JobController) resolveJobInfo() (j *model.JobStatusMO, err error) {
	jobID, err := strconv.Atoi(p.Ctx.Input.Param(":id"))
	if err != nil {
		p.InternalError(err)
		return
	}
	// Get the project info of this service
	j, err = service.GetJobByID(int64(jobID))
	if err != nil {
		p.InternalError(err)
		return
	}
	if j == nil {
		p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Invalid job ID: %d", jobID))
		return
	}
	return
}

func (p *JobController) DeployJobAction() {
	var config model.JobConfig
	err := p.ResolveBody(&config)
	if err != nil {
		return
	}
	//Judge authority
	project := p.ResolveUserPrivilegeByID(config.ProjectID)

	var newjob model.JobStatusMO
	newjob.Name = config.Name
	newjob.ProjectID = config.ProjectID
	newjob.Status = preparing // 0: preparing 1: running 2: suspending
	newjob.OwnerID = p.CurrentUser.ID
	newjob.OwnerName = p.CurrentUser.Username
	newjob.ProjectName = project.Name

	jobInfo, err := service.CreateJob(newjob)
	if err != nil {
		p.InternalError(err)
		return
	}

	jobDeployInfo, err := service.DeployJob(&config, c.RegistryBaseURI())
	if err != nil {
		p.ParseError(err, c.ParsePostK8sError)
		return
	}

	updateJob := model.JobStatusMO{ID: jobInfo.ID, Status: uncompleted, Yaml: string(jobDeployInfo.JobFileInfo)}
	_, err = service.UpdateJob(updateJob, "status", "yaml")
	if err != nil {
		p.InternalError(err)
		return
	}

	config.ID = jobInfo.ID
	p.RenderJSON(config)
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
		p.CustomAbortAudit(http.StatusBadRequest, err.Error())
		return
	}

	if pageIndex == 0 && pageSize == 0 {
		jobStatus, err := service.GetJobList(jobName, p.CurrentUser.ID)
		if err != nil {
			p.InternalError(err)
			return
		}
		err = service.SyncJobK8sStatus(jobStatus)
		if err != nil {
			p.InternalError(err)
			return
		}
		p.RenderJSON(jobStatus)
	} else {
		paginatedJobStatus, err := service.GetPaginatedJobList(jobName, p.CurrentUser.ID, pageIndex, pageSize, orderFieldValue, orderAsc)
		if err != nil {
			p.InternalError(err)
			return
		}
		err = service.SyncJobK8sStatus(paginatedJobStatus.JobStatusList)
		if err != nil {
			p.InternalError(err)
			return
		}
		p.RenderJSON(paginatedJobStatus)
	}
}

//get job
func (p *JobController) GetJobAction() {
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(j.ProjectID)

	jobStatus := []*model.JobStatusMO{j}
	err = service.SyncJobK8sStatus(jobStatus)
	if err != nil {
		p.InternalError(err)
		return
	}
	p.RenderJSON(jobStatus[0])
}

func (p *JobController) DeleteJobAction() {
	var err error
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(j.ProjectID)

	// stop service
	err = service.StopJobK8s(j)
	if err != nil {
		p.InternalError(err)
		return
	}
	//TODO: where is the deletion of kubernetes job object?, write it here or in service method? do we need another state to reference it?
	isSuccess, err := service.DeleteJob(j.ID)
	if err != nil {
		p.InternalError(err)
		return
	}
	if !isSuccess {
		p.CustomAbortAudit(http.StatusBadRequest, fmt.Sprintf("Failed to delete job with ID: %d", j.ID))
	}

}

func (p *JobController) GetJobStatusAction() {
	var err error
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(j.ProjectID)
	jobStatus, err := service.GetK8sJobByK8sassist(j.ProjectName, j.Name)
	if err != nil {
		p.ParseError(err, c.ParseGetK8sError)
		return
	}
	p.RenderJSON(jobStatus)
}

func (p *JobController) JobExists() {
	projectName := p.GetString("project_name")
	p.ResolveProjectMember(projectName)
	jobName := p.GetString("job_name")
	isJobExists, err := service.JobExists(jobName, projectName)
	if err != nil {
		p.InternalError(err)
		logs.Error("Check job name failed, error: %+v", err.Error())
		return
	}
	if isJobExists {
		p.CustomAbortAudit(http.StatusConflict, jobNameDuplicateErr.Error())
	}
}

func (p *JobController) GetJobPodsAction() {
	var err error
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(j.ProjectID)
	pods, err := service.GetK8sJobPods(j)
	if err != nil {
		p.ParseError(err, c.ParseGetK8sError)
		return
	}
	p.RenderJSON(pods)
}

func (p *JobController) GetJobLogsAction() {
	var err error
	j, err := p.resolveJobInfo()
	if err != nil {
		return
	}
	//Judge authority
	p.ResolveUserPrivilegeByID(j.ProjectID)
	podName := p.Ctx.Input.Param(":podname")
	readCloser, err := service.GetK8sJobLogs(j, podName, p.generatePodLogOptions())
	if err != nil {
		p.ParseError(err, c.ParseGetK8sError)
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
	p.ResolveProjectMember(projectName)
	logs.Info("Get selectable job list for", projectName)
	jobList, err := service.GetJobsByProjectName(projectName)
	if err != nil {
		logs.Error("Failed to get selectable jobs.")
		p.InternalError(err)
		return
	}
	p.RenderJSON(jobList)
}
