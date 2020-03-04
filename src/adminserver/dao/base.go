package dao

import (
	"database/sql"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

const AdminServerDbFile = "/root/adminServer.db"

var (
	GlobalCache cache.Cache
)

func InitGlobalCache() (err error) {
	GlobalCache, err = cache.NewCache("memory", `{"interval": 3600}`);
	return err;
}

func RegisterDatabase() error {
	if err := orm.RegisterDriver("sqlite3", orm.DRSqlite); err != nil {
		return err
	}
	if err := orm.RegisterDataBase("default", "sqlite3", AdminServerDbFile); err != nil {
		return err
	}
	return nil
}

func InitDatabase() error {
	db, err := sql.Open("sqlite3", AdminServerDbFile)
	if err != nil {
		return err
	}
	defer db.Close()
	nodeStateTable := `create table if not exists node_status(
                        id integer primary key autoincrement,
                        ip varchar(30) not null,
                        creation_time int not null
	                   );`
	if _, err := db.Exec(nodeStateTable); err != nil {
		return err
	}
	logs.Info("create node_status table successfully")

	nodeLogTable := `create table if not exists node_log(
                        id integer primary key autoincrement,
                        ip varchar(30) not null,
                        creation_time int not null,
                        log_type int,
                        success int,
                        pid int
	                   );`
	if _, err := db.Exec(nodeLogTable); err != nil {
		return err
	}
	logs.Info("create node_logs table successfully")

	nodeDetailTable := `create table if not exists node_log_detail_info(
                        id integer primary key autoincrement,
                        creation_time int not null,
                        detail text not null
	                   );`
	if _, err := db.Exec(nodeDetailTable); err != nil {
		return err
	}
	logs.Info("create node_log_detail table successfully")

	return nil
}
