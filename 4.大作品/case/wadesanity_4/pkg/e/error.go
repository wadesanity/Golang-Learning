package e

import (
	"errors"
)

var (
	AuthorizeError = errors.New("认证失败")
	ReqParamsError = errors.New("参数异常")

	DbQueryError        = errors.New("数据查询错误")
	DbQueryNotFound     = errors.New("数据未找到")
	DbQueryAlreadyFound = errors.New("数据已存在")
	DbCreateError       = errors.New("数据插入错误")
	DbUpdateError       = errors.New("数据更新错误")
	DbDeleteError       = errors.New("数据删除错误")

	UnknowError = errors.New("未知错误")

	CacheQueryError        = errors.New("数据查询错误")
	CacheQueryNotFound     = errors.New("数据未找到")
	CacheQueryAlreadyFound = errors.New("数据已存在")
	CacheCreateError       = errors.New("数据插入错误")
	CacheUpdateError       = errors.New("数据更新错误")
	CacheDeleteError       = errors.New("数据删除错误")

	RepeatActionError = errors.New("重复操作")
)

type ApiError struct {
	S          string
	HttpStatus int
}

func (e *ApiError) Error() string {
	return e.S
}

func NewApiError(httpStatus int, text string) error {
	return &ApiError{text, httpStatus}
}
