package db

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"king/config"
	"king/model"
	"log"
)

type dbm struct {
	Db *sql.DB
}

func (s *dbm) Orm() orm.Ormer {
	return Orm()
}

var Db = &dbm{nil}

func Orm() orm.Ormer {
	return orm.NewOrm()
}

func Connect() {
	if err := orm.RegisterDataBase("default", "mysql", config.GetString("dataSourceName"), 30); err != nil {
		glog.Errorln(err)
		return
	}
	orm.RegisterModel(
		new(model.Group),
		new(model.WebServer),
		new(model.Config),
		new(model.UpFile),
		new(model.Version),
	)
	log.Println("database connected")
	Db.Db, _ = orm.GetDB()
}

func Close() {
	Db.Db.Close()
}

func IsConnected() bool {
	return Db.Db != nil
}
func Truncate(table string) error {
	res, err := Orm().Raw("TRUNCATE `king`.`" + table + "`").Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return err
}
