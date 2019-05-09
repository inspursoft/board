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

//
func syncJobK8sStatus(jobList []*model.JobStatusMO) error {
	var err error
	reason := ""
	status := uncompleted
	// synchronize job status with the cluster system
	for _, jobStatusMO := range jobList {
		// Check the job status
		job, err := service.GetK8sJobByK8sassist(jobStatusMO.ProjectName, jobStatusMO.Name)
		if job == nil {
			logs.Info("Failed to get job", err)
			reason = "The job is not established in cluster system"
			status = uncompleted
		} else if job.Status.CompletionTime == nil {
			logs.Info("The job does not complete")
			reason = "The job does not complete"
			status = uncompleted
		} else {
			// read the doc https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
			success := false
			if job.Spec.Completions == nil {
				if job.Spec.Parallelism == nil {
					success = job.Status.Succeeded >= 1
				} else {
					success = job.Status.Succeeded >= *job.Spec.Parallelism
				}
			} else {
				success = job.Status.Succeeded >= *job.Spec.Completions
			}
			if success {
				logs.Debug("The desired completion number is reached",
					job.Status.Succeeded, job.Spec.Completions)
				status = completed
				reason = "The desired replicas number is reached"
			} else {
				logs.Debug("The desired completion number is not reached",
					job.Status.Succeeded, job.Spec.Completions)
				status = failed
				reason = "The desired replicas number is not reached"
			}
		}
		jobStatusMO.Status = status
		jobStatusMO.Comment = "Reason: " + reason
		_, err = service.UpdateJob(*jobStatusMO, "status", "comment")
		if err != nil {
			logs.Error("Failed to update job status.")
			break
		}
	}
	return err
}

//get job list
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
		err = syncJobK8sStatus(jobStatus)
		if err != nil {
			p.internalError(err)
			return
		}
		p.renderJSON(jobStatus)
	} else {
		paginatedJobStatus, err := service.GetPaginatedJobList(jobName, p.currentUser.ID, pageIndex, pageSize, orderField, orderAsc)
		if err != nil {
			p.internalError(err)
			return
		}
		err = syncJobK8sStatus(paginatedJobStatus.JobStatusList)
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
	_ = p.resolveUserPrivilegeByID(j.ProjectID)

	jobStatus := []*model.JobStatusMO{j}
	err = syncJobK8sStatus(jobStatus)
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
		return
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
		return
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
	opt.Previous, err = p.GetBool("privious", false)
	if err != nil {
		logs.Warn("Privious parameter %s is invalid: %+v", p.GetString("privious"), err)
	}
	opt.Timestamps, err = p.GetBool("timestamps", false)
	if err != nil {
		logs.Warn("Timestamps parameter %s is invalid: %+v", p.GetString("timestamps"), err)
	}

	if p.GetString("sinceSeconds") != "" {
		since, err := p.GetInt64("sinceSeconds")
		if err != nil {
			logs.Warn("SinceSeconds parameter %s is invalid: %+v", p.GetString("sinceSeconds"), err)
		} else {
			opt.SinceSeconds = &since
		}
	}

	since := p.GetString("sinceTime")
	if since != "" {
		sinceTime, err := time.Parse(time.RFC3339, since)
		if err != nil {
			logs.Warn("SinceTime parameter %s is invalid: %+v", since, err)
		} else {
			opt.SinceTime = &sinceTime
		}
	}

	tail, err := p.GetInt64("tailLines", -1)
	if err != nil {
		logs.Warn("TailLines parameter %s is invalid: %+v", p.GetString("tailLines"), err)
	} else if tail != -1 {
		opt.TailLines = &tail
	}

	limit, err := p.GetInt64("limitBytes", -1)
	if err != nil {
		logs.Warn("LimitBytes parameter %s is invalid: %+v", p.GetString("limitBytes"), err)
	} else if limit != -1 {
		opt.LimitBytes = &limit
	}

	return opt
}
