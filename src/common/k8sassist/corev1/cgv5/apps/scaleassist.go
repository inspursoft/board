// a temp file for building and guiding
package apps

import (
	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"

	k8sv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	rest "k8s.io/client-go/rest"
)

type scales struct {
	namespace string
	scale     ScaleInterface
}

func (s *scales) Update(kind string, scale *model.Scale) (*model.Scale, error) {
	k8sScale := types.ToK8sScale(scale)
	newk8sScale, err := s.scale.Update(kind, k8sScale)
	if err != nil {
		logs.Error("Update Scale of %s/%s failed. Err:%+v", scale.Name, s.namespace, err)
		return nil, err
	}

	modelScale := types.FromK8sScale(newk8sScale)
	return modelScale, nil
}

func (s *scales) Get(kind string, name string) (*model.Scale, error) {
	scaleinstance, err := s.scale.Get(kind, name)
	if err != nil {
		logs.Error("Get scale of %s failed. Err:%+v", name, err)
		return nil, err
	}

	return types.FromK8sScale(scaleinstance), nil
}

func NewScales(namespace string, scale ScaleInterface) *scales {
	return &scales{
		namespace: namespace,
		scale:     scale,
	}
}

// #############################################################################

// ScalesGetter has a method to return a ScaleInterface.
// A group's client should implement this interface.
type ScalesGetter interface {
	Scales(namespace string) ScaleInterface
}

// ScaleInterface has methods to work with Scale resources.
type ScaleInterface interface {
	ScaleExpansion
}

// k8sscales implements ScaleInterface
type k8sscales struct {
	client rest.Interface
	ns     string
}

// newScales returns a Scales
func newScales(c *v1beta1.ExtensionsV1beta1Client, namespace string) *k8sscales {
	return &k8sscales{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// The ScaleExpansion interface allows manually adding extra methods to the ScaleInterface.
type ScaleExpansion interface {
	Get(kind string, name string) (*k8sv1beta1.Scale, error)
	Update(kind string, scale *k8sv1beta1.Scale) (*k8sv1beta1.Scale, error)
}

// Get takes the reference to scale subresource and returns the subresource or error, if one occurs.
func (c *k8sscales) Get(kind string, name string) (result *k8sv1beta1.Scale, err error) {
	result = &k8sv1beta1.Scale{}

	// TODO this method needs to take a proper unambiguous kind
	fullyQualifiedKind := schema.GroupVersionKind{Kind: kind}
	resource, _ := meta.UnsafeGuessKindToResource(fullyQualifiedKind)

	err = c.client.Get().
		Namespace(c.ns).
		Resource(resource.Resource).
		Name(name).
		SubResource("scale").
		Do().
		Into(result)
	return
}

func (c *k8sscales) Update(kind string, scale *k8sv1beta1.Scale) (result *k8sv1beta1.Scale, err error) {
	result = &k8sv1beta1.Scale{}

	// TODO this method needs to take a proper unambiguous kind
	fullyQualifiedKind := schema.GroupVersionKind{Kind: kind}
	resource, _ := meta.UnsafeGuessKindToResource(fullyQualifiedKind)

	err = c.client.Put().
		Namespace(scale.Namespace).
		Resource(resource.Resource).
		Name(scale.Name).
		SubResource("scale").
		Body(scale).
		Do().
		Into(result)
	return
}
