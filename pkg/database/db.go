package database

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// singleton instance of database connection.
var (
	dbInstance        *gorm.DB
	dbOnce            sync.Once
	templatePostgres  = "host={host} port={port} user={username} dbname={name} password={password} {options}"
	templateMysql     = "{username}:{password}@tcp({host}:{port})/{name}?{options}"
	templateSqlServer = "sqlserver://{username}:{password}@{host}:{port}?database={name}&{options}"

	optionPlaceholders = map[string]string{
		"{username}": "username",
		"{password}": "password",
		"{host}":     "host",
		"{name}":     "name",
		"{options}":  "options",
	}
	maxOpenConn     = 25
	maxIdleConn     = 25
	maxLifeTimeConn = 300
	maxIdleTimeConn = 300
)

// Conn alias for DB().
func Conn() *gorm.DB {
	return DB()
}

// DB creates a new instance of gorm.DB if a connection is not established.
// return singleton instance.
func DB() *gorm.DB {
	if dbInstance == nil {
		dbOnce.Do(func() {
			openDB()
		})
	}
	return dbInstance
}

// openDB initialize gorm DB.
func openDB() {
	var dbDebug = true
	var logLevel = logger.Info

	if debug, err := beego.AppConfig.Bool("database::debug"); err == nil {
		dbDebug = debug
	}

	if !dbDebug {
		logLevel = logger.Silent
	}

	dbConfig, err := beego.AppConfig.GetSection("database")
	if err != nil {
		panic(err)
	}

	dbString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", dbConfig["username"], dbConfig["password"], dbConfig["host"], dbConfig["port"], dbConfig["name"], dbConfig["options"])

	// log.Println(dbString)
	gormDB, err := gorm.Open(mysql.Open(dbString), &gorm.Config{SkipDefaultTransaction: true,
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logLevel)})
	if err != nil {
		panic("cannot open database.")
	}
	dbInstance = gormDB
	sqlDb, err := dbInstance.DB()
	if err != nil {
		panic(err)
	}

	if parse, err := strconv.Atoi(dbConfig["maxopenconn"]); err == nil {
		maxOpenConn = parse
	}
	if parse, err := strconv.Atoi(dbConfig["maxidleconn"]); err == nil {
		maxIdleConn = parse
	}
	if parse, err := strconv.Atoi(dbConfig["maxlifetimeconn"]); err == nil {
		maxLifeTimeConn = parse
	}
	if parse, err := strconv.Atoi(dbConfig["maxidletimeconn"]); err == nil {
		maxIdleTimeConn = parse
	}
	sqlDb.SetMaxOpenConns(maxOpenConn)
	sqlDb.SetMaxIdleConns(maxIdleConn)
	sqlDb.SetConnMaxLifetime(time.Duration(maxLifeTimeConn) * time.Second)
	sqlDb.SetConnMaxIdleTime(time.Duration(maxIdleTimeConn) * time.Second)
	log.Println(("test"))

}
