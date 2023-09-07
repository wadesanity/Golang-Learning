package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"user/conf"
	pb "user/grpc/pb/user"
	"user/pkg/util"
	"user/repository/db"
	"user/service"
)

func main() {
	db.Init()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", conf.UserServerAddr))
	if err != nil {
		util.Logger.Fatalf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &service.UserService{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		util.Logger.Fatalf("failed to serve: %v", err)
	}
}
