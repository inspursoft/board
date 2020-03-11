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
	return o.Delete(&projectMember)
}

func GetProjectMembers(project model.Project) ([]*model.ProjectMember, error) {
	o := orm.NewOrm()
	sql := `select pm.id, pm.user_id, u.username, pm.project_id, pm.role_id
		from user u left join project_member pm 
				on u.id = pm.user_id 
		where pm.project_id = ?`
	var members []*model.ProjectMember
	_, err := o.Raw(sql, project.ID).QueryRows(&members)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return members, nil
}

func GetProjectMemberRole(project model.Project, user model.User) (*model.Role, error) {
	o := orm.NewOrm()
	sql := `select r.id, r.name, r.comment 
		from project_member pm left join role r
				on pm.role_id = r.id
		where pm.project_id = ?
				and pm.user_id = ?`
	var role model.Role
	err := o.Raw(sql, project.ID, user.ID).QueryRow(&role)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil

}
