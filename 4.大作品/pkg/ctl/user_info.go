package ctl

import (
	"context"
	"errors"
)

type userIDKeyType struct{}

var userIDKey = userIDKeyType{}

type UserInfo struct {
	Id uint `json:"id"`
}

func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	user, ok := fromContext(ctx)
	if !ok {
		return nil, errors.New("获取用户信息错误")
	}
	return user, nil
}

func NewContextWithUserIDKey(ctx context.Context, u *UserInfo) context.Context {
	return context.WithValue(ctx, userIDKey, u)
}

func fromContext(ctx context.Context) (*UserInfo, bool) {
	u, ok := ctx.Value(userIDKey).(*UserInfo)
	return u, ok
}
