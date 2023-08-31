package service

import (
	"fmt"
	"sync"
	"todolistGo/common"
	"todolistGo/dao"
)

var (
	UserServiceInstance *userService
	userServiceOnce     sync.Once
)

type userService struct {
}

func (*userService) Register(userName, userPwd string) (*common.ResponseOk, error) {
	i := dao.GetUserDAO()
	user, err := i.AddNew(userName, userPwd)
	if err != nil {
		return nil, err
	}
	r := make(map[string]interface{})
	r["id"] = user.ID
	r["userName"] = user.UserName
	r["userPwd"] = "xxx"
	r["createdTime"] = user.CreatedAt
	r["updatedTime"] = user.UpdatedAt
	item := make([]map[string]interface{}, 0)
	item = append(item, r)
	data := common.Data{
		Item:  item,
		Total: 1,
	}

	return &common.ResponseOk{
		Status: 200,
		Data:   data,
		Msg:    "添加用户成功",
	}, nil

}

func (*userService) Login(userName, userPwd string) (*common.ResponseOk, error) {
	userDao := dao.GetUserDAO()
	user := userDao.GetByOpts(dao.WithUserName(userName))
	if user == nil {
		return nil, fmt.Errorf("用户不存在:%s, %w", userName, common.UserNotExistsError)
	}
	if user.UserPwd != dao.Md5sumPwd(userPwd) {
		return nil, common.UserLoginParamsError
	}

	r := make(map[string]interface{})
	r["id"] = user.ID
	r["userName"] = user.UserName
	r["userPwd"] = "xxx"
	r["createdTime"] = user.CreatedAt
	r["updatedTime"] = user.UpdatedAt
	r["jwt"] = common.NewJwt(user.ID)
	item := make([]map[string]interface{}, 0)
	item = append(item, r)
	data := common.Data{
		Item:  item,
		Total: 1,
	}

	return &common.ResponseOk{
		Status: 200,
		Data:   data,
		Msg:    "用户登录成功",
	}, nil

}

func newUserService() *userService {
	return &userService{}
}

type UserService interface {
	Register(userName, userPwd string) (*common.ResponseOk, error)
	Login(userName, userPwd string) (*common.ResponseOk, error)
}

func GetUserServiceInstance() UserService {
	userServiceOnce.Do(func() {
		UserServiceInstance = newUserService()
	})
	return UserServiceInstance
}
