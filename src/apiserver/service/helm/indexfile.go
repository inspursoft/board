package helm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v2"

	"github.com/inspursoft/board/src/common/model"
)

var (
	// ErrNoChartVersion indicates that a chart with the given version is not found.
	ErrNoChartVersion = errors.New("no chart version found")
	// ErrNoChartName indicates that a chart with the given name is not found.
	ErrNoChartName = errors.New("no chart name found")
)

type IndexFile struct {
	Entries map[string][]*model.ChartVersion `json:"entries" yaml:"entries"`
}

// Has returns true if the index has an entry for a chart with the given name and exact version.
func (i IndexFile) Has(name, version string) bool {
	_, err := i.Get(name, version)
	return err == nil
}

// SortEntries sorts the entries by version in descending order.
//
// In canonical form, the individual version records should be sorted so that
// the most recent release for every version is in the 0th slot in the
// Entries.ChartVersions array. That way, tooling can predict the newest
// version without needing to parse SemVers.
func (i IndexFile) SortEntries() {
	for _, versions := range i.Entries {
		sort.Sort(sort.Reverse(ChartVersions(versions)))
	}
}

// Get returns the ChartVersion for the given name.
//
// If version is empty, this will return the chart with the highest version.
func (i IndexFile) Get(name, version string) (*model.ChartVersion, error) {
	vs, ok := i.Entries[name]
	if !ok {
		return nil, ErrNoChartName
	}
	if len(vs) == 0 {
		return nil, ErrNoChartVersion
	}

	var constraint *semver.Constraints
	if len(version) == 0 {
		constraint, _ = semver.NewConstraint("*")
	} else {
		var err error
		constraint, err = semver.NewConstraint(version)
		if err != nil {
			return nil, err
		}
	}

	for _, ver := range vs {
		test, err := semver.NewVersion(ver.Version)
		if err != nil {
			continue
		}

		if constraint.Check(test) {
			return ver, nil
		}
	}
	return nil, fmt.Errorf("No chart version found for %s-%s", name, version)
}

// WriteFile writes an index file to the given destination path.
//
// The mode on the file is set to 'mode'.
func (i IndexFile) WriteFile(dest string, mode os.FileMode) error {
	b, err := yaml.Marshal(i)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dest, b, mode)
}

// Merge merges the given index file into this index.
//
// This merges by name and version.
//
// If one of the entries in the given index does _not_ already exist, it is added.
// In all other cases, the existing record is preserved.
//
// This can leave the index in an unsorted state
func (i *IndexFile) Merge(f *IndexFile) {
	for _, cvs := range f.Entries {
		for _, cv := range cvs {
			if !i.Has(cv.Name, cv.Version) {
				e := i.Entries[cv.Name]
				i.Entries[cv.Name] = append(e, cv)
			}
		}
	}
}

// NewIndexFile initializes an index.
func NewIndexFile() *IndexFile {
	return &IndexFile{
		Entries: map[string][]*model.ChartVersion{},
	}
}
