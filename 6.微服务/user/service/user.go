package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
	"user/conf"
	pb "user/grpc/pb/user"
	"user/pkg/e"
	"user/pkg/util"
	"user/repository/db"
	"user/repository/db/dao"
	"user/repository/db/model"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (*UserService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	userDao := dao.NewUserDAO(ctx)
	_, err := userDao.GetByOpts(dao.WithNameInUser(req.Name))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.WithFields(logrus.Fields{
				"trace_id": ctx.Value(util.TraceIdKey),
				"req":      req,
				"detail":   err,
			}).Errorln("userDao.GetByOpts error.")
			return nil, status.Error(codes.Internal, e.DbQueryError.Error())
		}
	} else {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"req":      req,
			"detail":   err,
		}).Debugln("userDao.GetTotalByOpts AlreadyExists.")
		return nil, status.Error(codes.AlreadyExists, e.DbQueryAlreadyFound.Error())
	}
	var user = &model.User{}
	user.Name = req.Name
	user.Avatar = req.Avatar
	user.Md5sumPwd(req.Pwd)
	user, err = userDao.AddNew(user)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"req":      user,
			"detail":   err,
		}).Errorln("userDao.AddNew error.")
		return nil, status.Error(codes.Internal, e.DbCreateError.Error())
	}
	return &pb.RegisterRes{Res: true}, nil
}

func (*UserService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	userDao := dao.NewUserDAO(ctx)
	user, err := userDao.GetByOpts(dao.WithNameInUser(req.Name))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.WithFields(logrus.Fields{
				"trace_id": ctx.Value(util.TraceIdKey),
				"req":      user,
				"detail":   err,
			}).Warningln("userDao.GetByOpts not found.")
			util.Logger.Errorf("用户登录方法->用户查询方法->not found,请求形参:%v", req.Name)
			return nil, status.Error(codes.NotFound, e.DbQueryNotFound.Error())
		}
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"req":      user,
			"detail":   err,
		}).Errorln("userDao.GetByOpts error.")
		return nil, status.Error(codes.Internal, e.DbQueryError.Error())
	}
	if !user.CheckPwd(req.Pwd) {
		return nil, status.Error(codes.InvalidArgument, e.ReqParamsError.Error())
	}

	var isAuditor bool
	if user.Role == 1 {
		isAuditor = true
	}

	if user.Status == 1 {
		return nil, status.Error(codes.Unauthenticated, e.AuthorizeError.Error())
	}

	return &pb.LoginRes{Token: util.NewJwt(user.ID, isAuditor)}, nil

}

func (*UserService) ChangePwd(ctx context.Context, req *pb.ChangePwdReq) (*pb.ChangePwdRes, error) {
	userDao := dao.NewUserDAO(ctx)
	user, err := userDao.GetByOpts(dao.WhereID(req.Id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.WithFields(logrus.Fields{
				"trace_id": ctx.Value(util.TraceIdKey),
				"user_id":  req.Id,
				"detail":   err,
			}).Warningln("userDao.GetByOpts not found.")
			return nil, status.Error(codes.AlreadyExists, e.DbQueryNotFound.Error())
		}
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"user_id":  req.Id,
			"detail":   err,
		}).Errorln("userDao.GetByOpts error.")
		return nil, status.Error(codes.Internal, e.DbQueryError.Error())
	}
	if !user.CheckPwd(req.PwdOld) {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"user_id":  req.Id,
		}).Warningln("user CheckPwd error.")
		return nil, status.Error(codes.InvalidArgument, e.ReqParamsError.Error())
	}
	var userNew = &model.User{
		ID: user.ID,
	}
	userNew.Md5sumPwd(req.PwdNew)
	userNew, err = userDao.ChangeByModel(userNew)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"user":     userNew,
			"detail":   err,
		}).Errorln("user CheckPwd error.")
		return nil, status.Error(codes.Internal, e.DbUpdateError.Error())
	}
	return &pb.ChangePwdRes{Res: true}, nil
}

func (*UserService) ShowInfo(ctx context.Context, req *pb.ShowInfoReq) (*pb.ShowInfoRes, error) {
	userDao := dao.NewUserDAO(ctx)
	user, err := userDao.GetByOpts(dao.WhereID(req.Id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.WithFields(logrus.Fields{
				"trace_id": ctx.Value(util.TraceIdKey),
				"user_id":  req.Id,
			}).Warningln("user GetByOpts not found.")
			return nil, status.Error(codes.NotFound, e.DbQueryNotFound.Error())
		}
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"user_id":  req.Id,
			"detail":   err,
		}).Errorln("user GetByOpts err.")
		return nil, status.Error(codes.Internal, e.DbQueryError.Error())
	}
	return &pb.ShowInfoRes{
		Id:        user.ID,
		Name:      user.Name,
		Avatar:    user.Avatar,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Role:      pb.ShowInfoRes_Role(user.Role),
		Status:    pb.ShowInfoRes_Status(user.Status),
	}, nil
}

func (*UserService) ChangeAvatar(ctx context.Context, req *pb.ChangeAvatarReq) (*pb.ChangeAvatarRes, error) {
	userDao := dao.NewUserDAO(ctx)
	_, err := userDao.GetByOpts(dao.WhereID(req.Id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.WithFields(logrus.Fields{
				"trace_id": ctx.Value(util.TraceIdKey),
				"user_id":  req.Id,
			}).Warningln("user GetByOpts not found.")
			return nil, status.Error(codes.AlreadyExists, e.DbQueryNotFound.Error())
		}
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"user_id":  req.Id,
			"detail":   err,
		}).Errorln("user GetByOpts err.")
		return nil, status.Error(codes.Internal, e.DbQueryError.Error())
	}

	var userNew = &model.User{
		ID:     req.Id,
		Avatar: req.Avatar,
	}
	userNew, err = userDao.ChangeByModel(userNew)
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"user":     userNew,
			"detail":   err,
		}).Errorln("user ChangeByModel err.")
		return nil, status.Error(codes.Internal, e.DbUpdateError.Error())
	}
	return &pb.ChangeAvatarRes{Res: true}, nil
}

func (*UserService) List(ctx context.Context, req *pb.ListReq) (*pb.ListRes, error) {
	time.Sleep(time.Duration(conf.UserTimeoutMillisecond) * time.Millisecond)
	db1 := db.NewDBClient(ctx).Model(&model.User{})
	if req.Name != "" {
		db1 = db1.Where("name like ?", fmt.Sprintf("%%%s%%", req.Name))
	}
	if req.CreateStart != nil {
		db1 = db1.Where("created_at >= ?", req.CreateStart)
	}
	if req.CreateEnd != nil {
		db1 = db1.Where("created_at <= ?", req.CreateEnd, 0)
	}
	var total int64
	var res []*model.User
	err := db1.Count(&total).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.WithFields(logrus.Fields{
				"trace_id": ctx.Value(util.TraceIdKey),
				"req":      req,
			}).Infoln("user count not found.")
			return &pb.ListRes{
				Total: 0,
				List:  nil,
			}, nil
		}
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"req":      req,
			"detail":   err,
		}).Errorln("user count error.")
		return &pb.ListRes{
			Total: 0,
			List:  nil,
		}, status.Error(codes.Internal, e.DbQueryError.Error())
	}
	if total == 0 {
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"req":      req,
		}).Infoln("user count not found.")
		return &pb.ListRes{
			Total: 0,
			List:  nil,
		}, nil
	}

	order := "created_at ASC"
	if req.Order != "" {
		order = req.Order
	}
	err = db1.Order(order).Offset(int(req.Offset)).Limit(int(req.Limit)).Find(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.Logger.WithFields(logrus.Fields{
				"trace_id": ctx.Value(util.TraceIdKey),
				"req":      req,
			}).Infoln("user list not found.")
			return &pb.ListRes{
				Total: 0,
				List:  nil,
			}, nil
		}
		util.Logger.WithFields(logrus.Fields{
			"trace_id": ctx.Value(util.TraceIdKey),
			"req":      req,
			"detail":   err,
		}).Errorln("user list error.")
		return &pb.ListRes{
			Total: uint64(total),
			List:  nil,
		}, status.Error(codes.Internal, e.DbQueryError.Error())
	}
	var r []*pb.ShowInfoRes
	for _, u := range res {
		d := &pb.ShowInfoRes{
			Id:        u.ID,
			Name:      u.Name,
			Avatar:    u.Avatar,
			CreatedAt: timestamppb.New(u.CreatedAt),
			UpdatedAt: timestamppb.New(u.UpdatedAt),
			Role:      pb.ShowInfoRes_Role(u.Role),
			Status:    pb.ShowInfoRes_Status(u.Status),
		}
		r = append(r, d)
	}

	return &pb.ListRes{
		Total: uint64(total),
		List:  r,
	}, nil
}
