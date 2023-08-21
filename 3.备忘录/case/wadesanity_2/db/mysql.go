package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"todolistGo/case/wadesanity_2/logger"
	"todolistGo/case/wadesanity_2/model"
)

const dsn = "root:Liucheng_119@tcp(192.168.50.133:3306)/test_liu?charset=utf8mb4&parseTime=True&loc=Local"

var Db *gorm.DB

func init() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
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
