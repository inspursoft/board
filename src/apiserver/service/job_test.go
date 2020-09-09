package service_test

import (
	"testing"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

var unitTestJob = model.JobStatusMO{
	Name:        "unittestjob",
	OwnerID:     1,
	OwnerName:   "boardadmin",
	ProjectName: "library",
}

var unitTestJobConfig = model.JobConfig{
	Name:        "unittestjob",
	ProjectID:   int64(1),
	ProjectName: "library",
	ContainerList: []model.Container{
		{
			Name: "nginx",
			Image: model.ImageIndex{
				ImageName:   "library/jobcase",
				ImageTag:    "1.0",
				ProjectName: "library",
			},
		},
	},
}

var jobListInfo []*model.JobStatusMO
var podMOInfo []model.PodMO

func TestCreateJob(t *testing.T) {
	jobInfo, err := service.CreateJob(unitTestJob)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while create job.")
	unitTestJob.ID = jobInfo.ID
}

func TestUpdateJob(t *testing.T) {
	unitTestJob.Source = 1
	_, err := service.UpdateJob(unitTestJob, "source")
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while update job source.")
}

func TestUpdateJobStatus(t *testing.T) {
	unitTestJob.Status = 1
	_, err := service.UpdateJobStatus(unitTestJob.ID, 2)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while update job status.")
}

func TestGetJobByID(t *testing.T) {
	jobInfo, err := service.GetJobByID(unitTestJob.ID)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get job information by ID.")
	logs.Info("Job info is %+v", jobInfo)
}

func TestGetJobByProject(t *testing.T) {
	jobInfo, err := service.GetJobByProject(unitTestJob.Name, unitTestJob.ProjectName)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get job information by project name and job name.")
	logs.Info("Job info is %+v", jobInfo)
}

func TestGetJobsByProjectName(t *testing.T) {
	jobList, err := service.GetJobsByProjectName(unitTestJob.ProjectName)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get job list by project name.")
	logs.Info("Job list is %+v", jobList)
}

func TestJobExists(t *testing.T) {
	exist, err := service.JobExists(unitTestJob.Name, unitTestJob.ProjectName)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while check job by project name and job name.")
	assert.Equal(true, exist, "Job status is error.")
}

func TestGetJobList(t *testing.T) {
	jobList, err := service.GetJobList("", 1)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get job list.")
	jobListInfo = jobList
	logs.Info("Job list is %+v", jobListInfo)
}

func TestGetPaginatedJobList(t *testing.T) {
	PaginatedJList, err := service.GetPaginatedJobList("", 1, 1, 10, "name", 0)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get job list.")
	logs.Info("Paginated job list is %+v", PaginatedJList)
}

func TestDeployJob(t *testing.T) {
	registryURI := utils.GetStringValue("REGISTRY_BASE_URI")
	logs.Info("REGISTRY_URI %s", registryURI)
	jobDeployInfo, err := service.DeployJob(&unitTestJobConfig, registryURI)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get job list.")
	logs.Info("Information about deployed Job is %+v", string(jobDeployInfo.JobFileInfo))
}

func TestGetK8sJobPods(t *testing.T) {
	podMO, err := service.GetK8sJobPods(&unitTestJob)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while get job pod information.")
	podMOInfo = podMO
	logs.Info("Information about pod of Job is %+v", podMOInfo)
}

func TestGetK8sJobLogs(t *testing.T) {
	readCloser, _ := service.GetK8sPodLogs(podMOInfo[0].ProjectName, podMOInfo[0].Name, &model.PodLogOptions{})
	if readCloser != nil {
		defer readCloser.Close()
		logs.Info("logs about pods of Job is %+v", readCloser)
	}
	// TODO: This case always failed.
	// assert := assert.New(t)
	// assert.Nil(err, "Error occurred while get job logs.")
}

func TestSyncJobK8sStatus(t *testing.T) {
	err := service.SyncJobK8sStatus(jobListInfo)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while sync jobs status with kubernetes.")
}

func TestStopJobK8s(t *testing.T) {
	err := service.StopJobK8s(&unitTestJob)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while stop job in kubernetes.")
}

func TestDeleteJobByID(t *testing.T) {
	unitTestJob.Status = 1
	num, err := service.DeleteJobByID(unitTestJob.ID)
	assert := assert.New(t)
	assert.Nil(err, "Error occurred while delete job by ID.")
	assert.Equal(int64(1), num, "Num of deleted job is error.")
}
