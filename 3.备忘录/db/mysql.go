package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	mysqlLog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"
	"todolistGo/conf"
	"todolistGo/logger"
	"todolistGo/model"
)

var Db *gorm.DB

func Init() {
	var dsn string
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.MysqlUser,
		conf.MysqlPwd,
		conf.MysqlRemoteIp,
		conf.MysqlRemotePort,
		conf.MysqlDataBase,
	)
	var logLevel mysqlLog.LogLevel
	switch strings.ToLower(conf.MysqlLogLevel) {
	case "info":
		logLevel = mysqlLog.Info
	default:
		logLevel = mysqlLog.Silent
	}

	newLogger := mysqlLog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		mysqlLog.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:         "",
			SingularTable:       true,
			NameReplacer:        nil,
			NoLowerCase:         false,
			IdentifierMaxLength: 0,
		},
	})
	if err != nil {
		logger.Logger.Fatalf("数据库连接失败:%v", err)
		//panic(any(err))
		return
	}
	err = db.AutoMigrate(&model.User{}, &model.Task{})
	if err != nil {
		logger.Logger.Fatalf("数据库迁移失败:%v", err)
		return
	}
	Db = db
}
