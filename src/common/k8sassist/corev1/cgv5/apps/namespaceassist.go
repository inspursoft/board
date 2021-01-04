package apps

import (
	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
)

type namespaces struct {
	Namespace types.NamespaceInterface
}

func (c *namespaces) Create(namespace *model.Namespace) (*model.Namespace, error) {
	ns, err := c.Namespace.Create(types.ToK8sNamespace(namespace))
	if err != nil {
		logs.Error("Create namespace of %s failed. Err:%s", ns.Name, err.Error())
		return nil, err
	}

	return types.FromK8sNamespace(ns), nil
}

func (c *namespaces) Delete(namespaceName string) error {
	err := c.Namespace.Delete(namespaceName, &types.DeleteOptions{})
	if err != nil {
		logs.Error("Delete namespace of %s failed.", namespaceName)
		return err
	}

	return nil
}

func (c *namespaces) Get(namespaceName string) (*model.Namespace, error) {
	ns, err := c.Namespace.Get(namespaceName, types.GetOptions{})
	if err != nil {
		logs.Error("Get namespace of %s failed.", namespaceName)
		return nil, err
	}

	return types.FromK8sNamespace(ns), nil
}

func (c *namespaces) List() (*model.NamespaceList, error) {
	nsList, err := c.Namespace.List(types.ListOptions{})
	if err != nil {
		logs.Error("Get namespace list failed.")
		return nil, err
	}

	return types.FromK8sNamespaceList(nsList), nil
}

func (c *namespaces) Update(namespace *model.Namespace) (*model.Namespace, error) {
	ns, err := c.Namespace.Update(types.ToK8sNamespace(namespace))
	if err != nil {
		logs.Error("Update namespace of %s failed. Err:%s", ns.Name, err.Error())
		return nil, err
	}

	return types.FromK8sNamespace(ns), nil
}

func NewNamespaces(namespace types.NamespaceInterface) *namespaces {
	return &namespaces{
		Namespace: namespace,
	}
}
