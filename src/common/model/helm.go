package model

import (
	"bytes"
	"fmt"
	"io"
	"time"

	gyaml "github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type HelmRepository struct {
	ID   int64  `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
	URL  string `json:"url" orm:"column(url)"`
	Type int64  `json:"type" orm:"column(type)"`
}

type HelmRepositoryDetail struct {
	HelmRepository    `yaml:",inline"`
	ChartVersionsList []*ChartVersions `json:"charts"`
}

type PaginatedHelmRepositoryDetail struct {
	HelmRepository         `yaml:",inline"`
	PaginatedChartVersions `yaml:",inline"`
}

type PaginatedChartVersions struct {
	Pagination        *Pagination      `json:"pagination"`
	ChartVersionsList []*ChartVersions `json:"charts"`
}

type ChartVersions struct {
	Name     string          `json:"name" yaml:"name"`
	Versions []*ChartVersion `json:"versions" yaml:"versions"`
}

type ChartVersion struct {
	ChartMetadata `yaml:",inline"`
	URLs          []string `json:"urls" yaml:"urls"`
	Digest        string   `json:"digest,omitempty" yaml:"digest,omitempty"`
}

type ChartMetadata struct {
	Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
	Sources     []string `json:"sources,omitempty" yaml:"sources,omitempty"`
	Version     string   `json:"version,omitempty" yaml:"version,omitempty"`
	KubeVersion string   `json:"kubeVersion,omitempty" yaml:"kubeVersion,omitempty"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty" yaml:"keywords,omitempty"`
	Icon        string   `json:"icon,omitempty" yaml:"icon,omitempty"`
}

// 	Chart is a helm package that contains metadata, a default config, zero or more
// 	optionally parameterizable templates, and zero or more charts (dependencies).
type Chart struct {
	// Contents of the Chartfile.
	Metadata *ChartMetadata `json:"metadata,omitempty"`
	// Templates for this chart.
	Templates []*File `json:"templates,omitempty"`
	// Charts that this chart depends on.
	//Dependencies []*Chart `json:"dependencies,omitempty"`
	// Default config for this template.
	Values string `json:"values,omitempty"`
	// Miscellaneous files in a chart archive,
	// e.g. README, LICENSE, etc.
	Files []*File `json:"files,omitempty"`
}

type File struct {
	Name     string `json:"name,omitempty"`
	Contents string `json:"contents,omitempty"`
}

type Release struct {
	ID             int64     `json:"id,omitempty"`
	Name           string    `json:"name"`
	ProjectID      int64     `json:"project_id"`
	ProjectName    string    `json:"project_name"`
	RepositoryID   int64     `json:"repository_id"`
	RepositoryName string    `json:"repository"`
	Chart          string    `json:"chart"`
	ChartVersion   string    `json:"chartversion"`
	OwnerID        int64     `json:"owner_id,omitempty"`
	OwnerName      string    `json:"owner_name,omitempty"`
	Status         string    `json:"status,omitempty"`
	Values         string    `json:"values,omitempty"`
	UpdateTime     time.Time `json:"update_time,omitempty"`
	CreateTime     time.Time `json:"creation_time,omitempty"`
}

type ReleaseDetail struct {
	Release        `yaml:",inline"`
	Notes          string `json:"notes,omitempty" yaml:"notes,omitempty"`
	Workloads      string `json:"workloads,omitempty" yaml:"workloads,omitempty"`
	WorkloadStatus string `json:"workloadstatus,omitempty" yaml:"workloadstatus,omitempty"`
}

type ReleaseModel struct {
	ID             int64     `orm:"column(id)"`
	Name           string    `orm:"column(name)"`
	ProjectID      int64     `orm:"column(project_id)"`
	ProjectName    string    `orm:"column(project_name)"`
	RepositoryID   int64     `orm:"column(repository_id)"`
	RepostiroyName string    `orm:"column(repository)"`
	Workloads      string    `orm:"column(workloads)`
	OwnerID        int64     `orm:"column(owner_id)"`
	OwnerName      string    `orm:"column(owner_name)"`
	UpdateTime     time.Time `orm:"column(update_time)"`
	CreateTime     time.Time `orm:"column(creation_time)"`
}

func (rm *ReleaseModel) TableName() string {
	return "helm_release"
}

type VisitorFunc func([]*K8sInfo, error) error

// Modifier modify the object
type Modifier func(interface{}) (interface{}, error)

// GroupVersionKind unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

// Info contains temporary info to execute a REST call, or show the results
// of an already completed REST call.
type K8sInfo struct {
	// Namespace will be set if the object is namespaced and has a specified value.
	Namespace string
	Name      string
	GroupVersionKind

	// object source string
	Source string
}

// K8sHelper is a convenience struct for holding references to the interfaces
// needed to create K8sInfo for arbitrary objects.
type K8sHelper struct {
	meta.MetadataAccessor
	runtime.Decoder
}

// infoForData creates an K8sInfo object for the given data. An error is returned
// if any of the decoding or client lookup steps fail. Name and namespace will be
// set into K8sInfo if the mapping's MetadataAccessor can retrieve them.
func (m *K8sHelper) infoForData(data []byte) (*K8sInfo, error) {
	obj, gvk, err := m.Decode(data, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to decode %s: %v", string(data), err)
	}

	name, _ := m.MetadataAccessor.Name(obj)
	namespace, _ := m.MetadataAccessor.Namespace(obj)

	return &K8sInfo{
		Source:    string(data),
		Namespace: namespace,
		Name:      name,
		GroupVersionKind: GroupVersionKind{
			Group:   gvk.Group,
			Version: gvk.Version,
			Kind:    gvk.Kind,
		},
	}, nil
}

func (m *K8sHelper) Visit(info string, fn VisitorFunc) error {
	buffer := bytes.NewBufferString(info)
	d := yaml.NewYAMLOrJSONDecoder(buffer, 4096)
	var infos []*K8sInfo
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

func (m *K8sHelper) Transform(in string, fn Modifier) (string, error) {
	src := make(map[string]interface{})
	err := gyaml.Unmarshal([]byte(in), &src)
	if err != nil {
		return in, err
	}
	target, err := fn(src)
	if err != nil {
		return in, err
	}
	bs, err := gyaml.Marshal(target)
	if err != nil {
		return in, err
	}
	return string(bs), nil
}

func NewK8sHelper() *K8sHelper {
	return &K8sHelper{
		meta.NewAccessor(),
		unstructured.UnstructuredJSONScheme,
	}
}
