package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"grpc-example/pb"
)

// UnaryClientInterceptor 客户端拦截器
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	start := time.Now()
	defer func() {
		in, _ := json.Marshal(req)
		out, _ := json.Marshal(reply)
		inStr, outStr := string(in), string(out)
		duration := int64(time.Since(start) / time.Millisecond)

		log.Println("grpc", method, "in", inStr, "out", outStr, "err", err, "duration/ms", duration)

	}()
	return invoker(ctx, method, req, reply, cc, opts...)
}

// 附加头部内容 header
type extraMetadata struct {
	MyHeader string `json:"my-header"`
}

func (c extraMetadata) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"my-header": c.MyHeader,
	}, nil
}

func (c extraMetadata) RequireTransportSecurity() bool {
	return false
}

func main() {
	// 服务注册
	db := &ColonyBuilder{
		address: map[string][]string{"/colony": {"localhost:8001", "localhost:8002", "localhost:8003"}},
	}

	resolver.Register(db) // 注册

	conn, err := grpc.Dial(db.Scheme()+`:///colony`, // 负载集群
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "`+roundrobin.Name+`"}`), // 负载均衡策略: 轮询
		// conn, err := grpc.Dial(`localhost:8001`,// 请求单台机器
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 传输安全集
		// grpc.WithUnaryInterceptor(UnaryClientInterceptor), // 客户端拦截器

		grpc.WithPerRPCCredentials(extraMetadata{MyHeader: "from grpc.Dial grpc.PerRPCCredentials"}), // 附加头部内容 header
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)
	push(c)
	// pushStream(c)
}

// 发送普通信息
func push(c pb.HelloClient) {
	// 官方给的附加内容方式
	// ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
	// 	"metadata-my-header": "from ctx",
	// }))
	ctx := context.Background()
	for i := 0; i < 1; i++ {
		// for range time.NewTicker(time.Millisecond * 500).C {
		r, err := c.Hello(ctx, &pb.ParBody{Value: strconv.Itoa(i)}) // grpc.PerRPCCredentials(extraMetadata{MyHeader: "from c.Hello grpc.PerRPCCredentials"}), // 附加头部内容 header

		if err != nil {
			log.Fatal("could not greet: ", err)
		}
		log.Printf("返回: %s", r.GetValue())
	}
}

// 发送双向流
func pushStream(c pb.HelloClient) {
	stream, err := c.Channel(context.Background())
	if err != nil {
		fmt.Println(`c.Channel err`, err)
		return
	}
	t := time.NewTicker(time.Millisecond * 500)
	for i := 0; ; i++ {
		// 发送消息
		err = stream.Send(&pb.ParBody{Value: strconv.Itoa(i)})
		if err != nil {
			log.Fatal(err)
		}
		<-t.C
		recv, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("客户端接收服务端的消息", recv.Value)
	}

}
