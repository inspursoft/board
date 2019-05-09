package service

import (
	"errors"
	"io"
	"strings"

	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/k8sassist"
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"
	"github.com/astaxie/beego/logs"
)

type JobDeployInfo struct {
	Job         *model.Job
	JobFileInfo []byte
}

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
func GetK8sJobPods(job *model.JobStatusMO) ([]model.PodMO, error) {
	logs.Debug("Get Job pods %s/%s", job.ProjectName, job.Name)

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	k8sjob, _, err := k8sclient.AppV1().Job(job.ProjectName).Get(job.Name)

	if err != nil {
		return nil, err
	}
	var opts model.ListOptions
	if k8sjob.Spec.Selector != nil {
		opts.LabelSelector = types.LabelSelectorToString(k8sjob.Spec.Selector)
	}
	list, err := k8sclient.AppV1().Pod(job.ProjectName).List(opts)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, nil
	}
	var pods []model.PodMO
	for i := range list.Items {
		var containers []model.ContainerMO
		for j := range list.Items[i].Spec.Containers {
			containers = append(containers, model.ContainerMO{
				Name:  list.Items[i].Spec.Containers[j].Name,
				Image: list.Items[i].Spec.Containers[j].Image,
			})
		}
		pods = append(pods, model.PodMO{
			Name:        list.Items[i].Name,
			ProjectName: list.Items[i].Namespace,
			Spec: model.PodSpecMO{
				Containers: containers,
			},
		})
	}
	return pods, nil
}

func GetK8sJobLogs(job *model.JobStatusMO, podName string, opt *model.PodLogOptions) (io.ReadCloser, error) {
	logs.Debug("Get Job logs %s/%s/%s", job.ProjectName, job.Name, podName)

	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	return k8sclient.AppV1().Pod(job.ProjectName).GetLogs(podName, opt)
}
