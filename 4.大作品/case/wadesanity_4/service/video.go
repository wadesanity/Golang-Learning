package service

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
	"sync"
	"videoGo/case/wadesanity_4/pkg/e"
	"videoGo/case/wadesanity_4/pkg/util"
	"videoGo/case/wadesanity_4/repository/cache"
	"videoGo/case/wadesanity_4/repository/db"
	"videoGo/case/wadesanity_4/repository/db/model"
	typesReq "videoGo/case/wadesanity_4/types/req"
	typesRes "videoGo/case/wadesanity_4/types/res"
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
	video := &model.Video{
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
	db1 := db.NewDBClient(ctx).Model(&model.Video{}).
		Where(&model.Video{UserID: req.UserID, Status: 1})
	//Where(&model.Video{UserID: req.UserID})
	db2 := db1.Count(&total)
	err = db2.Error
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
	err = db1.Order("createdTime desc").
		Offset(*req.Offset).
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
		video.ViewCount = cache.GetVideoCount(ctx, video.ID, "viewCount")
		video.TimeCommentCount = cache.GetVideoCount(ctx, video.ID, "timeCommentCount")
	}
	res = videos
	return
}

func (*videoService) ShowOne(ctx context.Context, id, userID uint) (any, error) {
	count, err := cache.IncreaseVideoCount(ctx, id, "viewCount", 1)
	if err != nil {
		util.Logger.Errorf("videoOne cache error:%v", err)
		return nil, e.NewApiError(http.StatusInternalServerError, e.CacheUpdateError.Error())
	}
	_, _ = cache.IncreaseScoreVideoTop(ctx, strconv.Itoa(int(id)), 1)
	return typesRes.VideoOneRes{
		ViewCount:     uint(count),
		UpCount:       cache.GetVideoCount(ctx, id, "up"),
		BookmarkCount: cache.GetVideoCount(ctx, id, "bookmark"),
		ForwardCount:  cache.GetVideoCount(ctx, id, "forward"),
		IsUp:          cache.IsVideoDone(ctx, userID, id, "up"),
		IsBookMark:    cache.IsVideoDone(ctx, userID, id, "bookmark"),
		IsForward:     cache.IsVideoDone(ctx, userID, id, "forward"),
	}, nil
}

func (*videoService) ActionOne(ctx context.Context, req *typesReq.VideoActionReq) error {
	b := cache.IsVideoDone(ctx, req.UserID, req.ID, req.Action)
	switch req.Value {
	case 1:
		if b {
			return e.NewApiError(http.StatusBadRequest, e.RepeatActionError.Error())
		}
		_, err := cache.ActionVideo(ctx, req.UserID, req.Action, "add", req.ID)
		if err != nil {
			util.Logger.Errorf("video actionOne cache err:%v", err)
			return e.NewApiError(http.StatusInternalServerError, e.CacheUpdateError.Error())
		}
		return nil
	case -1:
		if !b {
			return e.NewApiError(http.StatusBadRequest, e.RepeatActionError.Error())
		}
		_, err := cache.ActionVideo(ctx, req.UserID, req.Action, "del", req.ID)
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
	comment := &model.Comment{
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
	var comments []*model.Comment
	var total int64
	db1 := db.NewDBClient(ctx).Model(&model.Comment{}).
		Where(&model.Comment{VideoID: req.VideoID, PID: *req.PID})
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
	var timeComment = &model.TimeComment{
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
	var res []*model.TimeComment
	err := db.NewDBClient(ctx).Model(&model.TimeComment{}).
		Where(&model.TimeComment{VideoID: req.VideoID}).
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
	idStringList, err := cache.GetVideoTop(ctx, 0, 99)
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
