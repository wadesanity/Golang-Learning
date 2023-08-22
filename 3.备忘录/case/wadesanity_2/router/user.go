package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"todolistGo/case/wadesanity_2/common"
	"todolistGo/case/wadesanity_2/logger"
	"todolistGo/case/wadesanity_2/service"
)



// ShowAccount godoc
// @Summary      register an account
// @Description  register by username and userpwd
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  common.ResponseOk
// @Failure      400  {object}  common.ResponseError
// @Failure      404  {object}  common.ResponseError
// @Failure      500  {object}  common.ResponseError
// @Router       /user/register [post]
func userRegister(c *gin.Context) {
	userName:=c.PostForm("userName")
	userPwd:=c.PostForm("userPwd")
	logger.Logger.Printf("userName:%v,userPwd:%v", userName,userPwd)
	if userName == "" || userPwd == "" {
		c.JSON(http.StatusBadRequest,common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserRegisterParamsError.Error(),
		})
		return
	}

	i:= service.GetUserServiceInstance()
	res, err := i.Register(userName, userPwd)
	if err != nil {
		if errors.Is(err, common.UserAlreadyExistsError){
			logger.Logger.Errorf("UserAlreadyExistsError:%v", err)
			c.JSON(http.StatusBadRequest,common.ResponseError{
				Status: http.StatusBadRequest,
				Error:  fmt.Errorf("注册用户已存在:%w",err).Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError,common.ResponseError{
			Status: http.StatusInternalServerError,
			Error:  fmt.Errorf("注册用户服务器内部错误:%w",err).Error(),
		})
		return
	}
	c.JSON(http.StatusOK, *res)
}

func userLogin(c *gin.Context)  {
	userName:=c.PostForm("userName")
	userPwd:=c.PostForm("userPwd")
	logger.Logger.Printf("userName:%v,userPwd:%v", userName,userPwd)
	if userName == "" || userPwd == "" {
		c.JSON(http.StatusBadRequest,common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserLoginParamsError.Error(),
		})
		return
	}

	i:= service.GetUserServiceInstance()
	res, err := i.Login(userName, userPwd)
	if err != nil {
		c.JSON(http.StatusInternalServerError,common.ResponseError{
			Status: http.StatusInternalServerError,
			Error:  fmt.Errorf("用户登录服务器内部错误:%w",err).Error(),
		})
		return
	}
	c.JSON(http.StatusOK, *res)
}