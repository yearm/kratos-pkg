package xgrpc

import (
	"fmt"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/yearm/kratos-pkg/config/gconfig"
	"github.com/yearm/kratos-pkg/errors"
)

// NewGRPCServer creates a GRPC server.
func NewGRPCServer(opts ...kgrpc.ServerOption) (*kgrpc.Server, error) {
	c, err := gconfig.GetServerGRPCConfig()
	if err != nil {
		return nil, errors.Wrap(err, "gconfig.GetServerGRPCConfig failed")
	}
	baseOptions := []kgrpc.ServerOption{
		kgrpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		kgrpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		kgrpc.Address(fmt.Sprintf("%s:%d", c.Host, c.Port)),
		kgrpc.Timeout(0), // Setting timeout to 0 delegates timeout control to client's context.
	}
	serverOptions := append(baseOptions, opts...)
	srv := kgrpc.NewServer(serverOptions...)
	if c.EnableHandlingTimeHistogram {
		grpc_prometheus.EnableHandlingTimeHistogram()
	}
	grpc_prometheus.Register(srv.Server)
	return srv, nil
}
