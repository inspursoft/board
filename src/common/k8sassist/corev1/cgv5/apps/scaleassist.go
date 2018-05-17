// a temp file for building and guiding
package apps

import (
	"git/inspursoft/board/src/common/model"
)

type scales struct {
	namespace string
}

func (d *scales) Update(kind string, scale *model.Scale) (*model.Scale, error) {
	return nil, nil
}

func (d *scales) Get(kind string, name string) (*model.Scale, error) {
	return nil, nil
}

func NewScales(namespace string) *scales {
	return &scales{
		namespace: namespace,
	}
}
