// a temp file for building and guiding
package apps

import (
	//api "k8s.io/client-go/pkg/api"
	//v1 "k8s.io/client-go/pkg/api/v1"
	//watch "k8s.io/client-go/pkg/watch"
	//rest "k8s.io/client-go/rest"
	"git/inspursoft/board/src/common/model"
)

type scales struct {
	ns string
}

func (d *scales) Update(kind string, scale *model.Scale) (*model.Scale, error) {
	return nil, nil
}

func (d *scales) Get(kind string, name string) (*model.Scale, error) {
	return nil, nil
}
