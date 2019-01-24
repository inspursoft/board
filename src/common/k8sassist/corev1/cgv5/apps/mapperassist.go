package apps

import (
	"bytes"
	"fmt"
	"io"

	"git/inspursoft/board/src/common/model"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Mapper is a convenience struct for holding references to the interfaces
// needed to create Info for arbitrary objects.
type mapper struct {
	meta.MetadataAccessor
	runtime.Decoder
}

// infoForData creates an Info object for the given data. An error is returned
// if any of the decoding or client lookup steps fail. Name and namespace will be
// set into Info if the mapping's MetadataAccessor can retrieve them.
func (m *mapper) infoForData(data []byte) (*model.Info, error) {
	obj, gvk, err := m.Decode(data, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to decode %s: %v", string(data), err)
	}

	name, _ := m.MetadataAccessor.Name(obj)
	namespace, _ := m.MetadataAccessor.Namespace(obj)

	return &model.Info{
		Source:    string(data),
		Namespace: namespace,
		Name:      name,
		GroupVersionKind: model.GroupVersionKind{
			Group:   gvk.Group,
			Version: gvk.Version,
			Kind:    gvk.Kind,
		},
	}, nil
}

func (m *mapper) Visit(info string, fn model.VisitorFunc) error {
	buffer := bytes.NewBufferString(info)
	d := yaml.NewYAMLOrJSONDecoder(buffer, 4096)
	var infos []*model.Info
	errs := []error(nil)
	for {
		ext := runtime.RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// TODO: This needs to be able to handle object in other encodings and schemas.
		ext.Raw = bytes.TrimSpace(ext.Raw)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}

		info, err := m.infoForData(ext.Raw)
		if err != nil {
			errs = append(errs, err)
		}
		infos = append(infos, info)
	}

	return fn(infos, utilerrors.NewAggregate(errs))
}

func NewMapper() *mapper {
	return &mapper{
		meta.NewAccessor(),
		unstructured.UnstructuredJSONScheme,
	}
}
