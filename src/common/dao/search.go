package dao

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type SearchProjectResult struct {
	OwnerName   string `json:"owner_name" orm:"column(owner_name)"`
	ProjectName string `json:"project_name" orm:"column(project_name)"`
	IsPublic    bool `json:"is_public"  orm:"column(is_public)"`
}
type SearchUserResult struct {
	UserName string `json:"user_name" orm:"column(user_name)"`
	RoleName string`json:"user_name" orm:"column(role_name)"`
	UserEmail string `json:"user_email" orm:"column(user_email)"`
}


func SearchPrivateProject(projectName string, usrName string) ([]SearchProjectResult, error) {
	var searchRes []SearchProjectResult
	sql := `
SELECT DISTINCT
  project.owner_name AS owner_name,
  project.name       AS project_name,
  project.public     AS is_public
FROM user
  JOIN project_member ON user_id = project_member.user_id
  JOIN project ON project_id = project.id
  JOIN role ON project_member.role_id = role.id
WHERE project.deleted = 0
      AND project.name LIKE ?
      AND ( user.username = ?
           OR (project.public = 1)
           OR ((SELECT DISTINCT user.system_admin
                FROM user
                WHERE user.username = ?) = 1));`
	o := orm.NewOrm()
	_, err := o.Raw(sql, "%"+projectName+"%", usrName, usrName).QueryRows(&searchRes)
	return searchRes, err
}
func SearchPublicProject(projectName string) ([]SearchProjectResult, error) {
	var searchRes []SearchProjectResult
	sql := fmt.Sprintf(`
	SELECT
	  owner_name AS owner_name,
	  project.name AS project_name,
	  project.public     AS is_public
	FROM project
	WHERE public = 1
	AND project.name LIKE ?;
		`)
	o := orm.NewOrm()
	_, err := o.Raw(sql, "%"+projectName+"%").QueryRows(&searchRes)
	return searchRes, err
}
func SearchUser(activeUser string, searchName string) ([]SearchUserResult, error) {
	var searchRes []SearchUserResult
	sql := `SELECT
  user.username AS user_name,
  role.name     AS role_name,
  user.email    AS user_email
FROM project_member
  JOIN role ON project_member.role_id = role.id
  JOIN user ON project_member.user_id = user.id
WHERE user.deleted = 0
      AND project_id = (SELECT project_id
                        FROM project_member
                          JOIN user ON project_member.user_id = user.id
                          JOIN project ON project_member.project_id = project.id
                        WHERE project.deleted = 0
                              AND user.username = ?)
      AND user.username LIKE ?;`
	o := orm.NewOrm()
	_, err := o.Raw(sql, activeUser, "%"+searchName+"%").QueryRows(&searchRes)
	return searchRes, err
}
