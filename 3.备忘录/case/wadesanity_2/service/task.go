package service

import (
	"sync"
	"todolistGo/case/wadesanity_2/common"
	"todolistGo/case/wadesanity_2/dao"
	"todolistGo/case/wadesanity_2/logger"
	"todolistGo/case/wadesanity_2/model"
)

var (
	TaskServiceInstance *taskService
	taskServiceOnce     sync.Once
)

type TaskService interface {
	SearchList(offset, limit, status int, userID uint, key string) (*common.ResponseOk, error)
	AddNewTask(task model.Task) (*common.ResponseOk, error)
	PutTask(taskIDList []uint, status uint) (*common.ResponseOk, error)
	PutTaskNew(taskIDList []uint, task model.Task) (*common.ResponseOk, error)

	DeleteTask(taskIDList []uint) (*common.ResponseOk, error)
}

func GetTaskServiceInstance() TaskService {
	taskServiceOnce.Do(func() {
		TaskServiceInstance = newTaskService()
	})
	return TaskServiceInstance
}

type taskService struct {
}

func newTaskService() *taskService {
	return &taskService{}
}

func (taskService) SearchList(offset, limit, status int, userID uint, key string) (*common.ResponseOk, error) {
	i := dao.GetTaskDAO()
	var taskList []model.Task
	var total int64
	var err error
	if status == -1 {
		taskList, total, err = i.SearchByKey(offset, limit,
			dao.WithUserIDAndOtherOr(userID, key),
		)
	} else {
		taskList, total, err = i.SearchByKey(offset, limit,
			dao.WithUserIDAndOtherOr(userID, key),
			dao.WithTaskStatus(uint(status)),
		)
	}
	if err != nil {
		return nil, err
	}
	item := make([]map[string]interface{}, 0)
	for _, task := range taskList {
		r := make(map[string]interface{})
		r["id"] = task.ID
		r["userID"] = task.UserID
		r["title"] = task.Title
		r["content"] = task.Content
		r["view"] = task.View
		r["Status"] = task.Status
		r["createdTime"] = task.CreatedAt
		r["updatedTime"] = task.UpdatedAt
		item = append(item, r)
	}

	data := common.Data{
		Item:  item,
		Total: int(total),
	}

	return &common.ResponseOk{
		Status: 200,
		Data:   data,
		Msg:    "查询任务成功",
	}, nil
}

func (taskService) AddNewTask(task model.Task) (*common.ResponseOk, error) {
	logger.Logger.Println("add task form:%v", task)
	i := dao.GetTaskDAO()
	err := i.AddNew(&task)
	if err != nil {
		return nil, err
	}
	r := make(map[string]interface{})
	r["id"] = task.ID
	r["userID"] = task.UserID
	r["Title"] = task.Title
	r["Content"] = task.Content
	r["View"] = task.View
	r["Status"] = task.Status
	r["createdTime"] = task.CreatedAt
	r["updatedTime"] = task.UpdatedAt
	item := make([]map[string]interface{}, 0)
	item = append(item, r)
	data := common.Data{
		Item:  item,
		Total: 1,
	}

	return &common.ResponseOk{
		Status: 200,
		Data:   data,
		Msg:    "添加任务成功",
	}, nil
}

func (taskService) PutTask(taskIDList []uint, status uint) (*common.ResponseOk, error) {
	logger.Logger.Printf("put task status:%v, idList:%v,", status, taskIDList)
	i := dao.GetTaskDAO()
	success, err := i.PutStatus(taskIDList, status)
	if err != nil {
		return nil, err
	}
	data := common.Data{
		Item:  nil,
		Total: int(success),
	}

	return &common.ResponseOk{
		Status: 200,
		Data:   data,
		Msg:    "更新任务成功",
	}, nil
}

func (taskService) PutTaskNew(taskIDList []uint, task model.Task) (*common.ResponseOk, error) {
	logger.Logger.Printf("put task task:%v, idList:%v,", task, taskIDList)
	i := dao.GetTaskDAO()
	success, err := i.PutStatusNew(taskIDList, task)
	if err != nil {
		return nil, err
	}
	data := common.Data{
		Item:  nil,
		Total: int(success),
	}

	return &common.ResponseOk{
		Status: 200,
		Data:   data,
		Msg:    "更新任务成功",
	}, nil
}

func (taskService) DeleteTask(taskIDList []uint) (*common.ResponseOk, error) {
	logger.Logger.Printf("delete task idList:%v", taskIDList)
	i := dao.GetTaskDAO()
	success, err := i.Delete(taskIDList)
	if err != nil {
		return nil, err
	}
	data := common.Data{
		Item:  nil,
		Total: int(success),
	}

	return &common.ResponseOk{
		Status: 200,
		Data:   data,
		Msg:    "更新任务成功",
	}, nil
}
