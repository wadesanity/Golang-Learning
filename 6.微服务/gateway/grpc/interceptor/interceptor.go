package interceptor

import (
	"context"
	"gateway/conf"
	"gateway/pkg/util"
	"github.com/wadesanity/hystrix-go/hystrix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func init() {
	hystrix.ConfigureCommand("etcd:///micro-case/user", hystrix.CommandConfig{
		Timeout:                conf.UserServerTimeoutSecond * 1000,
		MaxConcurrentRequests:  conf.UserServerMaxConcurrentRequests,
		RequestVolumeThreshold: conf.UserServerRequestVolumeThreshold,
		SleepWindow:            conf.UserServerSleepWindowSecond * 1000,
		ErrorPercentThreshold:  conf.UserServerErrorPercentThreshold,
	})
}

func TimeoutAndTraceIdClientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Pre-processor phase
	deadLine, ok := ctx.Deadline()
	if ok {
		util.Logger.Debugf("ctx already had timeout befor gin invoke to metadata. deadline:%v", deadLine)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(conf.RpcAllTimeout)*time.Second)
	defer cancel()
	deadLine, _ = ctx.Deadline()
	util.Logger.Debugf("ctx gen timeout from gin invoke to metadata. deadline:%v", deadLine)
	// Invoking the remote method
	traceId := ctx.Value(util.TraceIdKey).(string)
	util.Logger.Debugf("trace_id new and deadline from gin to metadata. trace_id:%v", traceId)
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"trace_id": traceId, "deadline": deadLine.String()}))
	err := invoker(ctx, method, req, reply, cc, opts...)

	return err
}

func UserHystrixClientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	util.Logger.Debugf("cc target:%v", cc.Target())
	return hystrix.DoC(ctx, cc.Target(), func(ctx context.Context) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}, func(ctx context.Context, err error) error {
		util.Logger.Errorf("fallback err:%v", err)
		return err
	})
}
