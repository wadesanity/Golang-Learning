package service

import (
	"context"
	pb "gateway/grpc/pb/user"
	"gateway/pkg/util"
	"gateway/service/mock"
	"gateway/types/req"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_UserService_Register2(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockUserServiceClient := mock.NewMockUserServiceClient(ctl)
	mockUserServiceClient.EXPECT().Register(context.Background(), &pb.RegisterReq{
		Name:   "1",
		Pwd:    "",
		Avatar: "",
	}).DoAndReturn(func(ctx context.Context, req *pb.RegisterReq, opts ...interface{}) (*pb.RegisterRes, error) {
		return &pb.RegisterRes{Res: true}, nil
	})
	userS := GetUserService(mockUserServiceClient)
	util.Init()
	err := userS.Register(context.Background(), &req.UserRegisterReq{
		Name:   "1",
		Pwd:    "",
		Avatar: "",
	})
	assert.Nil(t, err)
}
