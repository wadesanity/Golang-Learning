package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"todolistGo/common"
	"todolistGo/dao"
	"todolistGo/logger"
	"todolistGo/model"
	"todolistGo/service"
)

func taskGet(c *gin.Context) {
	userIDAny, ok := c.Get("userID")
	if !ok {
		logger.Logger.Errorf("userID not exists in context")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserIDJwtError.Error(),
		})
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		logger.Logger.Errorf("userID to uint error")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserIDJwtError.Error(),
		})
		return
	}

	offsetString, ok1 := c.GetQuery("offset")
	limitString, ok2 := c.GetQuery("limit")
	key, ok3 := c.GetQuery("key")
	if !ok1 || !ok2 || !ok3 {
		logger.Logger.Errorf("task query params invalid")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.TaskQueryParamsError.Error(),
		})
		return
	}

	offset, err1 := strconv.Atoi(offsetString)
	limit, err2 := strconv.Atoi(limitString)
	if err1 != nil || err2 != nil {
		logger.Logger.Errorf("task query params invalid")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.TaskQueryParamsError.Error(),
		})
	}

	statusString, ok4 := c.GetQuery("status")
	var status int
	if !ok4 {
		status = -1
	} else {
		statusInt, err := strconv.Atoi(statusString)
		if err != nil {
			logger.Logger.Errorf("task query params invalid")
			c.JSON(http.StatusBadRequest, common.ResponseError{
				Status: http.StatusBadRequest,
				Error:  common.TaskQueryParamsError.Error(),
			})
			return
		}
		status = statusInt
	}

	i := service.GetTaskServiceInstance()
	res, err := i.SearchList(offset, limit, status, userID, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseError{
			Status: http.StatusInternalServerError,
			Error:  fmt.Errorf("查询服务器内部错误:%w", err).Error(),
		})
		return
	}
	c.JSON(http.StatusOK, *res)
}

func taskPost(c *gin.Context) {
	userIDAny, ok := c.Get("userID")
	if !ok {
		logger.Logger.Errorf("userID not exists in context")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserIDJwtError.Error(),
		})
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		logger.Logger.Errorf("userID to uint error")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserIDJwtError.Error(),
		})
		return
	}

	ui := dao.GetUserDAO()
	if !ui.CheckExistByOpts(dao.WithID(userID)) {
		logger.Logger.Errorf("userID :%v not exists", userID)
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserNotExistsError.Error(),
		})
		return
	}

	var task model.Task
	err := c.ShouldBind(&task)
	if err != nil {
		logger.Logger.Errorf("task post params err:%v", err)
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.TaskPostParamsError.Error(),
		})
		return
	}
	task.UserID = userID
	i := service.GetTaskServiceInstance()
	res, err := i.AddNewTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseError{
			Status: http.StatusInternalServerError,
			Error:  fmt.Errorf("添加任务服务器内部错误:%w", err).Error(),
		})
		return
	}
	c.JSON(http.StatusOK, *res)
}

func taskPut(c *gin.Context) {
	userIDAny, ok := c.Get("userID")
	if !ok {
		logger.Logger.Errorf("userID not exists in context")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserIDJwtError.Error(),
		})
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		logger.Logger.Errorf("userID to uint error")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserIDJwtError.Error(),
		})
		return
	}

	ui := dao.GetUserDAO()
	if !ui.CheckExistByOpts(dao.WithID(userID)) {
		logger.Logger.Errorf("userID :%v not exists", userID)
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserNotExistsError.Error(),
		})
		return
	}

	ids, ok := c.GetPostForm("ids")
	if !ok {
		logger.Logger.Errorf("task ids invalid")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.TaskPostParamsError.Error(),
		})
		return
	}
	idStringList := strings.Split(ids, ",")
	idList := make([]uint, len(idStringList))
	for i_, idString := range idStringList {
		idUint, err := strconv.Atoi(idString)
		if err != nil {
			logger.Logger.Errorf("task ids invalid")
			c.JSON(http.StatusBadRequest, common.ResponseError{
				Status: http.StatusBadRequest,
				Error:  common.TaskPostParamsError.Error(),
			})
			return
		}
		idList[i_] = uint(idUint)
	}

	//statusString, ok := c.GetPostForm("status")
	//if !ok{
	//	logger.Logger.Errorf("task status invalid")
	//	c.JSON(http.StatusBadRequest, common.ResponseError{
	//		Status: http.StatusBadRequest,
	//		Error:  common.TaskPostParamsError.Error(),
	//	})
	//	return
	//}
	//status, err := strconv.Atoi(statusString)
	//if err!=nil{
	//	logger.Logger.Errorf("task status invalid")
	//	c.JSON(http.StatusBadRequest, common.ResponseError{
	//		Status: http.StatusBadRequest,
	//		Error:  common.TaskPostParamsError.Error(),
	//	})
	//	return
	//}

	i := service.GetTaskServiceInstance()
	//res, err := i.PutTask(idList,uint(status))

	var task model.Task
	err := c.ShouldBind(&task)
	if err != nil {
		logger.Logger.Errorf("task put params err:%v", err)
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.TaskPostParamsError.Error(),
		})
		return
	}
	res, err := i.PutTaskNew(idList, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseError{
			Status: http.StatusInternalServerError,
			Error:  fmt.Errorf("修改任务状态错误:%w", err).Error(),
		})
		return
	}
	c.JSON(http.StatusOK, *res)
}

func taskDelete(c *gin.Context) {
	userIDAny, ok := c.Get("userID")
	if !ok {
		logger.Logger.Errorf("userID not exists in context")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserIDJwtError.Error(),
		})
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		logger.Logger.Errorf("userID to uint error")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserIDJwtError.Error(),
		})
		return
	}

	ui := dao.GetUserDAO()
	if !ui.CheckExistByOpts(dao.WithID(userID)) {
		logger.Logger.Errorf("userID :%v not exists", userID)
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.UserNotExistsError.Error(),
		})
		return
	}

	ids, ok := c.GetPostForm("ids")
	if !ok {
		logger.Logger.Errorf("task ids invalid")
		c.JSON(http.StatusBadRequest, common.ResponseError{
			Status: http.StatusBadRequest,
			Error:  common.TaskPostParamsError.Error(),
		})
		return
	}
	idStringList := strings.Split(ids, ",")
	idList := make([]uint, len(idStringList))
	for i_, idString := range idStringList {
		idUint, err := strconv.Atoi(idString)
		if err != nil {
			logger.Logger.Errorf("task ids invalid")
			c.JSON(http.StatusBadRequest, common.ResponseError{
				Status: http.StatusBadRequest,
				Error:  common.TaskPostParamsError.Error(),
			})
			return
		}
		idList[i_] = uint(idUint)
	}

	i := service.GetTaskServiceInstance()
	res, err := i.DeleteTask(idList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ResponseError{
			Status: http.StatusInternalServerError,
			Error:  fmt.Errorf("修改任务状态错误:%w", err).Error(),
		})
		return
	}
	c.JSON(http.StatusOK, *res)
}
