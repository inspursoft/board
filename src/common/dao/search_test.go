package dao

import (
	"fmt"

	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSearchPrivite(t *testing.T) {
	fmt.Println(SearchPrivateProject("l", "Admin"))
	fmt.Println(SearchPublicProject("l"))
}
