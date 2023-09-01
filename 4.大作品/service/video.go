package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"videoGo/pkg/e"
	"videoGo/pkg/util"
	cache2 "videoGo/repository/cache"
	"videoGo/repository/db"
	model2 "videoGo/repository/db/model"
	typesReq "videoGo/types/req"
	typesRes "videoGo/types/res"
)

var (
	videoServiceInstance *videoService
	videoServiceOnce     sync.Once
)

type VideoService interface {
	CreateOne(ctx context.Context, req *typesReq.VideoCreateReq) (any, error)
	List(ctx context.Context, req *typesReq.VideoListReq) (videos any, total int64, err error)
	ShowOne(ctx context.Context, id, userID uint) (any, error)
	ActionOne(ctx context.Context, req *typesReq.VideoActionReq) error
	CreateComment(ctx context.Context, req *typesReq.VideoCommentReq) (any, error)
	ListComment(ctx context.Context, req *typesReq.VideoCommentListReq) (any, int64, error)
	CreateTimeComment(ctx context.Context, req *typesReq.VideoTimeCommentCreateReq) (any, error)
	ListTimeComment(ctx context.Context, req *typesReq.VideoTimeCommentList) (any, error)
	VideoTop(ctx context.Context) (any, error)
}

func GetVideoService() VideoService {
	videoServiceOnce.Do(func() {
		videoServiceInstance = newVideoService()
	})
	return videoServiceInstance
}

func newVideoService() *videoService {
	return &videoService{}
}

type videoService struct{}

func (*videoService) CreateOne(ctx context.Context, req *typesReq.VideoCreateReq) (res any, err error) {
	video := &model2.Video{
		UserID:     req.UserID,
		Title:      req.Title,
		StaticPath: req.StaticPath,
	}
	err = db.NewDBClient(ctx).Create(&video).Error
	if err != nil {
		util.Logger.Errorf("视频添加方法->error:%v,请求形参:%#v,", err, req)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbCreateError.Error())
	}
	return video, nil
}

func (*videoService) List(ctx context.Context, req *typesReq.VideoListReq) (res any, total int64, err error) {
	db1 := db.NewDBClient(ctx).Model(&model2.Video{}).
		Where(&model2.Video{Status: 1})

	if req.Title != nil {
		db1 = db1.Where("title like ?", fmt.Sprintf("%%%s%%", *req.Title))
	}

	if req.CreateStart != nil {
		db1 = db1.Where("createdTime > ?", time.Unix(*req.CreateStart, 0))
	}

	if req.CreateEnd != nil {
		db1 = db1.Where("createdTime <= ?", time.Unix(*req.CreateEnd, 0))
	}

	err = db1.Count(&total).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.Errorf("视频列表展示方法->not found,请求形参:%#v,", req)
			return nil, total, nil
		}
		util.Logger.Errorf("视频列表展示方法->error:%v,请求形参:%#v,", err, req)
		return nil, total, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	if total <= 0 {
		util.Logger.Errorf("视频列表展示方法->not found,请求形参:%#v,", req)
		return
	}
	var videos []*typesRes.VideoListRes
	if req.Order != nil {
		orderList := strings.Split(*req.Order, ",")
		for _, order := range orderList {
			db1 = db1.Order(order)
		}
	}
	err = db1.Offset(*req.Offset).
		Limit(req.Limit).
		Find(&videos).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.Errorf("视频列表展示方法->not found,请求形参:%#v,", req)
			return nil, total, e.NewApiError(http.StatusNotFound, e.DbQueryNotFound.Error())
		}
		util.Logger.Errorf("视频列表展示方法->error:%v,请求形参:%#v,", err, req)
		return nil, total, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	for _, video := range videos {
		video.ViewCount = cache2.GetVideoCount(ctx, video.ID, "viewCount")
		video.TimeCommentCount = cache2.GetVideoCount(ctx, video.ID, "timeCommentCount")
	}
	res = videos
	return
}

func (*videoService) ShowOne(ctx context.Context, id, userID uint) (any, error) {
	count, err := cache2.IncreaseVideoCount(ctx, id, "viewCount", 1)
	if err != nil {
		util.Logger.Errorf("videoOne cache error:%v", err)
		return nil, e.NewApiError(http.StatusInternalServerError, e.CacheUpdateError.Error())
	}
	_, _ = cache2.IncreaseScoreVideoTop(ctx, strconv.Itoa(int(id)), 1)
	return typesRes.VideoOneRes{
		ViewCount:     uint(count),
		UpCount:       cache2.GetVideoCount(ctx, id, "up"),
		BookmarkCount: cache2.GetVideoCount(ctx, id, "bookmark"),
		ForwardCount:  cache2.GetVideoCount(ctx, id, "forward"),
		IsUp:          cache2.IsVideoDone(ctx, userID, id, "up"),
		IsBookMark:    cache2.IsVideoDone(ctx, userID, id, "bookmark"),
		IsForward:     cache2.IsVideoDone(ctx, userID, id, "forward"),
	}, nil
}

func (*videoService) ActionOne(ctx context.Context, req *typesReq.VideoActionReq) error {
	b := cache2.IsVideoDone(ctx, req.UserID, req.ID, req.Action)
	switch req.Value {
	case 1:
		if b {
			return e.NewApiError(http.StatusBadRequest, e.RepeatActionError.Error())
		}
		_, err := cache2.ActionVideo(ctx, req.UserID, req.Action, "add", req.ID)
		if err != nil {
			util.Logger.Errorf("video actionOne cache err:%v", err)
			return e.NewApiError(http.StatusInternalServerError, e.CacheUpdateError.Error())
		}
		return nil
	case -1:
		if !b {
			return e.NewApiError(http.StatusBadRequest, e.RepeatActionError.Error())
		}
		_, err := cache2.ActionVideo(ctx, req.UserID, req.Action, "del", req.ID)
		if err != nil {
			util.Logger.Errorf("video actionOne cache err:%v", err)
			return e.NewApiError(http.StatusInternalServerError, e.CacheUpdateError.Error())
		}
		return nil
	default:
		return e.NewApiError(http.StatusBadRequest, e.ReqParamsError.Error())
	}
}

func (*videoService) CreateComment(ctx context.Context, req *typesReq.VideoCommentReq) (any, error) {
	comment := &model2.Comment{
		UserID:  req.UserID,
		VideoID: req.VideoID,
		Content: req.Content,
		PID:     *req.PID,
	}
	err := db.NewDBClient(ctx).Create(comment).Error
	if err != nil {
		util.Logger.Errorf("createComment err:%v", err)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbCreateError.Error())
	}
	return comment, nil
}

func (*videoService) ListComment(ctx context.Context, req *typesReq.VideoCommentListReq) (any, int64, error) {
	var comments []*model2.Comment
	var total int64
	db1 := db.NewDBClient(ctx).Model(&model2.Comment{}).
		Where(&model2.Comment{VideoID: req.VideoID, PID: *req.PID})
	err := db1.Count(&total).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.Errorf("ListComment count req:%#v, not found", req)
			return nil, total, nil
		}
		util.Logger.Errorf("ListComment count req:%#v, err:%v", req, err)
		return nil, total, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	err = db1.Order(clause.OrderByColumn{Column: clause.Column{Name: "createdTime"}, Desc: true}).
		Offset(*req.Offset).Limit(req.Limit).
		Find(&comments).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.Errorf("ListComment req:%#v, not found", req)
			return nil, total, nil
		}
		util.Logger.Errorf("ListComment req:%#v, err:%v", req, err)
		return nil, total, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	return comments, total, nil
}

func (*videoService) CreateTimeComment(ctx context.Context, req *typesReq.VideoTimeCommentCreateReq) (any, error) {
	var timeComment = &model2.TimeComment{
		UserID:  req.UserID,
		VideoID: req.VideoID,
		Content: req.Content,
		VideoAt: req.VideoAt,
	}
	err := db.NewDBClient(ctx).Create(timeComment).Error
	if err != nil {
		util.Logger.Errorf("createTimeComent err:%v", err)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbCreateError.Error())
	}
	return timeComment, nil
}

func (*videoService) ListTimeComment(ctx context.Context, req *typesReq.VideoTimeCommentList) (any, error) {
	var res []*model2.TimeComment
	err := db.NewDBClient(ctx).Model(&model2.TimeComment{}).
		Where(&model2.TimeComment{VideoID: req.VideoID}).
		Find(&res, "videoTime BETWEEN ? AND ?", req.VideoStart, req.VideoEnd).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.Debugf("ListTimeComment req:%#v, not found", req)
			return res, nil
		}
		util.Logger.Errorf("ListTimeComment req:%#v, err:%v", req, err)
		return res, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	return res, nil
}

func (*videoService) VideoTop(ctx context.Context) (any, error) {
	idStringList, err := cache2.GetVideoTop(ctx, 0, 99)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			util.Logger.Errorf("VideoTop cache not exist")
			return nil, nil
		}
		util.Logger.Errorf("VideoTop cache err:%v", err)
		return nil, e.NewApiError(http.StatusInternalServerError, e.CacheQueryError.Error())
	}
	return idStringList, nil
}
