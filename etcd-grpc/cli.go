package main

import (
	"fmt"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-example/etcd-grpc/src"
	"grpc-example/pb"
	"grpc-example/pkg/etcdhelp"
	"log"
	"strconv"
	"time"
)

func main() {
	etcdResolver, err := resolver.NewBuilder(etcdhelp.Cli)

	grpcCli, err := grpc.Dial(
		fmt.Sprintf(etcdResolver.Scheme()+":///"+src.SrvName),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 传输安全集
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "`+roundrobin.Name+`"}`), // 轮询
	)
	if err != nil {
		fmt.Println("连接服务器失败：", err)
	}
	defer grpcCli.Close()
	// 获得grpc句柄
	c := pb.NewHelloClient(grpcCli)
	t := time.NewTicker(time.Second)
	for i := 0; i < 100; i++ {
		r, err := c.Hello(context.Background(), &pb.ParBody{Value: strconv.Itoa(i)})
		if err != nil {
			log.Fatal("could not greet: ", err)
		}
		log.Printf("返回: %s", r.GetValue())
		<-t.C
	}
}
