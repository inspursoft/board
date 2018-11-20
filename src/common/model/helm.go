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
