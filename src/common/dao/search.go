package dao

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type SearchResult struct {
	OwnerName   string `json:"owner_name" orm:"column(owner_name)"`
	ProjectName string `json:"project_name" orm:"column(project_name)"`
}

func SearchPrivite(projectName string, usrName string) ([]SearchResult, error) {
	var serachRes []SearchResult
	sql := `
SELECT
  project.owner_name AS owner_name,
  project.name       AS project_name
FROM user
  JOIN project_member ON user_id = project_member.user_id
  JOIN project ON project_id = project.id
  JOIN role ON project_member.role_id = role.id
WHERE project.deleted = 0 AND (
  (user.username = ? AND project.name LIKE ?)
  OR (project.public = 1 AND project.name LIKE ?)
  OR ((SELECT max(role.id)
       FROM user
         JOIN project_member ON user_id = project_member.user_id
         JOIN project ON project_id = project.id
         JOIN role ON project_member.role_id = role.id
       WHERE user.username = ?) = 1 AND project.name LIKE ?));	`
	o := orm.NewOrm()
	_, err := o.Raw(sql, usrName, "%"+projectName+"%", "%"+projectName+"%", usrName, "%"+projectName+"%").QueryRows(&serachRes)
	return serachRes, err
}
func SearchPublic(projectName string) ([]SearchResult, error) {
	var serachRes []SearchResult
	sql := fmt.Sprintf(`
	SELECT
	  owner_name AS owner_name,
	  project.name AS project_name
	FROM project
	WHERE public = 1
	AND project.name LIKE ?;
		`)
	o := orm.NewOrm()
	_, err := o.Raw(sql, "%"+projectName+"%").QueryRows(&serachRes)
	return serachRes, err
}
