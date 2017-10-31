package collect_test

import (
	"git/inspursoft/board/src/collector/service/collect"
	"testing"
)

func TestGetNodeMachine(t *testing.T) {
	collect.SetInitVar("10.110.18.26", "8080")
	collect.GetNodeMachine("10.110.18.71")
}
