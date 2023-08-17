package main

import (
	"spiderGo/logger"
	"spiderGo/service"
)

func main() {
	logger.Logger.Println("start fuDa Spider")
	fs := service.GetFuDaSpiderService()
	fs.Do()
	logger.Logger.Println("end fuDa Spider")
}
