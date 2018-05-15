package k8sassist

import (
	"git/inspursoft/board/src/common/model"
)

type deployments struct {
	ns string
}

func (d *deployments) Create(*model.Deployment) (*model.Deployment, error) {
	return nil, nil
}

func (d *deployments) Update(*model.Deployment) (*model.Deployment, error) {
	return nil, nil
}

func (d *deployments) UpdateStatus(*model.Deployment) (*model.Deployment, error) {
	return nil, nil
}

func (d *deployments) Delete(name string) error {
	return nil
}

func (d *deployments) Get(name string) (*model.Deployment, error) {
	return nil, nil
}

func (d *deployments) List() (*model.DeploymentList, error) {
	return nil, nil
}

var _ DeploymentCliInterface = &deployments{}

// newNodes returns a Nodes
func NewDeployments(namespace string) (*deployments, error) {
	return &deployments{ns: namespace}, nil
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
