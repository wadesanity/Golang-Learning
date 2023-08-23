package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	apiReqType "videoGo/case/wadesanity_4/api/reqtype"
	apiResType "videoGo/case/wadesanity_4/api/restype"
	"videoGo/case/wadesanity_4/pkg/e"
	"videoGo/case/wadesanity_4/pkg/util"
	"videoGo/case/wadesanity_4/service"
)

func UserRegister(c *gin.Context) {
	var userRegisterReq apiReqType.UserRegisterReq
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
	c.JSON(http.StatusOK, apiResType.UserRegisterRes{
		Status: http.StatusOK,
		Msg:    "注册用户成功",
		Data:   *user,
	})
}

func UserLogin(c *gin.Context) {
	var userLoginReq apiReqType.UserLoginReq
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
	c.JSON(http.StatusOK, apiResType.UserLoginRes{
		Status: http.StatusOK,
		Msg:    "用户登录成功",
		Token:  s,
	})
	return
}

func UserChangePwd(c *gin.Context) {
	userIDAny, ok := c.Get("userID")
	if !ok {
		util.Logger.Errorf("用户修改密码handle,token获取id错误")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.AuthorizeError.Error(),
		})
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		util.Logger.Errorf("用户修改密码handle,token解析id类型错误")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.AuthorizeError.Error(),
		})
		return
	}

	var req apiReqType.UserChangePwdReq
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
	req.ID = userID
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
	c.JSON(http.StatusOK, apiResType.UserChangePwdRes{
		Status:   http.StatusOK,
		Msg:      "修改密码成功",
		UpdateAt: *t,
	})
}
