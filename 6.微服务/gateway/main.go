package main

import (
	"gateway/pkg/util"
	"gateway/router"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	router.InitRouter(engine)
	err := engine.Run(":8080")
	if err != nil {
		util.Logger.Panicf("engine run err:%v", err)
		return
	}
}
