package main

import (
	"spiderGo/logger"
	"spiderGo/service"
)

func main() {
	logger.Logger.Println("start fuDa Spider")
	fs := service.GetFuDaSpiderService()
	fs.Do() //缺少分页请求所以就不并发了
	logger.Logger.Println("end fuDa Spider")
}
