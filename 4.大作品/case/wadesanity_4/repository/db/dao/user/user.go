package dao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"videoGo/case/wadesanity_4/repository/db"
	"videoGo/case/wadesanity_4/repository/db/model"
)

var (
	UserDAO     *userDAO
	userDAOOnce sync.Once
)

type UserDAOInstance interface {
	CheckExistByOpts(opts ...Option) (bool, error)
	GetByOpts(opts ...Option) (*model.User, error)
	AddNew(user *model.User) (*model.User, error)
	ChangeByModel(user *model.User) (*model.User, error)
}

type userDAO struct {
	*gorm.DB
}

func newUserDAO(ctx context.Context) *userDAO {
	return &userDAO{db.Db.WithContext(ctx)}
}

func GetUserDAO(ctx context.Context) UserDAOInstance {
	userDAOOnce.Do(func() {
		UserDAO = newUserDAO(ctx)
	})
	return UserDAO
}

func (u *userDAO) AddNew(user *model.User) (*model.User, error) {
	res := u.DB.Create(&user)
	if res.Error != nil {
		return nil, fmt.Errorf("用户表添加数据库错误:%w", res.Error)
	}
	return user, nil
}

func (u *userDAO) CheckExistByOpts(opts ...Option) (bool, error) {
	var count int64
	db1 := u.DB.Model(&model.User{})
	for _, opt := range opts {
		db1 = opt(db1)
	}
	db1 = db1.Count(&count)
	return count == 1, db1.Error
}

func (u *userDAO) GetByOpts(opts ...Option) (*model.User, error) {
	var user model.User
	db1 := u.DB.Model(&model.User{})
	for _, opt := range opts {
		db1 = opt(db1)
	}
	res := db1.Take(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("用户表查询数据库错误:%w", res.Error)
	}
	return &user, nil
}

func (u *userDAO) ChangeByModel(user *model.User) (*model.User, error) {
	res := u.DB.Updates(user)
	if res.Error != nil {
		return nil, fmt.Errorf("用户表更新数据库错误:%w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("用户表更新数据库更新数为0:%w", gorm.ErrRecordNotFound)
	}
	return user, nil
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

func WithNameInUser(userName string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", userName)
	}
}

func WithIDInUser(userID uint) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", userID)
	}
}
