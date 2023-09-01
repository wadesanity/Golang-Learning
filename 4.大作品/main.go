package main

import (
	"videoGo/pkg/util"
	"videoGo/repository/cache"
	"videoGo/repository/db"
	"videoGo/router"
)

func main() {
	db.Init()
	cache.Init()
	e := router.SetupRouter()
	err := e.Run()
	if err != nil {
		util.Logger.Fatalf("gin run error:%v", err)
		return
	}
}
