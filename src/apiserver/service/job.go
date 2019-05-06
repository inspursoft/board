package service

import (
	"errors"
	"io"
	"os"
	"strings"

	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"github.com/astaxie/beego/logs"
)

type JobDeployInfo struct {
	Job         *model.Job
	JobFileInfo []byte
}

//func InitServiceConfig() (*model.ServiceConfig, error) {
//	return &model.ServiceConfig{}, nil
//}
//
//func SelectProject(config *model.ServiceConfig, projectID int64) (*model.ServiceConfig, error) {
//	config.Phase = "SELECT_PROJECT"
//	config.ProjectID = projectID
//	return config, nil
//}
//
//func ConfigureContainers(config *model.ServiceConfig, containers []yaml.Container) (*model.ServiceConfig, error) {
//	config.Phase = "CONFIGURE_CONTAINERS"
//	config.DeploymentYaml = yaml.Deployment{}
//	config.DeploymentYaml.ContainerList = containers
//	return config, nil
//}
//
//func ConfigureService(config *model.ServiceConfig, service yaml.Service, deployment yaml.Deployment) (*model.ServiceConfig, error) {
//	config.Phase = "CONFIGURE_SERVICE"
//	config.ServiceYaml = service
//	config.DeploymentYaml = deployment
//	return config, nil
//}
//
//func ConfigureTest(config *model.ServiceConfig) error {
//	config.Phase = "CONFIGURE_TESTING"
//	return nil
//}
//
//func Deploy(config *model.ServiceConfig) error {
//	config.Phase = "CONFIGURE_DEPLOY"
//	return nil
//}
//
func CreateJob(jobConfig model.JobStatusMO) (*model.JobStatusMO, error) {
	query := model.Project{Name: jobConfig.ProjectName}
	project, err := GetProject(query, "name")
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project is invalid")
	}

	jobConfig.ProjectID = project.ID
	jobID, err := dao.AddJob(jobConfig)
	if err != nil {
		return nil, err
	}
	jobConfig.ID = jobID
	return &jobConfig, err
}

func UpdateJob(j model.JobStatusMO, fieldNames ...string) (bool, error) {
	if j.ID == 0 {
		return false, errors.New("no Job ID provided")
	}
	_, err := dao.UpdateJob(j, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdateJobStatus(jobID int64, status int) (bool, error) {
	return UpdateJob(model.JobStatusMO{ID: jobID, Status: status, Deleted: 0}, "status", "deleted")
}

func DeleteJobByID(jobID int64) (int64, error) {
	if jobID == 0 {
		return 0, errors.New("no Job ID provided")
	}
	num, err := dao.DeleteJob(model.JobStatusMO{ID: jobID})
	if err != nil {
		return 0, err
	}
	return num, nil
}

func GetJobList(name string, userID int64) ([]*model.JobStatusMO, error) {
	query := model.JobStatusMO{Name: name}
	jobList, err := dao.GetJobData(query, userID)
	if err != nil {
		return nil, err
	}
	return jobList, err
}

func GetPaginatedJobList(name string, userID int64, pageIndex int, pageSize int, orderField string, orderAsc int) (*model.PaginatedJobStatus, error) {
	query := model.JobStatusMO{Name: name}
	paginatedJobStatus, err := dao.GetPaginatedJobData(query, userID, pageIndex, pageSize, orderField, orderAsc)
	if err != nil {
		return nil, err
	}
	return paginatedJobStatus, nil
}

func DeleteJob(jobID int64) (bool, error) {
	s := model.JobStatusMO{ID: jobID}
	_, err := dao.DeleteJob(s)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetK8sJob(jobURL string) (*model.Job, error) {
	var job model.Job
	logs.Debug("Get Job info jobURL(job): %+s", jobURL)
	err := k8sGet(&job, jobURL)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func GetK8sJobByK8sassist(pName string, jName string) (*model.Job, error) {
	logs.Debug("Get Job info %s/%s", pName, jName)

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	job, _, err := k8sclient.AppV1().Job(pName).Get(jName)

	if err != nil {
		return nil, err
	}
	return job, nil
}

func GetJob(job model.JobStatusMO, selectedFields ...string) (*model.JobStatusMO, error) {
	j, err := dao.GetJob(job, selectedFields...)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func GetJobByID(jobID int64) (*model.JobStatusMO, error) {
	return GetJob(model.JobStatusMO{ID: jobID, Deleted: 0}, "id", "deleted")
}

func GetJobByProject(jobName string, projectName string) (*model.JobStatusMO, error) {
	var jobquery model.JobStatusMO
	jobquery.Name = jobName
	jobquery.ProjectName = projectName
	job, err := GetJob(jobquery, "name", "project_name")
	if err != nil {
		return nil, err
	}
	return job, nil
}

//func SyncJobWithK8s(pName string) error {
//	logs.Debug("Sync Job of namespace %s", pName)
//	//obtain serviceList data of
//	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
//		KubeConfigPath: kubeConfigPath(),
//	})
//
//	jobList, err := k8sclient.AppV1().Job(pName).List()
//	if err != nil {
//		logs.Error("Failed to get job list with project name: %s", pName)
//		return err
//	}
//
//	//handle the jobList data
//	var jobquery model.JobStatus
//	for _, item := range jobList.Items {
//		project, err := GetProjectByName(item.Namespace)
//		if err != nil {
//			logs.Error("Failed to check project in DB %s", item.Namespace)
//			return err
//		}
//		if project == nil {
//			logs.Error("not found project in DB: %s", item.Namespace)
//			continue
//		}
//		if item.ObjectMeta.Name == k8sService {
//			continue
//		}
//		jobquery.Name = item.ObjectMeta.Name
//		jobquery.OwnerID = int64(project.OwnerID) //owner or admin TBD
//		jobquery.OwnerName = project.OwnerName
//		jobquery.ProjectName = project.Name
//		jobquery.ProjectID = project.ID
//		jobquery.Public = defaultPublic
//		jobquery.Comment = defaultComment
//		jobquery.Deleted = defaultDeleted
//		jobquery.Status = defaultStatus
//		jobquery.Source = k8s
//		jobquery.CreationTime, _ = time.Parse(time.RFC3339, item.CreationTimestamp.Format(time.RFC3339))
//		jobquery.UpdateTime, _ = time.Parse(time.RFC3339, item.CreationTimestamp.Format(time.RFC3339))
//		_, err = dao.SyncServiceData(servicequery)
//		if err != nil {
//			logs.Error("Sync Service %s failed.", servicequery.Name)
//		}
//	}
//
//	return nil
//}

func GetJobsByProjectName(pname string) ([]model.JobStatusMO, error) {
	jobList, err := dao.GetJobs("project_name", pname)
	if err != nil {
		return nil, err
	}
	return jobList, err
}

func JobExists(jobName string, projectName string) (bool, error) {
	var jobquery model.JobStatusMO
	jobquery.Name = jobName
	jobquery.ProjectName = projectName
	s, err := GetJob(jobquery, "name", "project_name")

	return s != nil, err
}

func GenerateJobYamlFileFromK8s(job *model.JobStatusMO, loadPath, filename string) error {
	clusterConfig := &k8sassist.K8sAssistConfig{KubeConfigPath: kubeConfigPath()}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)
	_, jobFileInfo, err := cli.AppV1().Job(job.ProjectName).Get(job.Name)
	if err != nil {
		return err
	}
	return utils.GenerateFile(jobFileInfo, loadPath, filename)
}

func SaveJobDeployYamlFiles(jobInfo *JobDeployInfo, loadPath, filename string) error {
	if jobInfo == nil {
		logs.Error("Job Deploy info is empty.")
		return errors.New("Job Deploy info is empty.")
	}
	return utils.GenerateFile(jobInfo.JobFileInfo, loadPath, filename)
}

//check yaml file config
func CheckJobYamlConfig(yamlFile io.Reader, projectName string) (*model.Job, error) {
	clusterConfig := &k8sassist.K8sAssistConfig{KubeConfigPath: kubeConfigPath()}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)

	jobInfo, err := cli.AppV1().Job(projectName).CheckYaml(yamlFile)
	if err != nil {
		logs.Error("Check job object by job.yaml failed, err:%+v\n", err)
		return nil, err
	}

	return jobInfo, nil
}

func DeployJob(jobConfig *model.JobConfig, registryURI string) (*JobDeployInfo, error) {
	clusterConfig := &k8sassist.K8sAssistConfig{KubeConfigPath: kubeConfigPath()}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)
	logs.Debug("Created job: ", jobConfig.Name)
	k8sjobConfig := MarshalJob(jobConfig, registryURI)
	//logs.Debug("Marshaled deployment: ", deploymentConfig)
	jobInfo, jobFileInfo, err := cli.AppV1().Job(jobConfig.ProjectName).Create(k8sjobConfig)
	if err != nil {
		logs.Error("Deploy job object of %s failed. error: %+v\n", jobConfig.Name, err)
		return nil, err
	}

	return &JobDeployInfo{
		Job:         jobInfo,
		JobFileInfo: jobFileInfo,
	}, nil
}

func DeployJobByYaml(projectName, jobAbsName string) error {
	clusterConfig := &k8sassist.K8sAssistConfig{KubeConfigPath: kubeConfigPath()}
	cli := k8sassist.NewK8sAssistClient(clusterConfig)

	jobFile, err := os.Open(jobAbsName)
	if err != nil {
		return err
	}
	defer jobFile.Close()

	_, err = cli.AppV1().Job(projectName).CreateByYaml(jobFile)
	if err != nil {
		logs.Error("Deploy job object by job.yaml failed, err:%+v\n", err)
		return err
	}
	return nil
}

func StopJobK8s(j *model.JobStatusMO) error {
	logs.Info("stop job in cluster %s", j.Name)
	// Stop deployment
	config := k8sassist.K8sAssistConfig{}
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	d := k8sclient.AppV1().Job(j.ProjectName)
	err := d.Delete(j.Name)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		logs.Error("Failed to delete job in cluster, error:%v", err)
		return err
	}
	return nil
}

func MarshalJob(jobConfig *model.JobConfig, registryURI string) *model.Job {
	if jobConfig == nil {
		return nil
	}
	podTemplate := model.PodTemplateSpec{
		ObjectMeta: model.ObjectMeta{
			Name:   jobConfig.Name,
			Labels: map[string]string{"app": jobConfig.Name},
		},
		Spec: model.PodSpec{
			Volumes:       setDeploymentVolumes(jobConfig.ContainerList),
			Containers:    setDeploymentContainers(jobConfig.ContainerList, registryURI),
			NodeSelector:  setDeploymentNodeSelector(jobConfig.NodeSelector),
			Affinity:      setDeploymentAffinity(jobConfig.AffinityList),
			RestartPolicy: model.RestartPolicyNever,
		},
	}

	return &model.Job{
		ObjectMeta: model.ObjectMeta{
			Name:      jobConfig.Name,
			Namespace: jobConfig.ProjectName,
		},
		Spec: model.JobSpec{
			Parallelism:           jobConfig.Parallelism,
			Completions:           jobConfig.Completions,
			ActiveDeadlineSeconds: jobConfig.ActiveDeadlineSeconds,
			BackoffLimit:          jobConfig.BackoffLimit,

			Template: podTemplate,
		},
	}
}
