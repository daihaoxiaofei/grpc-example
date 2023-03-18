package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"grpc-example/etcd-grpc/src"
	"grpc-example/pb"
	"grpc-example/pkg/closehelp"
	"grpc-example/server/hello"
	"net"
)

func main() {
	var gRPCHostAddr = "localhost:8002" // 服务器主机&监听端口

	// 将服务地址注册到etcd中
	err := src.Register(gRPCHostAddr)
	if err != nil {
		return
	}

	// 创建grpc句柄
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &hello.Server{})

	// 监听网络
	listener, err := net.Listen("tcp", gRPCHostAddr)
	if err != nil {
		fmt.Println("监听网络失败：", err)
		return
	}

	go func() {
		fmt.Println(`服务已开启`)
		err = s.Serve(listener)
		if err != nil {
			fmt.Println("监听异常：", err)
			return
		}
	}()

	closehelp.Register(func(ctx context.Context) {
		s.GracefulStop()
	})
	closehelp.SignalClose()

}
