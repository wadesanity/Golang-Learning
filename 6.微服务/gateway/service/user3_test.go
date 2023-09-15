package service

import (
	"context"
	pb "gateway/grpc/pb/user"
	"gateway/pkg/util"
	"gateway/types/req"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"testing"
)

type UserServiceMock3ClientImpl struct {
	mock.Mock
	client string
}

func (u *UserServiceMock3ClientImpl) Register(ctx context.Context, in *pb.RegisterReq, opts ...grpc.CallOption) (*pb.RegisterRes, error) {
	_ = u.Called(ctx, in, opts)
	return &pb.RegisterRes{Res: true}, nil
}

func Test_UserService_Register3(t *testing.T) {
	userServiceMock3ClientImpl := &UserServiceMock3ClientImpl{}
	userServiceMock3ClientImpl.On("Register", context.Background(), &pb.RegisterReq{
		Name:   "",
		Pwd:    "",
		Avatar: "",
	}, *new([]grpc.CallOption)).Return(&pb.RegisterRes{Res: true}, nil)
	userS := GetUserService(userServiceMock3ClientImpl)
	util.Init()
	err := userS.Register(context.Background(), &req.UserRegisterReq{
		Name:   "1",
		Pwd:    "",
		Avatar: "",
	})
	assert.Nil(t, err)
}
