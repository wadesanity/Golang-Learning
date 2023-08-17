package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"spiderGo/logger"
	"sync"
	"time"
)

const (
	driverName = "mysql"
	dataSourceName = "root:Liucheng_119@tcp(192.168.50.133)/test_liu"
)
var (
	Db *sql.DB
	dbOnce sync.Once
	)

func initDb() *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(any(err))
	}
	err = db.Ping()
	if err != nil {
		panic(any(err))
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	logger.Logger.Printf("initDb:%+v", db)
	return db
}

func GetDbSingleInstance() *sql.DB  {
	dbOnce.Do(func() {
		logger.Logger.Println("start initDb")
		Db=initDb()
		logger.Logger.Printf("Db:%+v", Db)
	})
	return Db
}


