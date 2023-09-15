package service

import (
	"context"
	pb "gateway/grpc/pb/user"
	"gateway/pkg/e"
	"gateway/pkg/util"
	"gateway/types/req"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	"sync"
)

var (
	UserServiceInstance *UserService
	userServiceOnce     sync.Once
)

//go:generate mockgen -source=user.go -destination=./mock/user_mock.go --package=mock
type UserServiceClient interface {
	Register(ctx context.Context, in *pb.RegisterReq, opts ...grpc.CallOption) (*pb.RegisterRes, error)
}

type UserService struct {
	client UserServiceClient
}

func newUserService(client UserServiceClient) *UserService {
	return &UserService{client: client}
}

func GetUserService(client UserServiceClient) *UserService {
	userServiceOnce.Do(func() {
		UserServiceInstance = newUserService(client)
	})
	return UserServiceInstance
}

func (s *UserService) Register(ctx context.Context, rq *req.UserRegisterReq) (err error) {
	in := &pb.RegisterReq{
		Name:   rq.Name,
		Pwd:    rq.Pwd,
		Avatar: rq.Avatar,
	}
	util.Logger.WithFields(logrus.Fields{
		"trace_id": ctx.Value(util.TraceIdKey),
		"in":       in,
	}).Trace("in info.")
	rs, err := s.client.Register(ctx, in)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"detail":   err,
		}).Errorln("grpc return err.")
		return ConvertGrpcError2http(err)
	}
	if !rs.Res {
		return e.NewApiError(http.StatusInternalServerError, "注册用户失败")
	}
	return
}

func (s *UserService) Login() {

}

func (s *UserService) ChangePwd() {

}

func (s *UserService) ShowInfo() {

}

func (s *UserService) ChangeAvatar() {

}

func (s *UserService) List() {

}
