package cache

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
	"videoGo/case/wadesanity_4/pkg/util"
	"videoGo/case/wadesanity_4/repository/db"
	"videoGo/case/wadesanity_4/repository/db/model"

	"github.com/redis/go-redis/v9"
)

const (
	videoViewCountKeyPrefix     = "video:viewCount:"
	videoUpCountKeyPrefix       = "video:upCount:"
	videoBookmarkCountKeyPrefix = "video:bookmarkCount:"
	videoForwardCountKeyPrefix  = "video:forwardCount:"
	videoTimeCommentCountPrefix = "video:timeCommentCount:"
	videoViewCountTopPrefix     = "video:top:"
)

func NewVideoCountKey(prefix string, id uint) string {
	return fmt.Sprintf("%s%s", prefix, strconv.Itoa(int(id)))
}

func NewVideoTopViewCountKey(prefix string) string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%s:%v:%s:%v", prefix, year, month, day)

}

var map1 = map[string]string{
	"viewCount":        videoViewCountKeyPrefix,
	"upCount":          videoUpCountKeyPrefix,
	"bookmarkCount":    videoBookmarkCountKeyPrefix,
	"ForwardCount":     videoForwardCountKeyPrefix,
	"timeCommentCount": videoTimeCommentCountPrefix,
}

func GetVideoCount(ctx context.Context, idUint uint, filedString string) (count uint) {
	prefix, _ := map1[filedString]
	idString := NewVideoCountKey(prefix, idUint)
	val, err := Rdb.Get(ctx, idString).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			util.Logger.Infof("redis key:%s does not exist", idString)
			_, _ = IncreaseVideoCount(ctx, idUint, filedString, 0)
		} else {
			util.Logger.Errorf("redis key:%s get error:%v", idString, err)
		}
		err = db.NewDBClient(ctx).Where(&model.Video{ID: idUint}).Select(filedString).Row().Scan(&count)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				util.Logger.Infof("mysql video id:%v, get:%v not found", idString, filedString)
				return 0
			}
			util.Logger.Infof("mysql video id:%v, get:%v Error:%v", idString, filedString, err)
			return 0
		}
		if count > 0 {
			_, _ = IncreaseVideoCount(ctx, idUint, filedString, int64(count))
		}
		return
	}
	util.Logger.Debugf("redis key:%s value:%v", idString, val)
	countInt, err := strconv.Atoi(val)
	if err != nil {
		util.Logger.Errorf("redis key:%s value:%v convert to int error:%v", idString, val, err)
		Rdb.Del(ctx, idString)
	}
	return uint(countInt)
}

func IncreaseVideoCount(ctx context.Context, idUint uint, filedString string, value int64) (int64, error) {
	prefix, _ := map1[filedString]
	return Rdb.IncrBy(ctx, NewVideoCountKey(prefix, idUint), value).Result()
}

func IncreaseScoreVideoTop(ctx context.Context, videoID string, value int) (float64, error) {
	return Rdb.ZIncrBy(ctx, NewVideoTopViewCountKey(videoViewCountTopPrefix), float64(value), videoID).Result()
}

func GetVideoTop(ctx context.Context, start, end int64) ([]string, error) {
	return Rdb.ZRevRange(ctx, NewVideoTopViewCountKey(videoViewCountTopPrefix), start, end).Result()
}
