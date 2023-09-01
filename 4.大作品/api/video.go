package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"videoGo/pkg/e"
	"videoGo/pkg/util"
	"videoGo/service"
	typesReq "videoGo/types/req"
	typesRes "videoGo/types/res"
)

func VideoCreateOne(c *gin.Context) {
	var req typesReq.VideoCreateReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("上传视频参数绑定->error:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.ReqParamsError.Error(),
		})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		util.Logger.Errorf("上传视频参数文件获取->error:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.ReqParamsError.Error(),
		})
		return
	}

	util.Logger.Debugln("fileName:", file.Filename)
	abs, err := filepath.Abs("static")
	if err != nil {
		util.Logger.Errorf("get abs static path error:%v", err)
		c.JSON(http.StatusInternalServerError, e.UnknowError.Error())
		return
	}
	dst := path.Join(abs, file.Filename)
	util.Logger.Debugln("file dst:", dst)

	// 上传文件至指定的完整文件路径
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		util.Logger.Errorf("上传视频参数文件保存->error:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.UnknowError.Error(),
		})
		return
	}
	req.StaticPath = dst
	req.UserID = c.GetUint("userID")
	videoService := service.GetVideoService()
	video, err := videoService.CreateOne(c.Request.Context(), &req)
	if err != nil {
		var eApiError *e.ApiError
		if errors.As(err, &eApiError) {
			c.JSON(eApiError.HttpStatus, eApiError.Error())
		}
		c.JSON(http.StatusInternalServerError, e.UnknowError.Error())
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResOk("视频上传成功", http.StatusOK, video))
}

func VideoList(c *gin.Context) {
	var req typesReq.VideoListReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("VideoList shouldBind err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	//req.UserID = c.GetUint("userID")
	videoService := service.GetVideoService()
	res, total, err := videoService.List(c.Request.Context(), &req)
	if err != nil {
		var eApiError *e.ApiError
		if errors.As(err, &eApiError) {
			c.JSON(eApiError.HttpStatus, typesRes.NewResError(eApiError.HttpStatus, eApiError))
			return
		}
		c.JSON(http.StatusInternalServerError, typesRes.NewResError(http.StatusInternalServerError, e.UnknowError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResList(res, total, "视频列表查询成功"))
}

func VideoOne(c *gin.Context) {
	userID := c.GetUint("userID")
	idString := c.Query("id")
	if idString == "" {
		util.Logger.Errorf("VideoOne id not exists")
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		util.Logger.Errorf("VideoOne id convert to int err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	id := uint(idInt)
	s := service.GetVideoService()
	res, err := s.ShowOne(c.Request.Context(), id, userID)
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResOk("视频详情成功", http.StatusOK, res))
}

func VideoAction(c *gin.Context) {
	var req typesReq.VideoActionReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("VideoAction shouldbind err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	req.UserID = c.GetUint("userID")
	s := service.GetVideoService()
	err = s.ActionOne(c.Request.Context(), &req)
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResOk("操作视频成功", http.StatusOK, nil))
}

func VideoComment(c *gin.Context) {
	var req typesReq.VideoCommentReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("VideoComment shouldBind err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	util.Logger.Debugf("VideoComment req:%v", req)
	req.UserID = c.GetUint("userID")
	s := service.GetVideoService()
	res, err := s.CreateComment(c.Request.Context(), &req)
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResOk("评论成功", http.StatusOK, res))
}

func VideoCommentList(c *gin.Context) {
	var req typesReq.VideoCommentListReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("VideoCommentList shouldBind err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	util.Logger.Debugf("VideoComment req:%v", req)
	s := service.GetVideoService()
	res, total, err := s.ListComment(c.Request.Context(), &req)
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResList(res, total, "展示评论列表成功"))
}

func VideoTimeCommentCreateOne(c *gin.Context) {
	var req typesReq.VideoTimeCommentCreateReq
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("VideoTimeCommentCreateOne shouldBind err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	req.UserID = c.GetUint("userID")
	util.Logger.Debugf("VideoTimeCommentCreateOne req:%v", req)
	s := service.GetVideoService()
	res, err := s.CreateTimeComment(c.Request.Context(), &req)
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResOk("添加弹幕成功", http.StatusOK, res))
}

func VideoTimeCommentList(c *gin.Context) {
	var req typesReq.VideoTimeCommentList
	err := c.ShouldBind(&req)
	if err != nil {
		util.Logger.Errorf("VideoTimeCommentList shouldBind err:%v", err)
		c.JSON(http.StatusBadRequest, typesRes.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	util.Logger.Debugf("VideoTimeCommentList req:%v", req)
	s := service.GetVideoService()
	res, err := s.ListTimeComment(c.Request.Context(), &req)
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResList(res, 0, "展示弹幕列表成功"))
}

func VideoTop(c *gin.Context) {
	s := service.GetVideoService()
	res, err := s.VideoTop(c.Request.Context())
	if err != nil {
		var apiError *e.ApiError
		errors.As(err, &apiError)
		c.JSON(apiError.HttpStatus, typesRes.NewResError(apiError.HttpStatus, apiError))
		return
	}
	c.JSON(http.StatusOK, typesRes.NewResList(res, 0, "视频播放排行榜成功"))
}
