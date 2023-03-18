package etcdhelp

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"grpc-example/pkg/closehelp"
	"time"
)

var Cli *clientv3.Client

func init() {
	var err error
	Cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(`etcd init err ` + err.Error())
	}

	closehelp.Register(func(ctx context.Context) {
		Cli.Close()
	})
}

// resp, err := m.client.Get(ctx, m.target, clientv3.WithPrefix(), clientv3.WithSerializable())
