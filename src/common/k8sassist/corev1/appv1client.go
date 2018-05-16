package v1

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/apps"
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"
)

type AppV1Client struct {
	Clientset *types.Clientset
}

type AppV1ClientInterface interface {
	Service() ServiceClientInterface
	Deployment() DeploymentClientInterface
	Node() NodeClientInterface
	Namespace() NamespaceClientInterface
}

func NewAppV1Client(clientset *types.Clientset) *AppV1Client {
	return &AppV1Client{
		Clientset: clientset,
	}
}

func (p *AppV1Client) Namespace() NamespaceClientInterface {
	return &apps.NamespaceClient{Namespace: p.Clientset.CoreV1().Namespaces()}
}

func (p *AppV1Client) Node() NodeClientInterface {
	return nil
}

func (p *AppV1Client) Deployment() DeploymentClientInterface {
	return nil
}

func (p *AppV1Client) Service() ServiceClientInterface {
	return nil
}

// NamespaceCli Interface has methods to work with Namespace resources.
type ServiceClientInterface interface {
	Create() (*model.Namespace, error)
	//Delete(*model.Namespace) error
	//Get(name string) (*model.Namespace, error)
	// List() (*model.NamespaceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Namespace, err error)
}

// NamespaceCli Interface has methods to work with Namespace resources.
type DeploymentClientInterface interface {
	Create(*model.Namespace) (*model.Namespace, error)
	//Delete(*model.Namespace) error
	//Get(name string) (*model.Namespace, error)
	// List() (*model.NamespaceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Namespace, err error)
}

// NamespaceCli Interface has methods to work with Namespace resources.
type NodeClientInterface interface {
	Create(*model.Namespace) (*model.Namespace, error)
	// Delete(*model.Namespace) error
	// Get(name string) (*model.Namespace, error)
	// List() (*model.NamespaceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Namespace, err error)
}

// NamespaceCli Interface has methods to work with Namespace resources.
type NamespaceClientInterface interface {
	Create(*model.Namespace) (*model.Namespace, error)
	Update(*model.Namespace) (*model.Namespace, error)
	Delete(string) error
	Get(string) (*model.Namespace, error)
	List() (*model.NamespaceList, error)
}
