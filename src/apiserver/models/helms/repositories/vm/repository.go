package vm

import ()

type PaginatedHelmRepositoryDetail struct {
	HelmRepository         `json:",inline"`
	PaginatedChartVersions `json:",inline"`
}

type HelmRepository struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Type int64  `json:"type"`
}

type PaginatedChartVersions struct {
	Pagination        *Pagination      `json:"pagination"`
	ChartVersionsList []*ChartVersions `json:"charts"`
}

type Pagination struct {
	PageIndex  int   `json:"page_index"`
	PageSize   int   `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	PageCount  int   `json:"page_count"`
}

type ChartVersions struct {
	Name     string          `json:"name"`
	Versions []*ChartVersion `json:"versions"`
}

type ChartVersion struct {
	ChartMetadata `json:",inline"`
	URLs          []string `json:"urls"`
	Digest        string   `json:"digest,omitempty"`
}

type ChartMetadata struct {
	Name        string   `json:"name,omitempty"`
	Sources     []string `json:"sources,omitempty"`
	Version     string   `json:"version,omitempty"`
	KubeVersion string   `json:"kubeversion,omitempty"`
	Description string   `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
	Icon        string   `json:"icon,omitempty"`
}
