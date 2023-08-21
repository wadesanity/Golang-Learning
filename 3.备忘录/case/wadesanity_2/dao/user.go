package dao

import (
	"crypto/md5"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"todolistGo/case/wadesanity_2/common"
	"todolistGo/case/wadesanity_2/db"
	"todolistGo/case/wadesanity_2/logger"
	"todolistGo/case/wadesanity_2/model"
)

var (
	UserDAO     *userDAO
	userDAOOnce sync.Once
)

type UserDAOInstance interface {
	CheckExistByOpts(opts ...Option) bool
	GetByOpts(opts ...Option) *model.User
	AddNew(userName, userPwd string) (model.User, error)
}

type userDAO struct {
}

func newUserDAO() *userDAO {
	return &userDAO{}
}

func GetUserDAO() UserDAOInstance {
	userDAOOnce.Do(func() {
		UserDAO = newUserDAO()
	})
	return UserDAO
}

func (userDAO) AddNew(userName, userPwd string) (user model.User, err error) {
	if db.Db.Where(&model.User{UserName: userName}).Take(&user).RowsAffected != 0 {
		return user, fmt.Errorf("用户已存在:%s, %w", userName, common.UserAlreadyExistsError)
	}

	user = model.User{
		UserName: userName,
		UserPwd:  Md5sumPwd(userPwd),
	}
	result := db.Db.Create(&user)
	if result.RowsAffected == 0 {
		return user, fmt.Errorf("用户注册失败:%#v, %w", user, result.Error)
	}
	return
}

func (userDAO) CheckExistByOpts(opts ...Option) bool {
	var count int64
	db1 := db.Db.Model(&model.User{})
	for _, opt := range opts {
		db1 = opt(db1)
	}
	db1.Count(&count)
	if count != 0 {
		return true
	}
	return false
}

func (userDAO) GetByOpts(opts ...Option) *model.User {
	var user model.User
	db1 := db.Db.Model(&model.User{})
	for _, opt := range opts {
		db1 = opt(db1)
	}
	res := db1.Take(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &user
}

func Md5sumPwd(userPwd string) string {
	logger.Logger.Println("userPwd before", userPwd)
	x := fmt.Sprintf("%x", md5.Sum([]byte(userPwd)))
	logger.Logger.Println("userPwd after", x)
	return x
}

func WithUserName(userName string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("userName = ?", userName)
	}
}

func WithID(userID uint) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", userID)
	}
}
