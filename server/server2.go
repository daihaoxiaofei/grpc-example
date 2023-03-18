package main

import (
	"google.golang.org/grpc"
	"grpc-example/pb"
	"grpc-example/server/hello"
	"log"
	"net"
)

func main() {
	// 建立 TCP 连接
	lis, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatalln("failed to listen: ", err)
	}
	// 创建 gRPC 服务
	s := grpc.NewServer()
	// s := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor))
	// 注册服务，两个参数，将 server 结构体的方法进行注册
	pb.RegisterHelloServer(s, &hello.Server{})
	log.Println("监听端口: ", lis.Addr())
	// 运行服务
	if err := s.Serve(lis); err != nil {
		log.Fatalln("failed to serve: ", err)
	}
}
