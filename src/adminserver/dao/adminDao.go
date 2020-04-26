package dao

import (
	"errors"
	"fmt"
	"git/inspursoft/board/src/adminserver/models"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func RegisterDB() error {
	b, err := ioutil.ReadFile(path.Join(models.DBconfigdir, "/env"))
	if err != nil {
		return err
	}
	DBpassword := strings.TrimPrefix(string(b), "DB_PASSWORD=")
	DBpassword = strings.Replace(DBpassword, "\n", "", 1)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("mysql-db2", "mysql", fmt.Sprintf("root:%s@tcp(%s:%d)/board?charset=utf8", DBpassword, "db", 3306))
	if err != nil {
		return err
	}
	logs.Info("register DB success")
	return nil
}

func CheckDB() error {
	o := orm.NewOrm()
	err := o.Using("mysql-db2")
	if err != nil {
		return err
	}
	_, err = o.Raw("SELECT 1").Exec()
	return err
}

func InitAdmin(user models.User) error {
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

func CacheAccountInfo(account models.Account) error {
	o := orm.NewOrm()
	o.Using("default")
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
	o := orm.NewOrm()
	o.Using("mysql-db2")
	err := o.Read(&user, "username", "system_admin", "deleted")
	if err == orm.ErrNoRows {
		return errors.New("Forbidden")
	}
	return nil
}

func LoginCheckPassword(user models.User) error {
	o := orm.NewOrm()
	o.Using("mysql-db2")
	err := o.Read(&user, "username", "password")
	if err == orm.ErrNoRows {
		return errors.New("Wrong password")
	}
	return nil
}

func GetUserByID(user models.User) error {
	o := orm.NewOrm()
	o.Using("mysql-db2")
	return o.Read(&user, "id", "deleted")
}

func UpdateUUIDToken(newtoken models.Token) error {
	existingToken := models.Token{Id: 1}
	o := orm.NewOrm()
	o.Using("default")
	if o.Read(&existingToken) == orm.ErrNoRows {
		if _, err := o.Insert(&newtoken); err != nil {
			return err
		}
	} else if (newtoken.Time - existingToken.Time) > 1800 {
		if _, err := o.Update(&newtoken); err != nil {
			return err
		}
	} else {
		return errors.New("another admin user has signed in other place")
	}
	return nil
}

func GetUUIDToken(token models.Token) error {
	o := orm.NewOrm()
	o.Using("default")
	return o.Read(&token)
}

func RemoveUUIDToken() error {
	o := orm.NewOrm()
	o.Using("default")
	if _, err := o.Delete(&models.Token{Id: 1}); err != nil {
		return err
	}
	return nil
}
