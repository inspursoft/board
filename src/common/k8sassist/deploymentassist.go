package k8sassist

import (
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type deployments struct {
	namespace string
	client    *kubernetes.Clientset
}

func (d *deployments) Create(deployment *model.Deployment) (*model.Deployment, error) {
	k8sDep := toK8sDeployment(deployment)
	k8sDep, err := d.client.AppsV1beta2().Deployments(d.namespace).Create(k8sDep)
	if err != nil {
		logs.Error("Create deployment of %s/%s failed. Err:%+v", deployment.Name, d.namespace, err)
		return nil, err
	}

	modelDep := fromK8sDeployment(k8sDep)
	return modelDep, nil
}

func (d *deployments) Update(deployment *model.Deployment) (*model.Deployment, error) {
	k8sDep := toK8sDeployment(deployment)
	k8sDep, err := d.client.AppsV1beta2().Deployments(d.namespace).Update(k8sDep)
	if err != nil {
		logs.Error("update deployment of %s/%s failed. Err:%+v", deployment.Name, d.namespace, err)
		return nil, err
	}

	modelDep := fromK8sDeployment(k8sDep)
	return modelDep, nil
}

func (d *deployments) UpdateStatus(deployment *model.Deployment) (*model.Deployment, error) {
	k8sDep := toK8sDeployment(deployment)
	k8sDep, err := d.client.AppsV1beta2().Deployments(d.namespace).UpdateStatus(k8sDep)
	if err != nil {
		logs.Error("update deployment status of %s/%s failed. Err:%+v", deployment.Name, d.namespace, err)
		return nil, err
	}

	modelDep := fromK8sDeployment(k8sDep)
	return modelDep, nil
}

func (d *deployments) Delete(name string) error {
	err := d.client.AppsV1beta2().Deployments(d.namespace).Delete(name, nil)
	if err != nil {
		logs.Error("delete deployment of %s/%s failed. Err:%+v", name, d.namespace, err)
	}
	return err
}

func (d *deployments) Get(name string) (*model.Deployment, error) {
	deployment, err := d.client.AppsV1beta2().Deployments(d.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("get deployment of %s/%s failed. Err:%+v", name, d.namespace, err)
		return nil, err
	}

	modelDep := fromK8sDeployment(deployment)
	return modelDep, nil
}

func (d *deployments) List() (*model.DeploymentList, error) {
	deploymentList, err := d.client.AppsV1beta2().Deployments(d.namespace).List(metav1.ListOptions{})
	if err != nil {
		logs.Error("list deployments failed. Err:%+v", err)
		return nil, err
	}

	modelDepList := fromK8sDeploymentList(deploymentList)
	return modelDepList, nil
}

var _ DeploymentCliInterface = &deployments{}

func NewDeployments(namespace string) DeploymentCliInterface {
	//TODO: init the clientset.
	var client *kubernetes.Clientset
	return &deployments{
		namespace: namespace,
		client:    client,
	}
}

// DeploymentCli has methods to work with Deployment resources in k8s-assist.
// How to:  deploymentCli, err := k8sassist.NewDeployments(nameSpace)
//          _, err := deploymentCli.Update(&deployment)
type DeploymentCliInterface interface {
	Create(*model.Deployment) (*model.Deployment, error)
	Update(*model.Deployment) (*model.Deployment, error)
	UpdateStatus(*model.Deployment) (*model.Deployment, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.Deployment, error)
	//List(opts v1.ListOptions) (*DeploymentList, error)
	List() (*model.DeploymentList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1beta1.Deployment, err error)
}
