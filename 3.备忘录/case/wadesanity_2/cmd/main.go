package main

import (
	"todolistGo/case/wadesanity_2/logger"
	"todolistGo/case/wadesanity_2/router"
)

func main() {
	e:=router.SetupRouter()
	err := e.Run()
	if err != nil {
		logger.Logger.Fatalf("gin run error:%v",err)
		return 
	}
}