package conf

import (
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"todolistGo/case/wadesanity_2/logger"
)

var (
	MysqlUser       string
	MysqlPwd        string
	MysqlRemoteIp   string
	MysqlRemotePort string
	MysqlDataBase   string
	MysqlLogLevel   string
)

func init() {
	abs, err := filepath.Abs("./case/wadesanity_2/conf/conf.ini")
	if err != nil {
		logger.Logger.Errorf("Fail to get filepath: %v", err)
		return
	}
	cfg, err := ini.Load(abs)
	if err != nil {
		logger.Logger.Errorf("Fail to read file: %v", err)
		os.Exit(1)
	}

	mysqlSection := cfg.Section("mysql")
	MysqlUser = mysqlSection.Key("user").String()
	MysqlPwd = mysqlSection.Key("pwd").String()
	MysqlRemoteIp = mysqlSection.Key("remote_ip").String()
	MysqlRemotePort = mysqlSection.Key("remote_port").String()
	MysqlDataBase = mysqlSection.Key("database").String()
	MysqlLogLevel = mysqlSection.Key("logLevel").String()

}
