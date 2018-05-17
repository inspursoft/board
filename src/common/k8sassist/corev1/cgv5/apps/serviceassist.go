package apps

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"
	"time"
)

type ServiceClient struct{}

func (c *ServiceClient) marshal(modelNamespace *model.Namespace) *types.Namespace {
	return &types.Namespace{
		TypeMeta: types.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: types.ObjectMeta{
			Name: modelNamespace.Name,
		},
	}
}

func (c *ServiceClient) unmarshal(typesNamespace *types.Namespace) *model.Namespace {
	ctime, _ := time.Parse(time.RFC3339, typesNamespace.ObjectMeta.CreationTimestamp.Format(time.RFC3339))
	return &model.Namespace{
		ObjectMeta: model.ObjectMeta{
			Name:              typesNamespace.ObjectMeta.Name,
			CreationTimestamp: ctime,
		},
		NamespacePhase: string(typesNamespace.Status.Phase),
	}
}

//create by struct or file
func (c *ServiceClient) Create(namespace *model.Namespace) (*model.Namespace, error) {
	// baseClient, err := NewBaseClient("")
	// if err != nil {
	// 	return nil, err
	// }
	// ns, err := baseClient.CoreV1().Namespaces().Create(c.marshal(namespace))
	// if err != nil {
	// 	logs.Error("Create namespace of %s failed. Err:%s", ns.Name, err.Error())
	// 	return nil, err
	// }

	// return c.unmarshal(ns), nil
	return nil, nil
}

// func (c *ClientSet) CreateService(service *types.Service) error {
// 	result, err := c.Client.CoreV1().Services(service.Namespace).Create(service)
// 	if err != nil {
// 		return err
// 	}
// 	logs.Debug("Created service %q.", result.GetObjectMeta().GetName())
// 	return nil
// }

// func (c *ClientSet) GetService(serviceName, serviceNamespace string) (*types.Service, error) {
// 	result, err := c.Client.CoreV1().Services(serviceNamespace).Get(serviceName, types.GetOptions{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	logs.Debug("Get service %q.", result.GetObjectMeta().GetName())
// 	return result, nil
// }

// func (c *ClientSet) GetServiceList(serviceNamespace string) (*types.ServiceList, error) {
// 	result, err := c.Client.CoreV1().Services(serviceNamespace).List(types.ListOptions{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	logs.Debug("Get service list.")
// 	return result, nil
// }

// func (c *ClientSet) DeleteService(serviceName, serviceNamespace string) error {
// 	err := c.Client.CoreV1().Services(serviceNamespace).Delete(serviceName, &types.DeleteOptions{})
// 	if err != nil {
// 		return err
// 	}
// 	logs.Debug("Delete service %q.", serviceName)
// 	return nil
// }
