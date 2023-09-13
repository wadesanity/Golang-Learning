package serve

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"net"
	"time"
	"user/conf"
	pb "user/grpc/pb/user"
	"user/interceptor"
	"user/pkg/util"
	"user/service"
)

func Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go registerEndPointToEtcd(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s", conf.UserListenAddr))
	if err != nil {
		panic(fmt.Errorf("failed to listen: %w", err))
	}
	util.Logger.Infof("start to serve at %v", lis.Addr())

	s := grpc.NewServer(grpc.ConnectionTimeout(1*time.Second),
		grpc.ChainUnaryInterceptor(interceptor.TraceIdFromMetadata2Context, interceptor.TimeoutControlInterceptor))
	pb.RegisterUserServiceServer(s, &service.UserService{})
	err = s.Serve(lis)
	if err != nil {
		panic(fmt.Errorf("failed to serve: %w", err))
	}

	util.Logger.Info("serve end!")
}

func registerEndPointToEtcd(ctx context.Context) {
	// 创建 etcd 客户端
	etcdClient, err := clientv3.NewFromURL(conf.EtcdAddr)
	if err != nil {
		panic(err)
	}
	util.Logger.WithFields(logrus.Fields{
		"etcdAddr": conf.EtcdAddr,
	}).Debug("NewFromURL info.")

	etcdManager, err := endpoints.NewManager(etcdClient, conf.UserServerAddr)
	if err != nil {
		panic(err)
	}
	util.Logger.WithFields(logrus.Fields{
		"endpoint_prefix": conf.UserServerAddr,
	}).Debug("NewManager info.")

	// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
	var ttl int64 = 10
	lease, _ := etcdClient.Grant(ctx, ttl)

	// 添加注册节点到 etcd 中，并且携带上租约 id
	epAddr := fmt.Sprintf("%s/%s", conf.UserServerAddr, conf.UserVisitAddr)
	err = etcdManager.AddEndpoint(ctx, epAddr, endpoints.Endpoint{Addr: conf.UserVisitAddr}, clientv3.WithLease(lease.ID))
	if err != nil {
		panic(err)
	}
	util.Logger.WithFields(logrus.Fields{
		"endpoint_addr": epAddr,
		"addr":          conf.UserVisitAddr,
	}).Debug("AddEndpoint info.")

	// 每隔 5 s进行一次延续租约的动作
	tickerC := time.Tick(5 * time.Second)
	for {
		select {
		case <-tickerC:
			// 续约操作
			resp, _ := etcdClient.KeepAliveOnce(ctx, lease.ID)
			util.Logger.Tracef("keep alive resp: %+v", resp)
		case <-ctx.Done():
			return
		}
	}
}
