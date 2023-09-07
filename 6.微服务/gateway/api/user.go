package api

import (
	"gateway/conf"
	pb "gateway/grpc/pb/user"
	"gateway/pkg/e"
	"gateway/pkg/util"
	"gateway/types/req"
	"gateway/types/res"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
)

func UserRegister(c *gin.Context) {
	var rq req.UserRegisterReq
	err := c.ShouldBind(&rq)
	if err != nil {
		util.Logger.Errorf("用户注册handler参数绑定错误:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.ReqParamsError.Error(),
		})
		return
	}
	util.Logger.Debugf("userRegisterReq:%v", rq)
	conn, err := grpc.Dial(conf.UserServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		util.Logger.Errorf("did not connect: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  e.GrpcDialError,
		})
		return
	}
	defer conn.Close()
	in := &pb.RegisterReq{
		Name:   rq.Name,
		Pwd:    rq.Pwd,
		Avatar: rq.Avatar,
	}
	client := pb.NewUserServiceClient(conn)
	util.Logger.Debugf("in:%v", in)
	rs, err := client.Register(c.Request.Context(), in)
	if err != nil {
		util.Logger.Errorf("用户注册handle错误:%v", err)
		s, ok := status.FromError(err)
		if !ok {
			util.Logger.Errorf("UserRegister FromError not ok")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  e.GrpcResError.Error(),
			})
			return
		}
		code, errString := ConvertCodeAndMessage2http(s)
		c.JSON(code, gin.H{
			"status": code,
			"error":  errString,
		})
		return
	}
	c.JSON(http.StatusOK, res.NewResOk("注册用户成功", http.StatusOK, rs.Res))
}

func UserLogin(c *gin.Context) {
	var rq req.UserLoginReq
	err := c.ShouldBind(&rq)
	if err != nil {
		util.Logger.Errorf("用户登录handler参数绑定错误:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.ReqParamsError.Error(),
		})
		return
	}
	util.Logger.Debugln("userLoginReq:", rq)
	conn, err := grpc.Dial(conf.UserServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		util.Logger.Errorf("did not connect: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  e.GrpcDialError,
		})
		return
	}
	defer conn.Close()
	var in = &pb.LoginReq{
		Name: rq.Name,
		Pwd:  rq.Pwd,
	}
	client := pb.NewUserServiceClient(conn)
	util.Logger.Debugf("in:%v", in)
	rs, err := client.Login(c.Request.Context(), in)
	if err != nil {
		util.Logger.Errorf("用户登录handle错误:%v", err)
		s, ok := status.FromError(err)
		if !ok {
			util.Logger.Errorf("UserLogin FromError not ok")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  e.GrpcResError.Error(),
			})
			return
		}
		code, errString := ConvertCodeAndMessage2http(s)
		c.JSON(code, gin.H{
			"status": code,
			"error":  errString,
		})
		return
	}
	c.JSON(http.StatusOK, res.NewResOk("用户登录成功", http.StatusOK, rs.Token))
}

func UserChangePwd(c *gin.Context) {
	var rq req.UserChangePwdReq
	err := c.ShouldBind(&rq)
	if err != nil {
		util.Logger.Errorf("用户修改密码handler参数绑定错误:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  e.ReqParamsError.Error(),
		})
		return
	}

	util.Logger.Debugf("UserChangePwdReq:%v", rq)
	conn, err := grpc.Dial(conf.UserServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		util.Logger.Errorf("did not connect: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  e.GrpcDialError,
		})
		return
	}
	defer conn.Close()
	var in = &pb.ChangePwdReq{
		Id:     uint32(c.GetUint("userID")),
		PwdOld: rq.PwdOld,
		PwdNew: rq.PwdNew,
	}
	client := pb.NewUserServiceClient(conn)
	util.Logger.Debugf("in:%v", in)
	rs, err := client.ChangePwd(c.Request.Context(), in)
	if err != nil {
		util.Logger.Errorf("用户改密码handle错误:%v", err)
		s, ok := status.FromError(err)
		if !ok {
			util.Logger.Errorf("UserLogin FromError not ok")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  e.GrpcResError.Error(),
			})
			return
		}
		code, errString := ConvertCodeAndMessage2http(s)
		c.JSON(code, gin.H{
			"status": code,
			"error":  errString,
		})
		return
	}
	c.JSON(http.StatusOK, res.NewResOk("修改密码成功", http.StatusOK, rs.Res))
}

func UserShowInfo(c *gin.Context) {
	conn, err := grpc.Dial(conf.UserServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		util.Logger.Errorf("did not connect: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  e.GrpcDialError,
		})
		return
	}
	defer conn.Close()
	id := c.GetUint("userID")
	if id == 0 {
		util.Logger.Errorf("userID not exists")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  e.AuthorizeError.Error(),
		})
		return
	}
	var in = &pb.ShowInfoReq{
		Id: uint32(id),
	}
	client := pb.NewUserServiceClient(conn)
	util.Logger.Debugf("in:%v", in)
	rs, err := client.ShowInfo(c.Request.Context(), in)
	if err != nil {
		util.Logger.Errorf("用户详情handle错误:%v", err)
		s, ok := status.FromError(err)
		if !ok {
			util.Logger.Errorf("UserLogin FromError not ok")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  e.GrpcResError.Error(),
			})
			return
		}
		code, errString := ConvertCodeAndMessage2http(s)
		c.JSON(code, gin.H{
			"status": code,
			"error":  errString,
		})
		return
	}
	util.Logger.Debugf("rs:%v", rs)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "用户信息查询成功", "data": gin.H{
		"id":         rs.Id,
		"name":       rs.Name,
		"avatar":     rs.Avatar,
		"created_at": rs.CreatedAt.Seconds,
		"updated_at": rs.UpdatedAt.Seconds,
		"role":       rs.Role.String(),
		"status":     rs.Status.String(),
	}})
}

func UserChangeAvatar(c *gin.Context) {
	var rq req.UserChangeAvatarReq
	err := c.ShouldBindWith(&rq, binding.Form)
	if err != nil {
		util.Logger.Errorf("用户修改头像handler参数绑定错误:%v", err)
		c.JSON(http.StatusBadRequest, e.ReqParamsError.Error())
		return
	}
	conn, err := grpc.Dial(conf.UserServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		util.Logger.Errorf("did not connect: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  e.GrpcDialError,
		})
		return
	}
	defer conn.Close()
	var in = &pb.ChangeAvatarReq{
		Id:     uint32(c.GetUint("userID")),
		Avatar: rq.Avatar,
	}
	client := pb.NewUserServiceClient(conn)
	util.Logger.Debugf("in:%v", in)
	rs, err := client.ChangeAvatar(c.Request.Context(), in)
	if err != nil {
		util.Logger.Errorf("用户改头像handle错误:%v", err)
		s, ok := status.FromError(err)
		if !ok {
			util.Logger.Errorf("UserLogin FromError not ok")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  e.GrpcResError.Error(),
			})
			return
		}
		code, errString := ConvertCodeAndMessage2http(s)
		c.JSON(code, gin.H{
			"status": code,
			"error":  errString,
		})
		return
	}
	c.JSON(http.StatusOK, res.NewResOk("用户头像信息修改成功", http.StatusOK, rs.Res))
}

func UserList(c *gin.Context) {
	var rq req.UserListReq
	err := c.ShouldBind(&rq)
	if err != nil {
		util.Logger.Errorf("UserList shouldBind err:%v", err)
		c.JSON(http.StatusBadRequest, res.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	conn, err := grpc.Dial(conf.UserServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		util.Logger.Errorf("did not connect: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  e.GrpcDialError,
		})
		return
	}
	defer conn.Close()
	var in = &pb.ListReq{
		Limit: uint32(rq.Limit),
	}
	if rq.Name != nil {
		in.Name = *rq.Name
	}
	if rq.CreateStart != nil {
		in.CreateStart = timestamppb.New(time.Unix(*rq.CreateStart, 0))
	}
	if rq.CreateEnd != nil {
		in.CreateEnd = timestamppb.New(time.Unix(*rq.CreateEnd, 0))

	}
	if rq.Order != nil {
		in.Order = *rq.Order
	}
	if rq.Offset != nil {
		in.Offset = uint32(*rq.Offset)
	}
	client := pb.NewUserServiceClient(conn)
	util.Logger.Debugf("in:%v", in)
	rs, err := client.List(c.Request.Context(), in)
	if err != nil {
		util.Logger.Errorf("用户列表handle错误:%v", err)
		s, ok := status.FromError(err)
		if !ok {
			util.Logger.Errorf("UserLogin FromError not ok")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  e.GrpcResError.Error(),
			})
			return
		}
		code, errString := ConvertCodeAndMessage2http(s)
		c.JSON(code, gin.H{
			"status": code,
			"error":  errString,
		})
		return
	}
	c.JSON(http.StatusOK, res.NewResList(rs.List, int64(rs.Total), "用户列表查询成功"))
}
