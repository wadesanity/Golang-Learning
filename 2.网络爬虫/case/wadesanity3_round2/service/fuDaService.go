package service

import (
	"spiderGo/dao"
	"spiderGo/fuDaSpider"
	"spiderGo/logger"
	"sync"
)

var (
	fs *fuDaSpiderService
	fuDaSpiderServiceOnce sync.Once
)

type fuDaSpiderService struct {
}

func (* fuDaSpiderService) Do() {
	fuDaSpiderInstance := fuDaSpider.GetFuDaSpider()
	rList, err := fuDaSpiderInstance.Spider()
	if err != nil{
		logger.Logger.Fatalln(err)
		return
	}
	fuDaNewsDaoInstance:= dao.GetFuDaNewsDao()
	totalCount, successCount, err := fuDaNewsDaoInstance.Insert(rList)
	if err != nil {
		logger.Logger.Fatalln(err)
		return
	}
	logger.Logger.Printf("Insert end, totalCount:%v, successCount:%v", totalCount, successCount)
	return
}

func newFuDaSpiderService() *fuDaSpiderService {
	return &fuDaSpiderService{}
}

type FuDaSpiderService interface {
	Do()
}

func GetFuDaSpiderService() FuDaSpiderService {
	fuDaSpiderServiceOnce.Do(func() {
		fs = newFuDaSpiderService()
	})
	return fs
}
