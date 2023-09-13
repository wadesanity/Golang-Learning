package conn

import (
	"context"
	"fmt"
	"gateway/conf"
	"gateway/grpc/interceptor"
	pb "gateway/grpc/pb/user"
	"gateway/pkg/util"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var UserClient pb.UserServiceClient

func Init() *grpc.ClientConn {
	cli, err := clientv3.NewFromURL(conf.EtcdAddr)
	if err != nil {
		panic(err)
	}
	etcdResolver, err := resolver.NewBuilder(cli)
	if err != nil {
		panic(err)
	}
	etcdTarget := fmt.Sprintf("etcd:///%s", conf.UserServerAddr)
	util.Logger.Infof("start etcdTarget:%v", etcdTarget)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	conn, err := grpc.DialContext(ctx,
		etcdTarget,
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.UserHystrixClientInterceptor,
			interceptor.TimeoutAndTraceIdClientInterceptor),
	)
	if err != nil {
		panic(err)
	}
	UserClient = pb.NewUserServiceClient(conn)
	return conn
}
