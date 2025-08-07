package xgrpc

import (
	"context"
	ctls "crypto/tls"
	"log"
	"sync"
	"time"

	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/yearm/kratos-pkg/config/gconfig"
	"github.com/yearm/kratos-pkg/errors"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
)

var (
	// connMap thread-safe cache for storing gRPC connections keyed by endpoint.
	connMap sync.Map

	// sg deduplicates concurrent connection requests.
	sg singleflight.Group
)

// GetGRPCClientConnByConfigKey creates a grpc client conn by config key.
func GetGRPCClientConnByConfigKey(key string, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	c, err := gconfig.GetClientGRPCConfig(key)
	if err != nil {
		return nil, errors.Wrap(err, "gconfig.GetClientGRPCConfig failed")
	}
	return GetGRPCClientConn(c.Endpoint, c.Timeout, c.TLS, opts...)
}

// GetGRPCClientConn creates a grpc client conn.
func GetGRPCClientConn(endpoint string, timeout int, tls bool, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint is required")
	}

	c, err, _ := sg.Do(endpoint, func() (interface{}, error) {
		if c, ok := connMap.Load(endpoint); ok {
			return c, nil
		}

		baseOptions := []kgrpc.ClientOption{
			kgrpc.WithEndpoint(endpoint),
			kgrpc.WithTimeout(time.Duration(timeout) * time.Second),
			kgrpc.WithOptions(
				grpc.WithIdleTimeout(0), // disable idle timeout
			),
			kgrpc.WithPrintDiscoveryDebugLog(false),
		}
		clientOptions := append(baseOptions, opts...)

		var (
			conn *grpc.ClientConn
			err  error
		)
		if tls {
			clientOptions = append(clientOptions, kgrpc.WithTLSConfig(&ctls.Config{}))
			conn, err = kgrpc.Dial(context.Background(), clientOptions...)
			if err != nil {
				return nil, errors.Wrap(err, "grpc.Dial failed")
			}
		} else {
			conn, err = kgrpc.DialInsecure(context.Background(), clientOptions...)
			if err != nil {
				return nil, errors.Wrap(err, "grpc.DialInsecure failed")
			}
		}

		log.Println("Connecting at", endpoint)
		connMap.Store(endpoint, conn)
		return conn, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "GetGRPCClientConn failed")
	}
	return c.(*grpc.ClientConn), nil
}
