package dao_test

import (
	"fmt"
	"git/inspursoft/board/src/common/dao"

	"testing"
)

func TestSearchPrivite(t *testing.T) {
	fmt.Println(dao.SearchPrivateProject("l", "Admin"))
	fmt.Println(dao.SearchPublicProject("l"))
}
