package api

import (
	"context"
	"fmt"
	"gateway/grpc/conn"
	pb "gateway/grpc/pb/user"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	http_test "github.com/stretchr/testify/http"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"net/http"
	"testing"
)

type UserServiceMockClient interface {
	Register(ctx context.Context, in *pb.RegisterReq, opts ...grpc.CallOption) (*pb.RegisterRes, error)
	Login(ctx context.Context, in *pb.LoginReq, opts ...grpc.CallOption) (*pb.LoginRes, error)
	ChangePwd(ctx context.Context, in *pb.ChangePwdReq, opts ...grpc.CallOption) (*pb.ChangePwdRes, error)
	ShowInfo(ctx context.Context, in *pb.ShowInfoReq, opts ...grpc.CallOption) (*pb.ShowInfoRes, error)
	ChangeAvatar(ctx context.Context, in *pb.ChangeAvatarReq, opts ...grpc.CallOption) (*pb.ChangeAvatarRes, error)
	List(ctx context.Context, in *pb.ListReq, opts ...grpc.CallOption) (*pb.ListRes, error)
}

type UserServiceMockClientImpl struct{}

func (*UserServiceMockClientImpl) Register(ctx context.Context, in *pb.RegisterReq, opts ...grpc.CallOption) (*pb.RegisterRes, error) {
	return &pb.RegisterRes{}, nil
}

func (*UserServiceMockClientImpl) Login(ctx context.Context, in *pb.LoginReq, opts ...grpc.CallOption) (*pb.LoginRes, error) {
	return &pb.LoginRes{}, nil
}

func (*UserServiceMockClientImpl) ChangePwd(ctx context.Context, in *pb.ChangePwdReq, opts ...grpc.CallOption) (*pb.ChangePwdRes, error) {
	return &pb.ChangePwdRes{}, nil
}

func (*UserServiceMockClientImpl) ShowInfo(ctx context.Context, in *pb.ShowInfoReq, opts ...grpc.CallOption) (*pb.ShowInfoRes, error) {
	return &pb.ShowInfoRes{}, nil
}

func (*UserServiceMockClientImpl) ChangeAvatar(ctx context.Context, in *pb.ChangeAvatarReq, opts ...grpc.CallOption) (*pb.ChangeAvatarRes, error) {
	return &pb.ChangeAvatarRes{}, nil
}

func (*UserServiceMockClientImpl) List(ctx context.Context, in *pb.ListReq, opts ...grpc.CallOption) (*pb.ListRes, error) {
	return &pb.ListRes{}, nil
}

type UserApiTestSuite struct {
	suite.Suite
	patches *gomonkey.Patches
}

// 在这个测试组开始之前运行一次
func (suite *UserApiTestSuite) SetupSuite() {
	fmt.Println("start SetupSuite")
	patches := gomonkey.ApplyGlobalVar(conn.UserClient.Register, new(UserServiceMockClientImpl).Register)
	suite.patches = patches
}

// 在这个测试组结束后运行一次
func (suite *UserApiTestSuite) TearDownSuite() {
	fmt.Println("end SetupSuite")
	suite.patches.Reset()
}

// 所有这些的 Test 开头的方法，都会在 go test 的时候运行
func (suite *UserApiTestSuite) TestUserRegister() {
	c, _ := gin.CreateTestContext(&http_test.TestResponseWriter{
		StatusCode: 0,
		Output:     "",
	})
	UserRegister(c)
	assert.Equal(suite.T(), http.StatusOK, c.Writer.Status())
}

func (suite *UserApiTestSuite) TestUserLogin() {
	c, _ := gin.CreateTestContext(&http_test.TestResponseWriter{
		StatusCode: 0,
		Output:     "",
	})
	UserLogin(c)
	assert.Equal(suite.T(), http.StatusOK, c.Writer.Status())
}

func (suite *UserApiTestSuite) TestUserChangePwd() {
	c, _ := gin.CreateTestContext(&http_test.TestResponseWriter{
		StatusCode: 0,
		Output:     "",
	})
	UserChangePwd(c)
	assert.Equal(suite.T(), http.StatusOK, c.Writer.Status())
}

func (suite *UserApiTestSuite) TestUserShowInfo() {
	c, _ := gin.CreateTestContext(&http_test.TestResponseWriter{
		StatusCode: 0,
		Output:     "",
	})
	UserShowInfo(c)
	assert.Equal(suite.T(), http.StatusOK, c.Writer.Status())
}

func (suite *UserApiTestSuite) TestUserChangeAvatar() {
	c, _ := gin.CreateTestContext(&http_test.TestResponseWriter{
		StatusCode: 0,
		Output:     "",
	})
	UserChangeAvatar(c)
	assert.Equal(suite.T(), http.StatusOK, c.Writer.Status())
}

func (suite *UserApiTestSuite) TestUserList() {
	c, _ := gin.CreateTestContext(&http_test.TestResponseWriter{
		StatusCode: 0,
		Output:     "",
	})
	UserList(c)
	assert.Equal(suite.T(), http.StatusOK, c.Writer.Status())
}

func Test_User_Api(t *testing.T) {
	suite.Run(t, &UserApiTestSuite{})
}
