package dao

import (
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/orm"
)

func InsertOrUpdateProjectMember(projectMember model.ProjectMember) (int64, error) {
	o := orm.NewOrm()
	ptmt, err := o.Raw(`insert into project_member
		 (id, project_id, user_id, role_id)
	 		values (?, ?, ?, ?) 
			on duplicate key 
			update role_id = ?`).Prepare()
	if err != nil {
		return 0, err
	}
	defer ptmt.Close()
	pmGeneratedID := projectMember.UserID + projectMember.ProjectID
	r, err := ptmt.
		Exec(pmGeneratedID,
			projectMember.ProjectID,
			projectMember.UserID,
			projectMember.RoleID,
			projectMember.RoleID)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return r.RowsAffected()
}

func DeleteProjectMember(projectMember model.ProjectMember) (int64, error) {
	o := orm.NewOrm()
	return o.Delete(&projectMember, "id")
}

func GetProjectMembers(project model.Project, user model.User) ([]*model.User, error) {
	o := orm.NewOrm()
	sql := `select u.id, u.username 
		from user u left join project_member pm 
				on u.id = pm.user_id 
	  where pm.project_id = ? 
				and pm.user_id = ?`
	var users []*model.User
	_, err := o.Raw(sql, project.ID, user.ID).QueryRows(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
