package apps

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"
	"time"

	"github.com/astaxie/beego/logs"
)

type NamespaceClient struct {
	Namespace types.NamespaceInterface
}

func (c *NamespaceClient) marshal(modelNamespace *model.Namespace) *types.Namespace {
	ns := &types.Namespace{
		TypeMeta: types.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: types.ObjectMeta{
			Name:      modelNamespace.Name,
			Namespace: modelNamespace.Namespace,
			Labels:    modelNamespace.Labels,
		},
	}

	return ns
}

func (c *NamespaceClient) unmarshal(typesNamespace *types.Namespace) *model.Namespace {
	ctime, _ := time.Parse(time.RFC3339, typesNamespace.ObjectMeta.CreationTimestamp.Format(time.RFC3339))
	return &model.Namespace{
		ObjectMeta: model.ObjectMeta{
			Name:              typesNamespace.ObjectMeta.Name,
			CreationTimestamp: ctime,
		},
		NamespacePhase: string(typesNamespace.Status.Phase),
	}
}

func (c *NamespaceClient) unmarshalList(typesNamespaceList *types.NamespaceList) *model.NamespaceList {
	modelNamespaceList := &model.NamespaceList{
		Items: make([]model.Namespace, 0),
	}
	for _, ns := range typesNamespaceList.Items {
		ctime, _ := time.Parse(time.RFC3339, ns.ObjectMeta.CreationTimestamp.Format(time.RFC3339))
		modelNamespaceList.Items = append(modelNamespaceList.Items, model.Namespace{
			ObjectMeta: model.ObjectMeta{
				Name:              ns.ObjectMeta.Name,
				CreationTimestamp: ctime,
			},
			NamespacePhase: string(ns.Status.Phase),
		})
	}

	return modelNamespaceList
}

func (c *NamespaceClient) Create(namespace *model.Namespace) (*model.Namespace, error) {
	ns, err := c.Namespace.Create(c.marshal(namespace))
	if err != nil {
		logs.Error("Create namespace of %s failed. Err:%s", ns.Name, err.Error())
		return nil, err
	}

	return c.unmarshal(ns), nil
}

func (c *NamespaceClient) Delete(namespaceName string) error {
	err := c.Namespace.Delete(namespaceName, &types.DeleteOptions{})
	if err != nil {
		logs.Error("Delete namespace of %s failed.", namespaceName)
		return err
	}

	return nil
}

func (c *NamespaceClient) Get(namespaceName string) (*model.Namespace, error) {
	ns, err := c.Namespace.Get(namespaceName, types.GetOptions{})
	if err != nil {
		logs.Error("Get namespace of %s failed.", namespaceName)
		return nil, err
	}

	return c.unmarshal(ns), nil
}

func (c *NamespaceClient) List() (*model.NamespaceList, error) {
	nsList, err := c.Namespace.List(types.ListOptions{})
	if err != nil {
		logs.Error("Get namespace list failed.")
		return nil, err
	}

	return c.unmarshalList(nsList), nil
}

func (c *NamespaceClient) Update(namespace *model.Namespace) (*model.Namespace, error) {
	ns, err := c.Namespace.Update(c.marshal(namespace))
	if err != nil {
		logs.Error("Update namespace of %s failed. Err:%s", ns.Name, err.Error())
		return nil, err
	}

	return c.unmarshal(ns), nil
}
