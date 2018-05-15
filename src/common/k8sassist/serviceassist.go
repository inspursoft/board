// a temp file for building and guiding
package k8sassist

import (
	//api "k8s.io/client-go/pkg/api"
	//v1 "k8s.io/client-go/pkg/api/v1"
	//watch "k8s.io/client-go/pkg/watch"
	//rest "k8s.io/client-go/rest"
	"git/inspursoft/board/src/common/model"
)

// ServiceCli interface has methods to work with Service resources in k8s-assist.
// How to:  serviceCli, err := k8sassist.NewServices(nameSpace)
//          service, err := serviceCli.Get(serviceName)
type ServiceCliInterface interface {
	Create(*model.Service) (*model.Service, error)
	Update(*model.Service) (*model.Service, error)
	UpdateStatus(*model.Service) (*model.Service, error)
	Delete(name string) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*model.Service, error)
	List() (*model.ServiceList, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.Service, err error)
}
