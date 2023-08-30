package api

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"videoGo/case/wadesanity_4/pkg/e"
	"videoGo/case/wadesanity_4/pkg/util"
	"videoGo/case/wadesanity_4/service"
	typesReq "videoGo/case/wadesanity_4/types/req"
	typesRes "videoGo/case/wadesanity_4/types/res"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var userRegisterReq typesReq.UserRegisterReq
	err := c.ShouldBind(&userRegisterReq)
	if err != nil {
		util.Logger.Errorf("用户注册handler参数绑定错误:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.ReqParamsError.Error(),
		})
		return
	}
	util.Logger.Debugf("userRegisterReq:%v", userRegisterReq)
	userServiceInstance := service.GetUserServiceInstance()
	user, err := userServiceInstance.Register(c.Request.Context(), &userRegisterReq)
	if err != nil {
		var eApiError *e.ApiError
		if errors.As(err, &eApiError) {
			c.JSON(eApiError.HttpStatus, gin.H{
				"status": eApiError.HttpStatus,
				"error":  eApiError.Error(),
			})
			return
		}
		util.Logger.Errorf("用户注册handle未知错误:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  e.UnknowError.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, typesRes.UserRegisterRes{
		Status: http.StatusOK,
		Msg:    "注册用户成功",
		Data:   *user,
	})
}

func UserLogin(c *gin.Context) {
	var userLoginReq typesReq.UserLoginReq
	err := c.ShouldBind(&userLoginReq)
	if err != nil {
		util.Logger.Errorf("用户登录handler参数绑定错误:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.ReqParamsError.Error(),
		})
		return
	}
	util.Logger.Debugln("userLoginReq:", userLoginReq)
	i := service.GetUserServiceInstance()
	s, err := i.Login(c.Request.Context(), &userLoginReq)
	if err != nil {
		var eApiError *e.ApiError
		if errors.As(err, &eApiError) {
			c.JSON(eApiError.HttpStatus, gin.H{
				"status": eApiError.HttpStatus,
				"error":  eApiError.Error(),
			})
			return
		}
		util.Logger.Errorf("用户登录handle未知错误:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  e.UnknowError.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, typesRes.UserLoginRes{
		Status: http.StatusOK,
		Msg:    "用户登录成功",
		Token:  s,
	})
	return
}

func UserChangePwd(c *gin.Context) {
	var req typesReq.UserChangePwdReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("用户修改密码handler参数绑定错误:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.ReqParamsError.Error(),
		})
		return
	}

	util.Logger.Debugf("UserChangePwdReq:%v", req)
	userService := service.GetUserServiceInstance()
	req.ID = c.GetUint("userID")
	t, err := userService.ChangePwd(c.Request.Context(), &req)
	if err != nil {
		var eApiError *e.ApiError
		if errors.As(err, &eApiError) {
			c.JSON(eApiError.HttpStatus, gin.H{
				"status": eApiError.HttpStatus,
				"error":  eApiError.Error(),
			})
			return
		}
		util.Logger.Errorf("用户修改密码handle未知错误:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  e.UnknowError.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, typesRes.UserChangePwdRes{
		Status:   http.StatusOK,
		Msg:      "修改密码成功",
		UpdateAt: *t,
	})
}

func UserShowInfo(c *gin.Context) {
	userService := service.GetUserServiceInstance()
	user, err := userService.ShowInfo(c.Request.Context(), c.GetUint("userID"))
	if err != nil {
		var eApiError *e.ApiError
		if errors.As(err, &eApiError) {
			c.JSON(eApiError.HttpStatus, eApiError.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, e.UnknowError.Error())
		return
	}
	c.JSON(http.StatusOK, typesRes.UserRegisterRes{
		Status: http.StatusOK,
		Msg:    "用户信息查询成功",
		Data:   *user,
	})
}

func UserChangeAvatar(c *gin.Context) {
	var req typesReq.UserChangeAvatarReq
	err := c.ShouldBindWith(&req, binding.Form)
	if err != nil {
		util.Logger.Errorf("用户修改头像handler参数绑定错误:%v", err)
		c.JSON(http.StatusBadRequest, e.ReqParamsError.Error())
		return
	}

	userService := service.GetUserServiceInstance()
	req.ID = c.GetUint("userID")
	user, err := userService.ChangeAvatar(c.Request.Context(), &req)
	if err != nil {
		var eApiError *e.ApiError
		if errors.As(err, &eApiError) {
			c.JSON(eApiError.HttpStatus, eApiError.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, e.UnknowError.Error())
		return
	}
	c.JSON(http.StatusOK, typesRes.UserRegisterRes{
		Status: http.StatusOK,
		Msg:    "用户头像信息修改成功",
		Data:   *user,
	})
}

func UserList(c *gin.Context) {
	var req typesReq.UserListReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("UserList shouldBind err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	util.Logger.Debugf("req:%#v, %+v", req, req)
	s := service.GetUserServiceInstance()
	res, total, err := s.List(c.Request.Context(), &req)
	util.Logger.Debugf("1res:%#v, %v", res, res == nil)
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResList(res, total, "用户列表查询成功"))
}
