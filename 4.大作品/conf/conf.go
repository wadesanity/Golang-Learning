package conf

import (
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"videoGo/pkg/util"
)

var (
	MysqlUser       string
	MysqlPwd        string
	MysqlRemoteIp   string
	MysqlRemotePort string
	MysqlDataBase   string
	MysqlLogLevel   string
	RedisPwd        string
	RedisRemoteIp   string
	RedisRemotePort int
	RedisDataBase   int
)

func init() {
	abs, err := filepath.Abs("./conf.ini")
	if err != nil {
		util.Logger.Errorf("Fail to get filepath: %v", err)
		return
	}
	cfg, err := ini.Load(abs)
	if err != nil {
		util.Logger.Errorf("Fail to read file: %v", err)
		os.Exit(1)
	}
	loadMysql(cfg)
	loadRedis(cfg)
}

func loadMysql(cfg *ini.File) {
	mysqlSection := cfg.Section("mysql")
	MysqlUser = mysqlSection.Key("user").String()
	MysqlPwd = mysqlSection.Key("pwd").String()
	MysqlRemoteIp = mysqlSection.Key("remote_ip").String()
	MysqlRemotePort = mysqlSection.Key("remote_port").String()
	MysqlDataBase = mysqlSection.Key("database").String()
	MysqlLogLevel = mysqlSection.Key("logLevel").String()
}

func loadRedis(cfg *ini.File) {
	redisSection := cfg.Section("redis")
	RedisPwd = redisSection.Key("pwd").String()
	RedisRemoteIp = redisSection.Key("remote_ip").String()
	RedisRemotePort = redisSection.Key("remote_port").MustInt(6379)
	RedisDataBase = redisSection.Key("database").MustInt(0)
}
