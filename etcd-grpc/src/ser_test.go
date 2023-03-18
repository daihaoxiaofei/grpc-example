package src

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-example/pb"
	"log"
	"testing"
)

// 获取注册的grpc服务集合
func TestName(t *testing.T) {
	list, err := manager.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range list {
		fmt.Println(v)
	}
	fmt.Println(`完成`)
}

// 持续关注变化
func TestWatch(t *testing.T) {
	WatchChannel, err := manager.NewWatchChannel(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for v := range WatchChannel {
		for _, update := range v {
			fmt.Println(*update)
		}

	}
}

//  测试grpc服务
func TestPush(t *testing.T) {
	conn, err := grpc.Dial(
		`localhost:8001`,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 传输安全集
	) // 请求单台机器

	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)
	r, err := c.Hello(context.Background(), &pb.ParBody{Value: `strconv.Itoa(i)`})

	if err != nil {
		log.Fatal("could not greet: ", err)
	}
	log.Printf("返回: %s", r.GetValue())
}
