package apps

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"
	"time"
)

type DeploymentClient struct{}

func (c *DeploymentClient) marshal(namespace *model.Namespace) *types.Namespace {
	return &types.Namespace{
		TypeMeta: types.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: types.ObjectMeta{
			Name: namespace.ObjectMeta.Name,
		},
	}
}

func (c *DeploymentClient) unmarshal(namespace *types.Namespace) *model.Namespace {
	ctime, _ := time.Parse(time.RFC3339, namespace.ObjectMeta.CreationTimestamp.Format(time.RFC3339))
	return &model.Namespace{
		ObjectMeta: model.ObjectMeta{
			Name:              namespace.ObjectMeta.Name,
			CreationTimestamp: ctime,
		},
		NamespacePhase: string(namespace.Status.Phase),
	}
}

//create by struct or file
// func (c *DeploymentClient) Create(namespace *model.Namespace) (*model.Namespace, error) {
// 	baseClient, err := NewBaseClient("")
// 	if err != nil {
// 		return nil, err
// 	}
// 	ns, err := baseClient.CoreV1().Namespaces().Create(c.marshal(namespace))
// 	if err != nil {
// 		logs.Error("Create namespace of %s failed. Err:%s", ns.Name, err.Error())
// 		return nil, err
// 	}

// 	return c.unmarshal(ns), nil
// }

// func (c *ClientSet) CreateDeployment(deployment *types.Deployment) error {
// 	result, err := c.Client.AppsV1beta1().Deployments(deployment.Namespace).Create(deployment)
// 	if err != nil {
// 		return err
// 	}
// 	logs.Debug("Created deployment %q.", result.GetObjectMeta().GetName())
// 	return nil
// }

// func (c *ClientSet) UpdateDeployment(deployment *types.Deployment) error {
// 	result, err := c.Client.AppsV1beta1().Deployments(deployment.Namespace).Update(deployment)
// 	if err != nil {
// 		return err
// 	}
// 	logs.Debug("Update deployment %q.", result.GetObjectMeta().GetName())
// 	return nil
// }

// func (c *ClientSet) GetDeployment(deploymentName, deploymentNamespace string) (*types.Deployment, error) {
// 	result, err := c.Client.AppsV1beta1().Deployments(deploymentNamespace).Get(deploymentName, types.GetOptions{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	logs.Debug("Get deployment %q.", result.GetObjectMeta().GetName())
// 	return result, nil
// }

// func (c *ClientSet) GetDeploymentList(deploymentNamespace string) (*types.DeploymentList, error) {
// 	result, err := c.Client.AppsV1beta1().Deployments(deploymentNamespace).List(types.ListOptions{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	logs.Debug("Get deployment list.")
// 	return result, nil
// }

// func (c *ClientSet) DeleteDeployment(deploymentName, deploymentNamespace string) error {
// 	deletePolicy := types.DeletePropagationForeground
// 	err := c.Client.AppsV1beta1().Deployments(deploymentNamespace).Delete(deploymentName, &types.DeleteOptions{PropagationPolicy: &deletePolicy})
// 	if err != nil {
// 		return err
// 	}
// 	logs.Debug("Delete deployment %q.", deploymentName)
// 	return nil
// }
