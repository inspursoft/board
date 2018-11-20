package helm

import (
	"bytes"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type VisitorFunc func([]*Info) error

// Info contains temporary info to execute a REST call, or show the results
// of an already completed REST call.
type Info struct {
	Accessor meta.MetadataAccessor
	// Namespace will be set if the object is namespaced and has a specified value.
	Namespace        string
	Name             string
	GroupVersionKind *schema.GroupVersionKind

	// Optional, Source is the filename or URL to template file (.json or .yaml),
	// or stdin to use to handle the resource
	Source string
	// Optional, this is the most recent value returned by the server if available. It will
	// typically be in unstructured or internal forms, depending on how the Builder was
	// defined. If retrieved from the server, the Builder expects the mapping client to
	// decide the final form. Use the AsVersioned, AsUnstructured, and AsInternal helpers
	// to alter the object versions.
	Object runtime.Object
}

// Mapper is a convenience struct for holding references to the interfaces
// needed to create Info for arbitrary objects.
type Mapper struct {
	meta.MetadataAccessor
	runtime.Decoder
}

// InfoForData creates an Info object for the given data. An error is returned
// if any of the decoding or client lookup steps fail. Name and namespace will be
// set into Info if the mapping's MetadataAccessor can retrieve them.
func (m *Mapper) InfoForData(data []byte) (*Info, error) {
	obj, gvk, err := m.Decode(data, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to decode %s: %v", string(data), err)
	}

	name, _ := m.MetadataAccessor.Name(obj)
	namespace, _ := m.MetadataAccessor.Namespace(obj)

	return &Info{
		Accessor:         m.MetadataAccessor,
		Source:           string(data),
		Namespace:        namespace,
		Name:             name,
		GroupVersionKind: gvk,
		Object:           obj,
	}, nil
}

func (m *Mapper) Visit(info string, fn VisitorFunc) error {
	buffer := bytes.NewBufferString(info)
	d := yaml.NewYAMLOrJSONDecoder(buffer, 4096)
	var infos []*Info
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

		info, err := m.InfoForData(ext.Raw)
		if err != nil {
			return err
		}
		infos = append(infos, info)
	}

	return fn(infos)
}

func NewMapper() *Mapper {
	return &Mapper{
		meta.NewAccessor(),
		unstructured.UnstructuredJSONScheme,
	}
}
