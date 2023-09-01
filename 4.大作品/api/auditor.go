package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"videoGo/pkg/e"
	"videoGo/pkg/util"
	"videoGo/service"
	types "videoGo/types/req"
	typesRes "videoGo/types/res"
)

func AuditorVideo(c *gin.Context) {
	var req types.AuditorReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("AuditorVideo shouldBind err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	util.Logger.Debugf("req:%v", req)
	s := service.GetAuditorService()
	res, err := s.VideoAudit(c.Request.Context(), &req)
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResOk("审计视频成功", http.StatusOK, res))
}

func AuditorUser(c *gin.Context) {
	var req types.AuditorReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("AuditorUser shouldBind err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	util.Logger.Debugf("req:%v", req)
	s := service.GetAuditorService()
	res, err := s.UserAudit(c.Request.Context(), &req)
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResOk("审计用户成功", http.StatusOK, res))
}

func AuditorComment(c *gin.Context) {
	var req types.AuditorReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("AuditorComment shouldBind err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	util.Logger.Debugf("req:%v", req)
	s := service.GetAuditorService()
	res, err := s.CommentAudit(c.Request.Context(), &req)
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResOk("审计评论成功", http.StatusOK, res))
}
