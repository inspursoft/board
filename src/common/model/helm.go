package model

import (
	"time"
)

type Repository struct {
	ID   int64  `json:"id" orm:"column(id)"`
	Name string `json:"name" orm:"column(name)"`
	URL  string `json:"url" orm:"column(url)"`
	Type int64  `json:"type" orm:"column(type)"`
}

type RepositoryDetail struct {
	Repository             `yaml:",inline"`
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
	ProjectId      int64     `json:"project_id"`
	ProjectName    string    `json:"project_name"`
	RepositoryId   int64     `json:"repoid"`
	RepositoryName string    `json:"repository"`
	Chart          string    `json:"chart"`
	ChartVersion   string    `json:"chartversion"`
	OwnerID        int64     `json:"owner_id,omitempty"`
	OwnerName      string    `json:"owner_name,omitempty"`
	Status         string    `json:"status,omitempty"`
	Values         string    `json:"values,omitempty"`
	UpdateTime     time.Time `json:"update_time,omitempty"`
	CreateTime     time.Time `json:"creation,omitempty"`
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
	ProjectId      int64     `orm:"column(project_id)"`
	ProjectName    string    `orm:"column(project_name)"`
	RepositoryId   int64     `orm:"column(repoid)"`
	RepostiroyName string    `orm:"column(repository)"`
	Workloads      string    `orm:"column(workloads)`
	OwnerID        int64     `orm:"column(owner_id)"`
	OwnerName      string    `orm:"column(owner_name)"`
	UpdateTime     time.Time `orm:"column(update_time)"`
	CreateTime     time.Time `orm:"column(creation_time)"`
}

func (rm *ReleaseModel) TableName() string {
	return "release"
}

type VisitorFunc func([]*Info, error) error

// GroupVersionKind unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

// Info contains temporary info to execute a REST call, or show the results
// of an already completed REST call.
type Info struct {
	// Namespace will be set if the object is namespaced and has a specified value.
	Namespace string
	Name      string
	GroupVersionKind

	// object source string
	Source string
}
