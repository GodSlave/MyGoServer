package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/GodSlave/MyGoServer/log"
)

type BaseSql struct {
	Url     string
	Engine  *xorm.Engine
	version int32
}

func (sql *BaseSql) InitDB() {
	var engine *xorm.Engine
	var err error
	engine, err = xorm.NewEngine("mysql", sql.Url)
	if err == nil {
		sql.Engine = engine
	} else {
		log.Error("db init error :  " + err.Error())
		panic(err)
	}
}

func (sql *BaseSql) CheckMigrate() {
	exit, err := sql.Engine.IsTableExist("Version")
	checkError(err)
	if exit {
		var version DBVersion
		err := sql.Engine.Limit(0, 1).Find(&version)
		version.Version = 2
		checkError(err)
		sql.migrate(version.Version)
	} else {
		sql.migrate(0)
	}

}

func (sql *BaseSql) migrate(oldversion int32) {
	if oldversion == 0 {
		err := sql.Engine.CreateTables(&DBVersion{})
		sql.Engine.Insert(&DBVersion{Version: 1})
		checkError(err)
	}

	if oldversion == 1 {
		sql.Engine.SQL("update d_b_version set version = 2 ;")
	}

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
