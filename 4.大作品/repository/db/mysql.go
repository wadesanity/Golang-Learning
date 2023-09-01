package db

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"
	"videoGo/conf"
	"videoGo/pkg/util"
	model2 "videoGo/repository/db/model"
)

var _db *gorm.DB

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}

func Init() {
	var dsn string
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.MysqlUser,
		conf.MysqlPwd,
		conf.MysqlRemoteIp,
		conf.MysqlRemotePort,
		conf.MysqlDataBase,
	)
	var logLevel logger.LogLevel
	switch strings.ToLower(conf.MysqlLogLevel) {
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Silent
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
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
			SingularTable: true,
		},
	})
	if err != nil {
		util.Logger.Fatalf("数据库连接失败:%v", err)
		//panic(any(err))
		return
	}
	_db = db
	migration()
}

func migration() {
	err := _db.AutoMigrate(&model2.User{}, &model2.Video{}, &model2.Comment{}, &model2.TimeComment{})
	if err != nil {
		util.Logger.Fatalf("数据库迁移失败:%v", err)
		return
	}
}
