package helm

import (
	"github.com/inspursoft/board/src/common/model"
	"github.com/Masterminds/semver"
)

type ChartVersions []*model.ChartVersion

// Len returns the length.
func (c ChartVersions) Len() int { return len(c) }

// Swap swaps the position of two items in the versions slice.
func (c ChartVersions) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less returns true if the version of entry a is less than the version of entry b.
func (c ChartVersions) Less(a, b int) bool {
	// Failed parse pushes to the back.
	i, err := semver.NewVersion(c[a].Version)
	if err != nil {
		return true
	}
	j, err := semver.NewVersion(c[b].Version)
	if err != nil {
		return false
	}
	return i.LessThan(j)
}

type SortChartVersionsByName []*model.ChartVersions

// Len returns the length.
func (c SortChartVersionsByName) Len() int { return len(c) }

// Swap swaps the position of two items in the versions slice.
func (c SortChartVersionsByName) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less returns true if the version of entry a is less than the version of entry b.
func (c SortChartVersionsByName) Less(a, b int) bool {
	return c[a].Name < c[b].Name
}
