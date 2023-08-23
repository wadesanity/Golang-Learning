package e

import (
	"errors"
)

var (
	AuthorizeError = errors.New("参数异常")
	ReqParamsError = errors.New("参数异常")

	UserIDJwtError = errors.New("未能获取用户ID")
	TaskQueryParamsError = errors.New("任务查询参数缺失")
	TaskPostParamsError = errors.New("查询添加参数错误")

	DbQueryError = errors.New("数据查询错误")
	DbQueryNotFound = errors.New("数据未找到")
	DbQueryAlreadyFound = errors.New("数据已存在")
	DbCreateError = errors.New("数据插入错误")
	DbUpdateError = errors.New("数据更新错误")
	DbDeleteError = errors.New("数据删除错误")

	UnknowError = errors.New("未知错误")
)

type ApiError struct {
	S string
	HttpStatus int
}

func (e *ApiError) Error() string {
	return e.S
}

func NewApiError(httpStatus int, text string) error{
	return &ApiError{text,httpStatus}
}