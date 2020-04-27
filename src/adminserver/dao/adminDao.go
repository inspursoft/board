package dao

import (
	"errors"
	"git/inspursoft/board/src/adminserver/models"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

var ErrAdminLogin = errors.New("another admin user has signed in other place")
var ErrForbidden = errors.New("Forbidden")
var ErrWrongPassword = errors.New("Wrong password")

func UsingDB(which string) orm.Ormer {
	o := orm.NewOrm()
	o.Using(which)
	return o
}

func CheckDB() error {
	o := UsingDB("mysql-db2")
	_, err := o.Raw("SELECT 1").Exec()
	return err
}

func InitAdmin(user models.User) error {
	o := UsingDB("mysql-db2")
	user.UpdateTime = time.Now()
	_, err := o.Update(&user, "password", "salt")
	if err != nil {
		return err
	}
	logs.Info("Admin password has been updated successfully.")
	return nil
}

func CacheAccountInfo(account models.Account) error {
	o := UsingDB("default")
	if o.Read(&models.Account{Id: 1}) == orm.ErrNoRows {
		if _, err := o.Insert(&account); err != nil {
			return err
		}
	} else {
		if _, err := o.Update(&account); err != nil {
			return err
		}
	}
	return nil
}

func LoginCheckAuth(user models.User) error {
	o := UsingDB("mysql-db2")
	err := o.Read(&user, "username", "system_admin", "deleted")
	if err == orm.ErrNoRows {
		return ErrForbidden
	}
	return nil
}

func LoginCheckPassword(user models.User) error {
	o := UsingDB("mysql-db2")
	err := o.Read(&user, "username", "password")
	if err == orm.ErrNoRows {
		return ErrWrongPassword
	}
	return nil
}

func GetUserByID(user models.User) error {
	o := UsingDB("mysql-db2")
	return o.Read(&user, "id", "deleted")
}

func UpdateUUIDToken(newtoken models.Token) error {
	existingToken := models.Token{Id: 1}
	o := UsingDB("default")
	if o.Read(&existingToken) == orm.ErrNoRows {
		if _, err := o.Insert(&newtoken); err != nil {
			return err
		}
	} else if (newtoken.Time - existingToken.Time) > 1800 {
		if _, err := o.Update(&newtoken); err != nil {
			return err
		}
	} else {
		return ErrAdminLogin
	}
	return nil
}

func GetUUIDToken(token models.Token) error {
	o := UsingDB("default")
	return o.Read(&token)
}

func RemoveUUIDToken() error {
	o := UsingDB("default")
	if _, err := o.Delete(&models.Token{Id: 1}); err != nil {
		return err
	}
	return nil
}
