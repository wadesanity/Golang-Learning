package common

import "errors"

var (
	UserAlreadyExistsError = errors.New("user already exists error")
	UserRegisterParamsError = errors.New("注册用户参数缺失")
	UserNotExistsError = errors.New("user not exists error")
	UserLoginParamsError = errors.New("用户登录密码错误")
	UserIDJwtError = errors.New("未能获取用户ID")
	TaskQueryParamsError = errors.New("任务查询参数缺失")
	TaskPostParamsError = errors.New("查询添加参数错误")
)

func a()  {
	UserAlreadyExistsError.Error()
}