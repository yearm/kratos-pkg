package xgrpc

import (
	"fmt"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/yearm/kratos-pkg/config/env"
)

// NewGRPCServer ...
func NewGRPCServer(opts ...kgrpc.ServerOption) *kgrpc.Server {
	var serverOpts = []kgrpc.ServerOption{
		kgrpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		kgrpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	}
	serverOpts = append(serverOpts, opts...)
	grpcHost := env.GetGRPCHost()
	grpcPort := env.GetGRPCPort()
	if grpcHost != "" && grpcPort > 0 {
		serverOpts = append(serverOpts, kgrpc.Address(fmt.Sprintf("%s:%d", grpcHost, grpcPort)))
	}
	// NOTE: kgrpc.Timeout 暂不用设置，因为 unaryServerInterceptor 中并没有 select case <-ctx.Done()
	// NOTE: 设置了反倒会改变 context 的超时传递时间，一般情况 client 的 context 带有超时时间
	// NOTE: 正常情况下 client 调用超时是为了避免链路阻塞堆积，server 端继续处理请求也属正常
	// NOTE: 未自定义设置服务端 timeout 时，kratos 框架默认设置为1秒，导致服务端调用时间过长或者链路较长时服务超时中断
	// NOTE: 所以此处设置 timeout 为0，即使用客户端调用传来的 ctx 中的超时控制
	serverOpts = append(serverOpts, kgrpc.Timeout(0))
	srv := kgrpc.NewServer(serverOpts...)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(srv.Server)
	return srv
}
