package app

import (
	"testing"
	"fmt"
)

func TestGetRunFlag(t *testing.T) {
	runFlag := getRunFlag()
	for k,v:=range runFlag{
		fmt.Println(k,v,*v)
	}

}

