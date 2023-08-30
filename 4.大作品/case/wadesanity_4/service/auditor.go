package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"sync"
	"videoGo/case/wadesanity_4/pkg/e"
	"videoGo/case/wadesanity_4/pkg/util"
	"videoGo/case/wadesanity_4/repository/cache"
	"videoGo/case/wadesanity_4/repository/db"
	"videoGo/case/wadesanity_4/repository/db/model"
	typesReq "videoGo/case/wadesanity_4/types/req"
)

var (
	auditorServiceInstance *auditorService
	auditorServiceOnce     sync.Once
)

type AuditorService interface {
	VideoAudit(ctx context.Context, req *typesReq.AuditorReq) (any, error)
	UserAudit(ctx context.Context, req *typesReq.AuditorReq) (any, error)
	CommentAudit(ctx context.Context, req *typesReq.AuditorReq) (any, error)
}

func GetAuditorService() AuditorService {
	auditorServiceOnce.Do(func() {
		auditorServiceInstance = newAuditorService()
	})
	return auditorServiceInstance
}

func newAuditorService() *auditorService {
	return &auditorService{}
}

type auditorService struct{}

func (*auditorService) VideoAudit(ctx context.Context, req *typesReq.AuditorReq) (any, error) {
	var video model.Video
	db1 := db.NewDBClient(ctx).Model(&model.Video{})
	err := db1.Where("id = ?", req.ID).Select("status").Take(&video).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewApiError(http.StatusNotFound, e.DbQueryNotFound.Error())
		}
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	if video.Status == *req.Status {
		return nil, e.NewApiError(http.StatusBadRequest, e.RepeatActionError.Error())
	}
	video.Status = *req.Status
	err = db1.Updates(&video).Error
	if err != nil {
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbUpdateError.Error())
	}
	return &video, nil
}

func (*auditorService) UserAudit(ctx context.Context, req *typesReq.AuditorReq) (any, error) {
	var user model.User
	db1 := db.NewDBClient(ctx).Model(&model.User{})
	err := db1.Where("id = ?", req.ID).Select("status").Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewApiError(http.StatusNotFound, e.DbQueryNotFound.Error())
		}
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	if user.Status == *req.Status {
		return nil, e.NewApiError(http.StatusBadRequest, e.RepeatActionError.Error())
	}
	user.Status = *req.Status
	util.Logger.Debugf("user:%v", user)
	err = db1.Updates(&user).Error
	if err != nil {
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbUpdateError.Error())
	}
	_, _ = cache.SetBlackBitUserID(ctx, req.ID, int(*req.Status))
	return &user, nil
}

func (*auditorService) CommentAudit(ctx context.Context, req *typesReq.AuditorReq) (any, error) {
	tx := db.NewDBClient(ctx).Model(&model.Comment{}).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbUpdateError.Error())
	}
	var status uint
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.ID).Select("status").Row().Scan(&status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, e.NewApiError(http.StatusNotFound, e.DbQueryNotFound.Error())
		}
		tx.Rollback()
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	if status == *req.Status {
		tx.Rollback()
		return nil, e.NewApiError(http.StatusBadRequest, e.RepeatActionError.Error())
	}
	err = tx.Where("id = ?", req.ID).Select("status").Updates(&model.Comment{Status: *req.Status}).Error
	if err != nil {
		tx.Rollback()
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbUpdateError.Error())
	}
	err = tx.Commit().Error
	if err != nil {
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbUpdateError.Error())
	}
	return &status, nil
}
