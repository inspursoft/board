// a temp file for building and guiding
package k8sassist

import (
	"git/inspursoft/board/src/common/model"
	//"git/inspursoft/k8sassist/k8sassist/v1/adapter/alphe/types"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/astaxie/beego/logs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cliv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type NamespaceClient struct {
	namespace cliv1.NamespaceInterface
}

func (c *NamespaceClient) marshal(modelNamespace *model.Namespace) *v1.Namespace {
	return &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: modelNamespace.Name,
		},
	}
}

func (c *NamespaceClient) unmarshal(typesNamespace *v1.Namespace) *model.Namespace {
	ctime, _ := time.Parse(time.RFC3339, typesNamespace.ObjectMeta.CreationTimestamp.Format(time.RFC3339))
	return &model.Namespace{
		ObjectMeta: model.ObjectMeta{
			Name:              typesNamespace.ObjectMeta.Name,
			CreationTimestamp: ctime,
		},
		NamespacePhase: string(typesNamespace.Status.Phase),
	}
}

func (c *NamespaceClient) Create(namespace *model.Namespace) (*model.Namespace, error) {
	ns, err := c.namespace.Create(c.marshal(namespace))
	if err != nil {
		logs.Error("Create namespace of %s failed. Err:%s", ns.Name, err.Error())
		return nil, err
	}

	return c.unmarshal(ns), nil
}

func (c *NamespaceClient) Delete(namespace *model.Namespace) error {

	return nil
}

func (c *NamespaceClient) Get(name string) (*model.Namespace, error) {

	return nil, nil
}

// NamespaceCli Interface has methods to work with Namespace resources.
type NamespaceCliInterface interface {
	Create(*model.Namespace) (*model.Namespace, error)
	Delete(*model.Namespace) error
	Get(name string) (*model.Namespace, error)
	List() (*model.NamespaceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Namespace, err error)
}
