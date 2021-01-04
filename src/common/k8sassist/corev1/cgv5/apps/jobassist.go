package apps

import (
	"encoding/json"
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"
	"io"
	"io/ioutil"

	"github.com/astaxie/beego/logs"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	betchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

type job struct {
	namespace string
	job       betchv1.JobInterface
}

func NewJob(namespace string, jobIf betchv1.JobInterface) *job {
	return &job{
		namespace: namespace,
		job:       jobIf,
	}
}

func (j *job) processJobHandler(job *model.Job, handler func(*types.Job) (*types.Job, error)) (customModel *model.Job, primitiveData []byte, err error) {
	k8sJob := types.ToK8sJob(job)
	logs.Debug("handler k8s job: ", k8sJob)
	handledDep, err := handler(k8sJob)
	if err != nil {
		logs.Error("Failed to handle job of %s/%s failej. Err:%+v", handledDep.Name, handledDep.Namespace, err)
		return nil, nil, err
	}
	customModel = types.FromK8sJob(handledDep)
	primitiveData, err = yaml.Marshal(types.GenerateJobConfig(handledDep))
	if err != nil {
		logs.Error("Failed to marshal primitive from job config, error: %+v", err)
		return
	}
	return
}

func (j *job) Create(job *model.Job) (*model.Job, []byte, error) {
	return j.processJobHandler(job, j.job.Create)
}

func (j *job) Update(job *model.Job) (*model.Job, []byte, error) {
	return j.processJobHandler(job, j.job.Update)
}

func (j *job) Delete(name string) error {
	deletePolicy := types.DeletePropagationForeground
	err := j.job.Delete(name, &types.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logs.Error("Delete job of %s/%s failed. Err:%+v", name, j.namespace, err)
	}
	return err
}

func (j *job) Get(name string) (*model.Job, []byte, error) {
	job, err := j.job.Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("Get job of %s/%s failed. Err:%+v", name, j.namespace, err)
		return nil, nil, err
	}
	jobConfig := types.GenerateJobConfig(job)
	jobfileInfo, err := yaml.Marshal(jobConfig)
	if err != nil {
		logs.Error("Marshal job failed, error: %v", err)
		return nil, nil, err
	}
	modelDep := types.FromK8sJob(job)
	return modelDep, jobfileInfo, nil
}

func (j *job) List() (*model.JobList, error) {
	jobList, err := j.job.List(metav1.ListOptions{})
	if err != nil {
		logs.Error("List job failed. Err:%+v", err)
		return nil, err
	}
	modelJobList := types.FromK8sJobList(jobList)
	return modelJobList, nil
}

func (j *job) UpdateStatus(job *model.Job) (*model.Job, []byte, error) {
	return j.processJobHandler(job, j.job.UpdateStatus)
}

func (j *job) Patch(name string, pt model.PatchType, data []byte, subresources ...string) (result *model.Job, err error) {
	k8sJob, err := j.job.Patch(name, k8stypes.PatchType(string(pt)), data, subresources...)
	if err != nil {
		logs.Error("Patch job of %s/%s failed. Err:%+v", name, j.namespace, err)
		return nil, err
	}
	modelJob := types.FromK8sJob(k8sJob)
	return modelJob, nil
}

func (j *job) PatchToK8s(name string, pt model.PatchType, job *model.Job) (*model.Job, []byte, error) {
	k8sJob := types.ToK8sJob(job)
	serviceRollConfig, err := json.Marshal(k8sJob)
	if err != nil {
		logs.Debug("Marshal rollingUpdateConfig failed %+v\n", k8sJob)
		return nil, nil, err
	}

	k8sJob, err = j.job.Patch(name, k8stypes.PatchType(pt), serviceRollConfig)
	if err != nil {
		logs.Error("PatchK8s Job of %s/%s failej. Err:%+v", job.Name, j.namespace, err)
		return nil, nil, err
	}

	jobConfig := types.GenerateJobConfig(k8sJob)
	jobfileInfo, err := yaml.Marshal(jobConfig)
	if err != nil {
		logs.Error("Marshal job failed, error: %v", err)
		return nil, nil, err
	}
	modelDep := types.FromK8sJob(k8sJob)
	return modelDep, jobfileInfo, nil
}

func (j *job) CreateByYaml(r io.Reader) (*model.Job, error) {
	context, err := ioutil.ReadAll(r)
	if err != nil {
		logs.Error("Read file failed, error: %v", err)
		return nil, err
	}

	var job types.Job
	err = yaml.Unmarshal(context, &job)
	if err != nil {
		logs.Error("Unmarshal job failed, error: %v", err)
		return nil, err
	}

	if job.ObjectMeta.Namespace != j.namespace {
		logs.Error(namespacesErr.Error())
		return nil, namespacesErr
	}

	jobInfo, err := j.job.Create(&job)
	if err != nil {
		logs.Error("Create job failed, error: %v", err)
		return nil, err
	}

	return types.FromK8sJob(jobInfo), nil
}

func (j *job) CheckYaml(r io.Reader) (*model.Job, error) {
	context, err := ioutil.ReadAll(r)
	if err != nil {
		logs.Error("Read file failed, error: %v", err)
		return nil, err
	}

	var job types.Job
	err = yaml.Unmarshal(context, &job)
	if err != nil {
		logs.Error("Unmarshal job failed, error: %v", err)
		return nil, err
	}

	if job.ObjectMeta.Namespace != j.namespace {
		logs.Error(namespacesErr.Error())
		return nil, namespacesErr
	}

	return types.FromK8sJob(&job), nil
}
