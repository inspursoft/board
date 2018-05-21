package apps

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"

	"io"
	"io/ioutil"

	"github.com/astaxie/beego/logs"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/typed/apps/v1beta2"
)

type deployments struct {
	namespace string
	deploy    v1beta2.DeploymentInterface
}

func (d *deployments) Create(deployment *model.Deployment) (*model.Deployment, error) {
	k8sDep := types.ToK8sDeployment(deployment)
	k8sDep, err := d.deploy.Create(k8sDep)
	if err != nil {
		logs.Error("Create deployment of %s/%s failed. Err:%+v", deployment.Name, d.namespace, err)
		return nil, err
	}

	modelDep := types.FromK8sDeployment(k8sDep)
	return modelDep, nil
}

func (d *deployments) Update(deployment *model.Deployment) (*model.Deployment, error) {
	k8sDep := types.ToK8sDeployment(deployment)
	k8sDep, err := d.deploy.Update(k8sDep)
	if err != nil {
		logs.Error("Update deployment of %s/%s failed. Err:%+v", deployment.Name, d.namespace, err)
		return nil, err
	}

	modelDep := types.FromK8sDeployment(k8sDep)
	return modelDep, nil
}

func (d *deployments) UpdateStatus(deployment *model.Deployment) (*model.Deployment, error) {
	k8sDep := types.ToK8sDeployment(deployment)
	k8sDep, err := d.deploy.UpdateStatus(k8sDep)
	if err != nil {
		logs.Error("Update deployment status of %s/%s failed. Err:%+v", deployment.Name, d.namespace, err)
		return nil, err
	}

	modelDep := types.FromK8sDeployment(k8sDep)
	return modelDep, nil
}

func (d *deployments) Delete(name string) error {
	err := d.deploy.Delete(name, nil)
	if err != nil {
		logs.Error("Delete deployment of %s/%s failed. Err:%+v", name, d.namespace, err)
	}
	return err
}

func (d *deployments) Get(name string) (*model.Deployment, error) {
	deployment, err := d.deploy.Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("Get deployment of %s/%s failed. Err:%+v", name, d.namespace, err)
		return nil, err
	}

	modelDep := types.FromK8sDeployment(deployment)
	return modelDep, nil
}

func (d *deployments) List() (*model.DeploymentList, error) {
	deploymentList, err := d.deploy.List(metav1.ListOptions{})
	if err != nil {
		logs.Error("List deployments failed. Err:%+v", err)
		return nil, err
	}

	modelDepList := types.FromK8sDeploymentList(deploymentList)
	return modelDepList, nil
}

func (d *deployments) Patch(name string, pt model.PatchType, data []byte, subresources ...string) (result *model.Deployment, err error) {
	k8sDep, err := d.deploy.Patch(name, k8stypes.PatchType(string(pt)), data, subresources...)
	if err != nil {
		logs.Error("Patch deployment of %s/%s failed. Err:%+v", name, d.namespace, err)
		return nil, err
	}

	modelDep := types.FromK8sDeployment(k8sDep)
	return modelDep, nil
}

func (d *deployments) CreateByYaml(r io.Reader) (*model.Deployment, error) {
	context, err := ioutil.ReadAll(r)
	if err != nil {
		logs.Error("Read file failed, error: %v", err)
		return nil, err
	}

	var deployment types.Deployment
	err = yaml.Unmarshal(context, &deployment)
	if err != nil {
		logs.Error("Unmarshal deployment failed, error: %v", err)
		return nil, err
	}

	err = types.CheckDeploymentConfig(d.namespace, deployment)
	if err != nil {
		logs.Error("Deployment config is error, error: %v", err)
		return nil, err
	}

	deploymentInfo, err := d.deploy.Create(&deployment)
	if err != nil {
		logs.Error("Create deployment failed, error: %v", err)
		return nil, err
	}

	return types.FromK8sDeployment(deploymentInfo), nil
}

func NewDeployments(namespace string, deploy v1beta2.DeploymentInterface) *deployments {
	return &deployments{
		namespace: namespace,
		deploy:    deploy,
	}
}
