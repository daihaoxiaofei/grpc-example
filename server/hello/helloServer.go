package hello

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"grpc-example/pb"
	"io"
	"log"
	"strconv"
	"time"
)

type Server struct {
	pb.UnimplementedHelloServer // 可以被嵌入以具有向前兼容的实现。
}

func (s *Server) Hello(ctx context.Context, in *pb.ParBody) (*pb.ParBody, error) {
	res := &pb.ParBody{Value: ":8001 " + in.GetValue()}
	log.Printf("service 收到: %v", in.GetValue())
	// time.Sleep(time.Second * 3)
	return res, nil
}

func (s *Server) Channel(stream pb.Hello_ChannelServer) error {
	t := time.NewTicker(time.Millisecond * 500)
	for i := 0; ; i++ {
		// 接收消息
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		fmt.Println("服务端接收客户端的消息", recv)
		<-t.C
		// 发送消息
		rsp := &pb.ParBody{Value: strconv.Itoa(i)}
		err = stream.Send(rsp)
		if err != nil {
			return err
		}
	}
}

// UnaryServerInterceptor 普通拦截器 获取一些基本信息
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	remote, _ := peer.FromContext(ctx)
	remoteAddr := remote.Addr.String()

	start := time.Now()
	defer func() {
		in, _ := json.Marshal(req)
		out, _ := json.Marshal(resp)
		log.Println("ip", remoteAddr, "access_end", info.FullMethod, "in", string(in), "out", string(out),
			"err", err, "duration/ms", time.Since(start).Milliseconds())
	}()

	// 获取附加信息
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		accessTokenList := md.Get("access-token")
		fmt.Println(`accessTokenList`, accessTokenList, len(accessTokenList))

		header := md.Get("my-header")
		fmt.Println(`my-header`, header, len(header))

		metadataHeader := md.Get("metadata-my-header")
		fmt.Println(`metadata-my-header`, metadataHeader, len(metadataHeader))
	}

	resp, err = handler(ctx, req)

	return
}

//  鉴权连接器 以下
// 所有不需要用户登录的接口，就放在这里，否则则必须登录
var publicAPIMapper = map[string]bool{
	"/cashapp.CashApp/PingPong": true,
	"/cashapp.CashApp/Register": true,
	"/cashapp.CashApp/Login":    true,
}

func IsPublicAPI(fullMethodName string) bool {
	return publicAPIMapper[fullMethodName]
}

// SentryUnaryServerInterceptor 鉴权拦截器
func SentryUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (result interface{}, err error) {
		fullMethodName := info.FullMethod

		if !IsPublicAPI(fullMethodName) {
			accessToken := GetAccessToken(ctx)
			userID, err := GetUserIDByAccessToken(accessToken)
			if err != nil || userID == 0 {
				log.Printf("failed to find user by %s: %s", accessToken, err)
				return nil, status.Errorf(codes.Unauthenticated, ``, `err.`)
			}

			ctx = context.WithValue(ctx, "access_token", accessToken)
			ctx = context.WithValue(ctx, "user_id", userID)
		}

		log.Printf("got request with %v", req)
		result, err = handler(ctx, req)

		return
	}
}

// GetUserIDByAccessToken 通过accessToken得到用户id   accessToken 应该是需要调用另一些登录方法获得的
func GetUserIDByAccessToken(accessToken string) (int, error) {
	return 1, nil
}

// GetAccessToken 从ctx中拿metadata，然后从中拿access-token
// 这个时候，客户端只需要统一添加一个 Access-Token 的头部，值为对应用户的 access_token 即可。
func GetAccessToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	accessTokenList := md.Get("access-token")
	if len(accessTokenList) == 1 {
		return accessTokenList[0]
	}

	return ""
}
