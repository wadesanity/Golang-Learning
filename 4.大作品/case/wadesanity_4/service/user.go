package service

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
	"videoGo/case/wadesanity_4/pkg/e"
	"videoGo/case/wadesanity_4/pkg/util"
	"videoGo/case/wadesanity_4/repository/db/dao"
	"videoGo/case/wadesanity_4/repository/db/model"
	types "videoGo/case/wadesanity_4/types/req"

	"gorm.io/gorm"
)

var (
	userServiceInstance *userService
	userServiceOnce     sync.Once
)

type UserService interface {
	Register(context.Context, *types.UserRegisterReq) (*model.User, error)
	Login(context.Context, *types.UserLoginReq) (string, error)
	ChangePwd(context.Context, *types.UserChangePwdReq) (*time.Time, error)
	ShowInfo(context.Context, uint) (*model.User, error)
	ChangeAvatar(context.Context, *types.UserChangeAvatarReq) (*model.User, error)
}

func GetUserServiceInstance() UserService {
	userServiceOnce.Do(func() {
		userServiceInstance = newUserService()
	})
	return userServiceInstance
}

func newUserService() *userService {
	return &userService{}
}

type userService struct{}

func (*userService) Register(ctx context.Context, req *types.UserRegisterReq) (*model.User, error) {
	userDao := dao.NewUserDAO(ctx)
	b, err := userDao.GetTotalByOpts(dao.WithNameInUser(req.Name))
	if err != nil {
		util.Logger.Errorf("用户注册方法->用户存在方法->error:%v,请求形参:%v", err, req.Name)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	if b > 0 {
		util.Logger.Errorf("用户注册方法->用户存在方法->already found,请求参数:%v", req.Name)
		return nil, e.NewApiError(http.StatusBadRequest, e.DbQueryAlreadyFound.Error())
	}
	var user = &model.User{}
	user.Name = req.Name
	user.Avatar = req.Avatar
	user.Md5sumPwd(req.Pwd)
	user, err = userDao.AddNew(user)
	if err != nil {
		util.Logger.Errorf("用户注册方法->用户添加方法错误:%v,请求形参:%#v,", err, user)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbCreateError.Error())
	}
	return user, nil
}

func (*userService) Login(ctx context.Context, req *types.UserLoginReq) (string, error) {
	userDao := dao.NewUserDAO(ctx)
	user, err := userDao.GetByOpts(dao.WithNameInUser(req.Name))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.Errorf("用户登录方法->用户查询方法->not found,请求形参:%v", req.Name)
			return "", e.NewApiError(http.StatusNotFound, e.DbQueryNotFound.Error())
		}
		util.Logger.Errorf("用户登录方法->用户查询方法->error:%v,请求形参:%v", err, req.Name)
		return "", e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	if !user.CheckPwd(req.Pwd) {
		return "", e.NewApiError(http.StatusBadRequest, e.ReqParamsError.Error())
	}

	var isAuditor bool
	if user.Role == 1 {
		isAuditor = true
	}

	if user.Status == 1 {
		return "", e.NewApiError(http.StatusForbidden, e.AuthorizeError.Error())
	}

	return util.NewJwt(user.ID, isAuditor), nil

}

func (*userService) ChangePwd(ctx context.Context, req *types.UserChangePwdReq) (*time.Time, error) {
	userDao := dao.NewUserDAO(ctx)
	user, err := userDao.GetByOpts(dao.WhereID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.Errorf("用户修改密码方法->用户查询方法->not found,用户id:%v", req.ID)
			return nil, e.NewApiError(http.StatusNotFound, e.DbQueryNotFound.Error())
		}
		util.Logger.Errorf("用户修改密码方法->用户查询方法->error:%v,用户id:%v", err, req.ID)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}

	if !user.CheckPwd(req.PwdOld) {
		util.Logger.Errorf("用户修改密码方法->密码校验->false,用户名:%v", req.ID)
		return nil, e.NewApiError(http.StatusBadRequest, e.ReqParamsError.Error())
	}
	var userNew = &model.User{
		ID: user.ID,
	}
	userNew.Md5sumPwd(req.PwdNew)
	userNew, err = userDao.ChangeByModel(userNew)
	if err != nil {
		util.Logger.Errorf("用户修改密码方法->用户修改方法->error:%v,用户model:%#v", err, userNew)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbUpdateError.Error())
	}
	return &user.UpdatedAt, nil
}

func (*userService) ShowInfo(ctx context.Context, userID uint) (*model.User, error) {
	userDao := dao.NewUserDAO(ctx)
	user, err := userDao.GetByOpts(dao.WhereID(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.Errorf("用户信息展示方法->not found,用户id:%v", userID)
			return nil, e.NewApiError(http.StatusNotFound, e.DbQueryNotFound.Error())
		}
		util.Logger.Errorf("用户信息展示方法->error:%v,用户id:%v", err, userID)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	return user, nil
}

func (*userService) ChangeAvatar(ctx context.Context, req *types.UserChangeAvatarReq) (*model.User, error) {
	userDao := dao.NewUserDAO(ctx)
	_, err := userDao.GetByOpts(dao.WhereID(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.Errorf("用户修改头像方法->查询->not found,用户id:%v", req.ID)
			return nil, e.NewApiError(http.StatusNotFound, e.DbQueryNotFound.Error())
		}
		util.Logger.Errorf("用户修改头像方法->查询->error:%v,用户id:%v", err, req.ID)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}

	var userNew = &model.User{
		ID:     req.ID,
		Avatar: req.Avatar,
	}
	userNew, err = userDao.ChangeByModel(userNew)
	if err != nil {
		util.Logger.Errorf("用户修改头像方法->修改->error:%v", err)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbUpdateError.Error())
	}
	return userNew, nil
}
