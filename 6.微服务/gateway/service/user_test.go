package service

import (
	"context"
	"errors"
	"fmt"
	pb "gateway/grpc/pb/user"
	"gateway/pkg/e"
	"gateway/pkg/util"
	"gateway/types/req"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wadesanity/hystrix-go/hystrix"
	"google.golang.org/grpc"
	"net/http"
	"reflect"
	"testing"
)

type UserServiceMockClientImpl struct {
	client string
}

func NewUserMockClient(client string) *UserServiceMockClientImpl {
	return &UserServiceMockClientImpl{client: client}
}

func (*UserServiceMockClientImpl) Register(ctx context.Context, in *pb.RegisterReq, opts ...grpc.CallOption) (*pb.RegisterRes, error) {
	return &pb.RegisterRes{Res: true}, nil
}

//func (*UserServiceMockClientImpl) Login(ctx context.Context, in *pb.LoginReq, opts ...grpc.CallOption) (*pb.LoginRes, error) {
//	return &pb.LoginRes{}, nil
//}
//
//func (*UserServiceMockClientImpl) ChangePwd(ctx context.Context, in *pb.ChangePwdReq, opts ...grpc.CallOption) (*pb.ChangePwdRes, error) {
//	return &pb.ChangePwdRes{}, nil
//}
//
//func (*UserServiceMockClientImpl) ShowInfo(ctx context.Context, in *pb.ShowInfoReq, opts ...grpc.CallOption) (*pb.ShowInfoRes, error) {
//	return &pb.ShowInfoRes{}, nil
//}
//
//func (*UserServiceMockClientImpl) ChangeAvatar(ctx context.Context, in *pb.ChangeAvatarReq, opts ...grpc.CallOption) (*pb.ChangeAvatarRes, error) {
//	return &pb.ChangeAvatarRes{}, nil
//}
//
//func (*UserServiceMockClientImpl) List(ctx context.Context, in *pb.ListReq, opts ...grpc.CallOption) (*pb.ListRes, error) {
//	return &pb.ListRes{}, nil
//}

type UserServiceRegisterSuit struct {
	suite.Suite
	userClient  UserServiceClient
	userService *UserService
	patches     *gomonkey.Patches
}

func (u *UserServiceRegisterSuit) SetupTest() {
	util.Init()
	u.userClient = NewUserMockClient("test")
	u.userService = GetUserService(u.userClient)
}

func (u *UserServiceRegisterSuit) TearDownTest() {
	if u.patches != nil {
		u.patches.Reset()
	}
}

func (u *UserServiceRegisterSuit) TestTrue() {
	err := u.userService.Register(context.Background(), &req.UserRegisterReq{
		Name:   "1",
		Pwd:    "1",
		Avatar: "1",
	})
	assert.Nil(u.T(), err)
}

func (u *UserServiceRegisterSuit) TestFalse() {
	u.patches = gomonkey.ApplyMethod(reflect.TypeOf(u.userClient), "Register", func(_ *UserServiceMockClientImpl) (*pb.RegisterRes, error) {
		return &pb.RegisterRes{Res: false}, nil
	})
	err := u.userService.Register(context.Background(), &req.UserRegisterReq{
		Name:   "1",
		Pwd:    "1",
		Avatar: "1",
	})
	assert.NotNil(u.T(), err)
}

func (u *UserServiceRegisterSuit) TestErrorUnKnow() {
	u.patches = gomonkey.ApplyMethod(reflect.TypeOf(u.userClient), "Register", func(_ *UserServiceMockClientImpl) (*pb.RegisterRes, error) {
		return &pb.RegisterRes{Res: false}, errors.New("test err")
	})
	err := u.userService.Register(context.Background(), &req.UserRegisterReq{
		Name:   "1",
		Pwd:    "1",
		Avatar: "1",
	})
	assert.NotNil(u.T(), err)
	var apiError *e.ApiError
	assert.ErrorAs(u.T(), err, &apiError)
	assert.Equal(u.T(), http.StatusInternalServerError, apiError.HttpStatus)
}

func (u *UserServiceRegisterSuit) TestErrorHystrixTimeout() {
	u.patches = gomonkey.ApplyMethod(reflect.TypeOf(u.userClient), "Register", func(_ *UserServiceMockClientImpl) (*pb.RegisterRes, error) {
		return &pb.RegisterRes{Res: false}, fmt.Errorf("hystrix: %w", hystrix.ErrTimeout)
	})
	err := u.userService.Register(context.Background(), &req.UserRegisterReq{
		Name:   "1",
		Pwd:    "1",
		Avatar: "1",
	})
	assert.NotNil(u.T(), err)
	var apiError *e.ApiError
	assert.ErrorAs(u.T(), err, &apiError)
	assert.Equal(u.T(), http.StatusGatewayTimeout, apiError.HttpStatus)
}

func Test_UserService_Register(t *testing.T) {
	suite.Run(t, new(UserServiceRegisterSuit))
}
