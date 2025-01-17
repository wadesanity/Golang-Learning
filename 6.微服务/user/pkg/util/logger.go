package util

import (
	"github.com/sirupsen/logrus"
	"user/conf"
)

var Logger *logrus.Logger

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Logger = logrus.New()
	//file, err := os.OpenFile("./log/log1.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	//if err != nil {
	//	panic(fmt.Errorf("openfile err:%w", err))
	//}
	//Logger.SetOutput(file)
	Logger.SetNoLock()
	// Only log the warning severity or above.
	Logger.SetLevel(conf.LogLevel)
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   false,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
