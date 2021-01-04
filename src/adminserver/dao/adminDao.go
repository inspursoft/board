package dao

import (
	"fmt"
	"github.com/inspursoft/board/src/adminserver/common"
	"github.com/inspursoft/board/src/adminserver/models"
	"github.com/inspursoft/board/src/common/model"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func CheckDB() error {
	o := orm.NewOrm()
	err := o.Using("mysql-db2")
	if err != nil {
		return err
	}
	_, err = o.Raw("SELECT 1").Exec()
	return err
}

func InitAdmin(user model.User) error {
	o := orm.NewOrm()
	o.Using("mysql-db2")
	user.UpdateTime = time.Now()
	_, err := o.Update(&user, "password", "salt")
	if err != nil {
		return err
	}
	logs.Info("Admin password has been updated successfully.")
	return nil
}

func LoginCheckAuth(user model.User) (model.User, error) {
	o := orm.NewOrm()
	o.Using("mysql-db2")
	err := o.Read(&user, "username", "system_admin", "deleted")
	if err == orm.ErrNoRows {
		return user, common.ErrForbidden
	}
	return user, nil
}

func LoginCheckPassword(user model.User) (model.User, error) {
	o := orm.NewOrm()
	o.Using("mysql-db2")
	err := o.Read(&user, "username", "password")
	if err == orm.ErrNoRows {
		return user, common.ErrWrongPassword
	}
	return user, nil
}

func GetUserByID(user model.User) (model.User, error) {
	o := orm.NewOrm()
	o.Using("mysql-db2")
	return user, o.Read(&user, "id", "deleted")
}

func UpdateUUIDToken(newtoken models.Token) error {
	existingToken := models.Token{Id: 1}
	o := orm.NewOrm()
	if o.Read(&existingToken) == orm.ErrNoRows {
		if _, err := o.Insert(&newtoken); err != nil {
			return err
		}
	} else if (newtoken.Time - existingToken.Time) > 1800 {
		if _, err := o.Update(&newtoken); err != nil {
			return err
		}
	} else {
		return common.ErrAdminLogin
	}
	return nil
}

func GetUUIDToken(token models.Token) (models.Token, error) {
	o := orm.NewOrm()
	return token, o.Read(&token)
}

func RemoveUUIDToken() error {
	o := orm.NewOrm()
	if _, err := o.Delete(&models.Token{Id: 1}); err != nil {
		return err
	}
	return nil
}

func RegisterDB() error {
	DBpassword, err := common.ReadCfgItem("db_password", "/go/cfgfile/board.cfg")
	if err != nil {
		return err
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("mysql-db2", "mysql", fmt.Sprintf("root:%s@tcp(%s:%d)/board?charset=utf8", DBpassword, "db", 3306))
	if err != nil {
		return err
	}
	logs.Info("register DB success")
	return nil
}
