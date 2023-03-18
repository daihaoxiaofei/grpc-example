package main

import "google.golang.org/grpc/resolver"

type ColonyBuilder struct {
	address map[string][]string
}

// Build 参数 1 目标 2 连接
func (b ColonyBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &ColonyResolver{}
	paths := b.address[target.URL.Path]
	addresses := make([]resolver.Address, len(paths))
	for i, s := range paths {
		addresses[i] = resolver.Address{Addr: s}
	}
	// 更新状态
	_ = cc.UpdateState(resolver.State{Addresses: addresses})
	return r, nil
}

func (b ColonyBuilder) Scheme() string {
	return `resolver`
}

type ColonyResolver struct {
}

func (c ColonyResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (c ColonyResolver) Close() {
}
