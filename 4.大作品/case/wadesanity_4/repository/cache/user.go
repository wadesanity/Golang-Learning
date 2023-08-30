package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"videoGo/case/wadesanity_4/pkg/util"
	"videoGo/case/wadesanity_4/repository/db"
	"videoGo/case/wadesanity_4/repository/db/model"
)

const (
	UserUpVideoIDKeyPrefix       = "user:UpVideoIDSet:"
	UserBookmarkVideoIDKeyPrefix = "user:BookmarkVideoIDSet:"
	UserForwardVideoIDKeyPrefix  = "user:ForwardVideoIDSet:"
	UserBlackUserIDKey           = "user:BlackUserIDBitmap"
)

func NewUserKey(prefix string, userID uint) string {
	return fmt.Sprintf("%s%s", prefix, strconv.Itoa(int(userID)))
}

var map2 = map[string]string{
	"up":       UserUpVideoIDKeyPrefix,
	"bookmark": UserBookmarkVideoIDKeyPrefix,
	"forward":  UserForwardVideoIDKeyPrefix,
}

func IsVideoDone(ctx context.Context, userID, videoID uint, actionString string) bool {
	prefix, _ := map2[actionString]
	idString := NewUserKey(prefix, userID)
	b, err := Rdb.SIsMember(ctx, idString, videoID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			util.Logger.Infof("redis key:%s does not exist", idString)
			_, _ = ActionVideo(ctx, userID, actionString, "add", 0)
		} else {
			util.Logger.Errorf("redis key:%s get error:%v", idString, err)
		}
		var actionListString string
		err = db.NewDBClient(ctx).Where(&model.User{ID: userID}).Select(actionString).Take(&actionListString).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				util.Logger.Infof("mysql userid id:%v, get:%v not found", userID, actionString)
				return false
			}
			util.Logger.Infof("mysql userid id:%v, get:%v error:%v", userID, actionString, err)
			return false
		}
		actionList := strings.Split(actionListString, ",")
		if len(actionList) > 0 {
			_, _ = ActionVideo(ctx, userID, actionString, "add", actionList)
		}
	}
	util.Logger.Debugf("redis key:%s value:%v", idString, b)
	return b
}

func ActionVideo(ctx context.Context, userID uint, actionString, action string, value any) (int64, error) {
	prefix, _ := map2[actionString]
	keyString := NewUserKey(prefix, userID)
	if action == "del" {
		return Rdb.SRem(ctx, keyString, value).Result()
	}
	return Rdb.SAdd(ctx, keyString, value).Result()
}

func GetBlackBitUserID(ctx context.Context, userID uint) (int64, error) {
	return Rdb.GetBit(ctx, UserBlackUserIDKey, int64(userID)).Result()
}

func SetBlackBitUserID(ctx context.Context, userID uint, value int) (int64, error) {
	return Rdb.SetBit(ctx, UserBlackUserIDKey, int64(userID), value).Result()
}
