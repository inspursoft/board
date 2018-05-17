package apps

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"
	"time"
)

type NodeClient struct{}

func (c *NodeClient) marshal(modelNamespace *model.Namespace) *types.Namespace {
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

func (c *NodeClient) unmarshal(typesNamespace *types.Namespace) *model.Namespace {
	ctime, _ := time.Parse(time.RFC3339, typesNamespace.ObjectMeta.CreationTimestamp.Format(time.RFC3339))
	return &model.Namespace{
		ObjectMeta: model.ObjectMeta{
			Name:              typesNamespace.ObjectMeta.Name,
			CreationTimestamp: ctime,
		},
		NamespacePhase: string(typesNamespace.Status.Phase),
	}
}

func (c *NodeClient) Create(namespace *model.Namespace) (*model.Namespace, error) {
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
