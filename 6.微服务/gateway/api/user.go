package api

import (
	"gateway/grpc/conn"
	pb "gateway/grpc/pb/user"
	"gateway/pkg/e"
	"gateway/pkg/util"
	"gateway/types/req"
	"gateway/types/res"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
)

// UserRegister godoc
//
//	@Summary		Register an account
//	@Description	Register by name and pwd
//	@Tags			users
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name	formData		string	true	"username"
//	@Param			pwd		formData		string	true	"pwd"
//	@Param			avatar	formData		string	false	"avatar"
//	@Success		200		{object}	res.Response
//	@Failure		400		{object}	res.Response
//	@Failure		404		{object}	res.Response
//	@Failure		500		{object}	res.Response
//	@Router			/user_register [post]
func UserRegister(c *gin.Context) {
	var rq req.UserRegisterReq
	err := c.ShouldBind(&rq)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": c.Request.Context().Value(util.TraceIdKey),
			"detail":   err,
		}).Errorln("rq bind return err")
		c.JSON(http.StatusBadRequest, res.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	util.Logger.Debugf("userRegisterReq:%v", rq)
	in := &pb.RegisterReq{
		Name:   rq.Name,
		Pwd:    rq.Pwd,
		Avatar: rq.Avatar,
	}
	util.Logger.Debugf("in:%v", in)
	rs, err := conn.UserClient.Register(c.Request.Context(), in)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"detail": err,
		}).Errorln("grpc return err")
		ConvertGrpcError2http(err, c)
		return
	}
	c.JSON(http.StatusOK, res.NewResOk("注册用户成功", http.StatusOK, rs.Res))
}

func UserLogin(c *gin.Context) {
	var rq req.UserLoginReq
	err := c.ShouldBind(&rq)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": c.Request.Context().Value(util.TraceIdKey),
			"detail":   err,
		}).Errorln("rq bind return err")
		c.JSON(http.StatusBadRequest, res.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	util.Logger.Debugln("userLoginReq:", rq)
	var in = &pb.LoginReq{
		Name: rq.Name,
		Pwd:  rq.Pwd,
	}
	util.Logger.Debugf("in:%v", in)
	rs, err := conn.UserClient.Login(c.Request.Context(), in)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"detail": err,
		}).Errorln("grpc return err")
		ConvertGrpcError2http(err, c)
		return
	}
	c.JSON(http.StatusOK, res.NewResOk("用户登录成功", http.StatusOK, rs.Token))
}

func UserChangePwd(c *gin.Context) {
	var rq req.UserChangePwdReq
	err := c.ShouldBind(&rq)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": c.Request.Context().Value(util.TraceIdKey),
			"detail":   err,
		}).Errorln("rq bind return err")
		c.JSON(http.StatusBadRequest, res.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}

	util.Logger.Debugf("UserChangePwdReq:%v", rq)

	var in = &pb.ChangePwdReq{
		Id:     uint32(c.GetUint("userID")),
		PwdOld: rq.PwdOld,
		PwdNew: rq.PwdNew,
	}
	util.Logger.Debugf("in:%v", in)
	rs, err := conn.UserClient.ChangePwd(c.Request.Context(), in)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"detail": err,
		}).Errorln("grpc return err")
		ConvertGrpcError2http(err, c)
		return
	}
	c.JSON(http.StatusOK, res.NewResOk("修改密码成功", http.StatusOK, rs.Res))
}

func UserShowInfo(c *gin.Context) {
	id := c.GetUint("userID")
	if id == 0 {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": c.Request.Context().Value(util.TraceIdKey),
		}).Errorln("rq userID not exists")
		c.JSON(http.StatusBadRequest, res.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	var in = &pb.ShowInfoReq{
		Id: uint32(id),
	}
	util.Logger.Debugf("in:%v", in)
	rs, err := conn.UserClient.ShowInfo(c.Request.Context(), in)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"detail": err,
		}).Errorln("grpc return err")
		ConvertGrpcError2http(err, c)
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
		util.Logger.WithFields(logrus.Fields{
			"trace_id": c.Request.Context().Value(util.TraceIdKey),
			"detail":   err,
		}).Errorln("rq bind return err")
		c.JSON(http.StatusBadRequest, res.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
	var in = &pb.ChangeAvatarReq{
		Id:     uint32(c.GetUint("userID")),
		Avatar: rq.Avatar,
	}
	util.Logger.Debugf("in:%v", in)
	rs, err := conn.UserClient.ChangeAvatar(c.Request.Context(), in)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"detail": err,
		}).Errorln("grpc return err")
		ConvertGrpcError2http(err, c)
		return
	}
	c.JSON(http.StatusOK, res.NewResOk("用户头像信息修改成功", http.StatusOK, rs.Res))
}

func UserList(c *gin.Context) {
	var rq req.UserListReq
	err := c.ShouldBind(&rq)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": c.Request.Context().Value(util.TraceIdKey),
			"detail":   err,
		}).Errorln("rq bind return err")
		c.JSON(http.StatusBadRequest, res.NewResError(http.StatusBadRequest, e.ReqParamsError))
		return
	}
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
	util.Logger.WithFields(logrus.Fields{
		"trace_id": c.Request.Context().Value(util.TraceIdKey),
		"in":       in,
		"rq":       rq,
	}).Debugln("rq and in info")
	//ctx, cancel := context.WithCancel(c.Request.Context())
	//go func() {
	//	time.Sleep(200 * time.Millisecond)
	//	cancel()
	//}()
	rs, err := conn.UserClient.List(c.Request.Context(), in)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": c.Request.Context().Value(util.TraceIdKey),
			"detail":   err,
		}).Error("grpc return err")
		ConvertGrpcError2http(err, c)
		return
	}
	c.JSON(http.StatusOK, res.NewResList(rs.List, int64(rs.Total), "用户列表查询成功"))
}
