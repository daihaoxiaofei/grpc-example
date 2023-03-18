package src

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"grpc-example/pkg/etcdhelp"
)

var (
	SrvName = `grpcService`
	cli     = etcdhelp.Cli
	manager endpoints.Manager
	err     error
)

func init() {
	manager, err = endpoints.NewManager(cli, SrvName)
	if err != nil {
		panic(err)
	}
}

// Register 注册并保持
func Register(addr string) error {
	// get leaseID 续约组
	leaseGrantResponse, err := cli.Grant(context.TODO(), 1)
	if err != nil {
		return err
	}

	err = manager.AddEndpoint(
		context.Background(),
		SrvName+`/`+addr,
		endpoints.Endpoint{Addr: addr},
		clientv3.WithLease(leaseGrantResponse.ID), // 续约组
	)
	if err != nil {
		return err
	}

	ch, err := cli.KeepAlive(context.Background(), leaseGrantResponse.ID)
	if err != nil {
		return err
	}
	go func() {
		for { // 需要不断的取出lease的response
			<-ch
		}
	}()
	return nil
}
