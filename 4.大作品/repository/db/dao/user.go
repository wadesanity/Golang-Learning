package dao

import (
	"context"
	"gorm.io/gorm"
	"videoGo/repository/db"
	"videoGo/repository/db/model"
)

//var (
//UserDAO     *userDAO
//userDAOOnce sync.Once
//)

type UserDAOInstance interface {
	GetTotalByOpts(opts ...Option) (int64, error)
	GetByOpts(opts ...Option) (*model.User, error)
	AddNew(user *model.User) (*model.User, error)
	ChangeByModel(user *model.User) (*model.User, error)
}

type userDAO struct {
	*gorm.DB
}

func NewUserDAO(ctx context.Context) UserDAOInstance {
	if ctx == nil {
		ctx = context.Background()
	}
	return &userDAO{db.NewDBClient(ctx)}
}

//func GetUserDAO(ctx context.Context) UserDAOInstance {
//	userDAOOnce.Do(func() {
//		UserDAO = newUserDAO(ctx)
//	})
//	return UserDAO
//}

func (dao *userDAO) AddNew(user *model.User) (*model.User, error) {
	return user, dao.DB.Create(&user).Error
}

func (dao *userDAO) GetTotalByOpts(opts ...Option) (int64, error) {
	var count int64
	db1 := dao.DB.Model(&model.User{})
	for _, opt := range opts {
		db1 = opt(db1)
	}
	db1 = db1.Count(&count)
	return count, db1.Error
}

func (dao *userDAO) GetByOpts(opts ...Option) (*model.User, error) {
	var user model.User
	db1 := dao.DB.Model(&model.User{})
	for _, opt := range opts {
		db1 = opt(db1)
	}
	res := db1.Take(&user)
	return &user, res.Error
}

func (dao *userDAO) ChangeByModel(user *model.User) (*model.User, error) {
	return user, dao.DB.Updates(user).Error
}

type Option func(db *gorm.DB) *gorm.DB

func WithNameInUser(userName string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", userName)
	}
}

func WhereID(id uint) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

type fieldStringOrSlice interface {
	string | []string
}

func SelectField[F fieldStringOrSlice](f F) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(f)
	}
}
