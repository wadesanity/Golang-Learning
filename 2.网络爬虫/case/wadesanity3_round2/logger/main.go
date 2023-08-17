package logger

import (
	"log"
	"os"
)

const logPath = "fuDaSpider.log"

var Logger *log.Logger

func init() {
	Logger = &log.Logger{}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
	Logger.SetOutput(f)
	Logger.SetFlags(log.LstdFlags|log.Lshortfile)
}