package main

import (
	"context"
	"google.golang.org/grpc"
	"grpc-example/pb"
	"grpc-example/pkg/closehelp"
	"grpc-example/server/hello"
	"log"
	"net"
)

func main() {
	// 建立 TCP 连接
	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalln("failed to listen: ", err)
	}
	// 创建 gRPC 服务
	// s := grpc.NewServer()
	s := grpc.NewServer(
	// grpc.UnaryInterceptor(hello.UnaryServerInterceptor), // 普通拦截器 获取一些基本信息
	// grpc.UnaryInterceptor(hello.SentryUnaryServerInterceptor()), // 鉴权拦截器
	)

	// 注册服务，两个参数，将 server 结构体的方法进行注册
	pb.RegisterHelloServer(s, &hello.Server{})

	go func() {
		log.Println("监听端口: ", listener.Addr())
		if err := s.Serve(listener); err != nil {
			log.Fatalln("failed to serve: ", err)
		}
	}()

	closehelp.Register(func(ctx context.Context) {
		s.GracefulStop()
		// listener.Close()
		// s.Stop()
	})
	closehelp.SignalClose()
}
