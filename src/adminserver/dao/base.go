package dao

import (
	"database/sql"
	"os"
	"path"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

const AdminServerDbPath = "/data/adminserver/database/"
const AdminServerDbFile = "adminServer.db"

var (
	GlobalCache cache.Cache
)

func InitGlobalCache() (err error) {
	GlobalCache, err = cache.NewCache("memory", `{"interval": 3600}`)
	return err
}

func RegisterDatabase(dbFileName string) error {
	if err := orm.RegisterDriver("sqlite3", orm.DRSqlite); err != nil {
		return err
	}
	if err := orm.RegisterDataBase("default", "sqlite3", dbFileName); err != nil {
		return err
	}
	return nil
}

func InitDatabase() {
	if _, err := os.Stat(AdminServerDbPath); os.IsNotExist(err) {
		os.MkdirAll(AdminServerDbPath, os.ModePerm)
	}

	adminServerDbFileName := path.Join(AdminServerDbPath, AdminServerDbFile)
	if _, err := os.Stat(adminServerDbFileName); os.IsNotExist(err) {
		if errInitDb := InitDbTables(adminServerDbFileName); errInitDb != nil {
			logs.Error(errInitDb)
		}
	}
	if err := RegisterDatabase(adminServerDbFileName); err != nil {
		logs.Error(err)
	}
}

func InitDbTables(dbFileName string) error {
	db, err := sql.Open("sqlite3", dbFileName)
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
                        completed int,
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

	tokenTable := `create table if not exists token(
		id integer primary key autoincrement,
		time int not null,
		token varchar(30) not null
		);`
	if _, err := db.Exec(tokenTable); err != nil {
		return err
	}
	logs.Info("create token table successfully")

	return nil
}
