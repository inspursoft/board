package service

import (
	"fmt"

	"testing"
)

func TestK8sCliFactory(t *testing.T) {
	defer func() { recover() }()
	s, d := Suspend("10.110.18.71")
	fmt.Println(s, d)
}
