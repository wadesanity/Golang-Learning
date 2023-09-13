package wrapper

import (
	"gateway/conf"
	"time"
)

var UserHystrixGroup *Group

func init() {
	UserHystrixGroup = NewServiceWrapper("user")
}

func NewServiceWrapper(name string) *Group {
	c := &Config{
		Namespace:              name,
		Timeout:                time.Duration(conf.UserServerTimeoutSecond) * time.Second,
		MaxConcurrentRequests:  conf.UserServerMaxConcurrentRequests,
		RequestVolumeThreshold: uint64(conf.UserServerRequestVolumeThreshold),
		SleepWindow:            time.Duration(conf.UserServerSleepWindowSecond) * time.Second,
		ErrorPercentThreshold:  conf.UserServerErrorPercentThreshold,
	}

	return NewGroup(c)
}
