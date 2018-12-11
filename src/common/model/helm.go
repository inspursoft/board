package model

import ()

type Repository struct {
	ID       int64  `json:"id" orm:"column(id)"`
	Name     string `json:"name" orm:"column(name)"`
	URL      string `json:"url" orm:"column(url)"`
	Username string `json:"username" orm:"column(username)"`
	Password string `json:"password" orm:"column(password)"`
	Cert     string `json:"cert" orm:"column(cert)"`
	Key      string `json:"key" orm:"column(key)"`
	CA       string `json:"ca" orm:"column(ca)"`
	Type     int64  `json:"type" orm:"column(type)"`
}

type RepositoryDetail struct {
	Repository *Repository `yaml:",inline"`
	IndexFile  *IndexFile  `json:"index" yaml:"index"`
}

type IndexFile struct {
	Entries map[string]ChartVersions `json:"entries" yaml:"entries"`
}

type ChartVersions []*ChartVersion

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

type ReleaseService struct {
	ReleaseId   int64  `json:"releaseid,omitempty"`
	ServiceId   int64  `json:"serviceid"`
	ServiceName string `json:"servicename,omitempty"`
}

type Release struct {
	ID           int64            `json:"id,omitempty"`
	Name         string           `json:"name"`
	ProjectId    int64            `json:"project_id"`
	RepositoryId int64            `json:"repoid"`
	Chart        string           `json:"chart"`
	ChartVersion string           `json:"chartversion"`
	Value        string           `json:"value,omitempty"`
	Workloads    string           `json:"workload,omitempty"`
	Services     []ReleaseService `json:"services,omitempty"`
}

type ReleaseModel struct {
	ID           int64                 `orm:"column(id)"`
	Name         string                `orm:"column(name)"`
	ProjectId    int64                 `orm:"column(project_id)"`
	RepositoryId int64                 `orm:"column(repoid)"`
	Chart        string                `orm:"column(chart)"`
	ChartVersion string                `orm:"column(chartversion)"`
	Value        string                `orm:"column(value)"`
	Workloads    string                `orm:"column(workload)"`
	Services     []ReleaseServiceModel `orm:"-"`
}

func (rm *ReleaseModel) TableName() string {
	return "release"
}

type ReleaseServiceModel struct {
	ReleaseId int64
	ServiceId int64
}

type VisitorFunc func([]*Info) error

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

	// Optional, Source is the filename or URL to template file (.json or .yaml),
	// or stdin to use to handle the resource
	Source string
}
