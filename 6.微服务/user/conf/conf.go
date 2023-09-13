package conf

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"path/filepath"
)

var (
	MysqlUser              string
	MysqlPwd               string
	MysqlRemoteIp          string
	MysqlRemotePort        string
	MysqlDataBase          string
	MysqlLogLevel          string
	RedisPwd               string
	RedisRemoteIp          string
	RedisRemotePort        int
	RedisDataBase          int
	UserServerAddr         string
	UserListenAddr         string
	UserVisitAddr          string
	LogLevel               logrus.Level
	EtcdAddr               string
	UserTimeoutMillisecond int
)

func init() {
	abs, err := filepath.Abs("./conf/conf.ini")
	if err != nil {
		panic(fmt.Errorf("fail to get filepath: %w", err))
	}
	cfg, err := ini.Load(abs)
	if err != nil {
		panic(fmt.Errorf("fail to read file: %w", err))
	}
	loadMysql(cfg)
	//loadRedis(cfg)
	loadUser(cfg)
	loadLog(cfg)
	loadEtcd(cfg)

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

func loadUser(cfg *ini.File) {
	userSection := cfg.Section("user")
	UserServerAddr = userSection.Key("server_addr").String()
	UserListenAddr = userSection.Key("listen_addr").MustString(":8080")
	UserVisitAddr = userSection.Key("visit_addr").MustString("user:8080")
	UserTimeoutMillisecond = userSection.Key("timeout_millisecond").MustInt(5000)

}

func loadLog(cfg *ini.File) {
	logSection := cfg.Section("log")
	LogLevel = logrus.Level(logSection.Key("log_level").MustInt(5))
}

func loadEtcd(cfg *ini.File) {
	etcdSection := cfg.Section("etcd")
	EtcdAddr = etcdSection.Key("addr").MustString("http://etcd:2379")
}
