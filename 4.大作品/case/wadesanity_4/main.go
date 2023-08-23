package main

import (
	"videoGo/case/wadesanity_4/pkg/util"
	"videoGo/case/wadesanity_4/router"
)


func main() {
	e:=router.SetupRouter()
	err := e.Run()
	if err != nil {
		util.Logger.Fatalf("gin run error:%v",err)
		return
	}
}
