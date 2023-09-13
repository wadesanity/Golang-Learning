package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
	"user/pkg/util"
)

func TraceIdFromMetadata2Context(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s := fmt.Sprintf("userService TraceIdFromReq2Context metadata get error, context:%v, req:%v", ctx, req)
		return nil, status.Error(codes.Unauthenticated, s)
	}
	v := md.Get("trace_id")
	if len(v) < 1 {
		s := fmt.Sprintf("userService TraceIdFromReq2Context traceId get error, context:%v, req:%v", ctx, req)
		return nil, status.Error(codes.Unauthenticated, s)
	}
	ctx = context.WithValue(ctx, util.TraceIdKey, v[0])
	return handler(ctx, req)
}

func TimeoutControlInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	traceId, _ := ctx.Value(util.TraceIdKey).(string)
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	deadline, _ := ctx.Deadline()
	defer cancel()
	s := fmt.Sprintf("traceId:%v, userService TimeoutControlInterceptor set context deadline:%v", traceId, deadline)
	util.Logger.Debug(s)
	return handler(ctx, req)
}
