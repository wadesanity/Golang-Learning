package service

import (
	"context"
	"net/http"
	"sync"
	"time"
	apiReqType "videoGo/case/wadesanity_4/api/reqtype"
	"videoGo/case/wadesanity_4/pkg/e"
	"videoGo/case/wadesanity_4/pkg/util"
	daoUser "videoGo/case/wadesanity_4/repository/db/dao/user"
	"videoGo/case/wadesanity_4/repository/db/model"
)

type UserService interface {
	Register(context.Context, *apiReqType.UserRegisterReq) (*model.User, error)
	Login(context.Context, *apiReqType.UserLoginReq) (string, error)
	ChangePwd(context.Context, *apiReqType.UserChangePwdReq) (*time.Time, error)
}

func GetUserServiceInstance() UserService {
	userServiceOnce.Do(func() {
		UserServiceInstance = newUserService()
	})
	return UserServiceInstance
}

var (
	UserServiceInstance *userService
	userServiceOnce     sync.Once
)

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (*userService) Register(ctx context.Context,req *apiReqType.UserRegisterReq) (*model.User, error) {
	userDao := daoUser.GetUserDAO(ctx)
	b, err := userDao.CheckExistByOpts(daoUser.WithNameInUser(req.Name))
	if err != nil {
		util.Logger.Errorf("用户注册方法->用户存在方法->error:%v,请求形参:%v", err, req.Name)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	if b {
		util.Logger.Errorf("用户注册方法->用户存在方法->already found,请求参数:%v", req.Name)
		return nil, e.NewApiError(http.StatusBadRequest, e.DbQueryAlreadyFound.Error())
	}
	var user = &model.User{}
	user.Name = req.Name
	user.Avatar = req.Avatar
	user.Md5sumPwd(req.Pwd)
	user, err = userDao.AddNew(user)
	if err != nil {
		util.Logger.Errorf("用户注册方法->用户添加方法错误:%v,请求形参:%#v,",err,user)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbCreateError.Error())
	}
	return user, nil
}

func (*userService) Login(ctx context.Context, req *apiReqType.UserLoginReq) (string, error) {
	userDao := daoUser.GetUserDAO(ctx)
	user, err := userDao.GetByOpts(daoUser.WithNameInUser(req.Name))
	if err != nil {
		util.Logger.Errorf("用户登录方法->用户查询方法->error:%v,请求形参:%v", err, req.Name)
		return "", e.NewApiError(http.StatusInternalServerError,e.DbQueryError.Error())
	}
	if user == nil {
		util.Logger.Errorf("用户登录方法->用户查询方法->not found,请求形参:%v", req.Name)
		return "", e.NewApiError(http.StatusNotFound, e.DbQueryNotFound.Error())
	}
	if !user.CheckPwd(req.Pwd) {
		return "", e.NewApiError(http.StatusBadRequest, e.ReqParamsError.Error())
	}

	return util.NewJwt(user.ID), nil

}

func (*userService) ChangePwd(ctx context.Context, req *apiReqType.UserChangePwdReq) (*time.Time, error) {
	userDao := daoUser.GetUserDAO(ctx)
	user, err := userDao.GetByOpts(daoUser.WithIDInUser(req.ID))
	if err != nil {
		util.Logger.Errorf("用户修改密码方法->用户查询方法->error:%v,用户id:%v", err, req.ID)
		return nil, e.NewApiError(http.StatusInternalServerError, e.DbQueryError.Error())
	}
	if user == nil {
		util.Logger.Errorf("用户修改密码方法->用户查询方法->not found,用户id:%v", req.ID)
		return nil, e.NewApiError(http.StatusNotFound, e.DbQueryNotFound.Error())
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
		return nil, e.NewApiError(http.StatusInternalServerError,e.DbUpdateError.Error())
	}
	return &user.UpdatedAt, nil
}
