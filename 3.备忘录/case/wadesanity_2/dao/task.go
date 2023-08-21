package dao

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
	"todolistGo/case/wadesanity_2/db"
	"todolistGo/case/wadesanity_2/model"
)

var (
	TaskDAO     *taskDAO
	taskDAOOnce sync.Once
)

type TaskDAOInstance interface {
	SearchByKey(offset, limit int, opts ...Option) (taskList []model.Task, total int64, err error)
	AddNew(task *model.Task) error
	PutStatus(taskIDList []uint, status uint) (success int64, err error)
	PutStatusNew(taskIDList []uint, task model.Task) (success int64, err error)
	Delete(taskIDList []uint) (success int64, err error)
}

func GetTaskDAO() TaskDAOInstance {
	taskDAOOnce.Do(func() {
		TaskDAO = newTaskDAO()
	})
	return TaskDAO
}

type taskDAO struct{}

func newTaskDAO() *taskDAO {
	return &taskDAO{}
}

func (taskDAO) SearchByKey(offset, limit int, opts ...Option) (taskList []model.Task, total int64, err error) {
	db1 := db.Db.Model(&model.Task{}).Debug()
	for _, opt := range opts {
		db1 = opt(db1)
	}
	db1 = db1.Count(&total)
	db1 = db1.Order("updatedTime DESC").Offset(offset).Limit(limit)
	db1 = db1.Find(&taskList)
	if db1.RowsAffected == 0 {
		return taskList, total, db1.Error
	}

	taskIDList := make([]uint, db1.RowsAffected)
	for i, task := range taskList {
		taskList[i].View += 1
		taskIDList[i] = task.ID
	}
	db2 := db.Db.Model(&model.Task{}).Debug()
	db2 = db2.Where("id IN ?", taskIDList)
	db2 = db2.UpdateColumn("view", gorm.Expr("view + ?", 1))

	return taskList, total, db2.Error

}

func (taskDAO) AddNew(task *model.Task) error {
	db1 := db.Db.Model(&model.Task{}).Debug()
	//var taskNew model.Task
	//taskNew.Title = task.Title
	//taskNew.Content = task.Content
	db1 = db1.Create(task)
	if db1.RowsAffected != 1 {
		return fmt.Errorf("插入任务失败:%v, err:%w", task, db1.Error)
	}
	return nil
}

func (taskDAO) PutStatus(taskIDList []uint, status uint) (success int64, err error){
	db1:= db.Db.Model(&model.Task{}).Debug()
	db1 = db1.Where("id IN ?", taskIDList)
	//db1 = db1.UpdateColumn("status", status)
	db1 = db1.Updates(model.Task{
		Status:    &status,
		//UpdatedAt: time.Time{},
	})

	return db1.RowsAffected, db1.Error
}

func (taskDAO) PutStatusNew(taskIDList []uint, task model.Task) (success int64, err error) {
	db1:= db.Db.Model(&model.Task{}).Debug()
	db1 = db1.Where("id IN ?", taskIDList)
	//db1 = db1.UpdateColumn("status", status)
	db1 = db1.Updates(task)
	return db1.RowsAffected, db1.Error
}

func (taskDAO) Delete(taskIDList []uint) (success int64, err error){
	db1:= db.Db.Debug()
	db1 = db1.Where("id IN ?", taskIDList)
	//db1 = db1.UpdateColumn("status", status)
	db1 = db1.Delete(&model.Task{})
	return db1.RowsAffected, db1.Error
}

type Option func(db *gorm.DB) *gorm.DB

func WithUserID(userID uint) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("userID = ?", userID)
	}
}

func WithTaskStatus(status uint) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}


func WithOrTaskTitle(title string) Option {
	return func(db *gorm.DB) *gorm.DB {
		ftitle := fmt.Sprintf("%%%s%%", title)
		return db.Or("title like ?", ftitle)
	}
}

func WithOrTaskContent(Content string) Option {
	return func(db *gorm.DB) *gorm.DB {
		fContent := fmt.Sprintf("%%%s%%", Content)
		return db.Or("content like ?", fContent)
	}
}

func WithUserIDAndOtherOr(userID uint, key string) Option {
	return func(db *gorm.DB) *gorm.DB {
		fKey := fmt.Sprintf("%%%s%%", key)
		return db.Where("userID = ? AND (title like ? OR content like ?)", userID, fKey, fKey)
	}
}
