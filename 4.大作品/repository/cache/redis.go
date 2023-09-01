package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"videoGo/conf"
	"videoGo/pkg/util"
)

//func main() {
//	ctx := context.Background()
//
//	rdb := redis.NewClient(&redis.Options{
//		Addr:	  "localhost:6379",
//		Password: "", // no password set
//		DB:		  0,  // use default DB
//	})
//
//	err := rdb.Set(ctx, "key", "value", 0).Err()
//	if err != nil {
//		panic(err)
//	}
//
//	val, err := rdb.Get(ctx, "key").Result()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("key", val)
//
//	val2, err := rdb.Get(ctx, "key2").Result()
//	if err == redis.Nil {
//		fmt.Println("key2 does not exist")
//	} else if err != nil {
//		panic(err)
//	} else {
//		fmt.Println("key2", val2)
//	}
//	// Output: key value
//	// key2 does not exist
//}

var Rdb *redis.Client

func Init() {
	addr := fmt.Sprintf("%s:%d", conf.RedisRemoteIp, conf.RedisRemotePort)
	util.Logger.Infof("addr:%v", addr)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.RedisPwd,      // no password set
		DB:       conf.RedisDataBase, // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		util.Logger.Errorf("redis连接异常:%v", err)
		panic(any(err))
	}
	Rdb = rdb
}
