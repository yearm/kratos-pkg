package xgrpc

import (
	"context"
	"crypto/tls"
	"fmt"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yearm/kratos-pkg/registry"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"sync"
	"time"
)

// GetRPCClientConn 配置结构规则 config/env/env.go:58
func GetRPCClientConn(configPath string, opts ...kgrpc.ClientOption) *grpc.ClientConn {
	endpoint := viper.GetString(fmt.Sprintf("%s.endpoint", configPath))
	timeout := viper.GetInt(fmt.Sprintf("%s.timeout", configPath))
	dialWithCredentials := viper.GetBool(fmt.Sprintf("%s.dialWithCredentials", configPath))
	if endpoint == "" {
		logrus.Panicln("endpoint is nil, config path:", configPath)
	}
	conn, err := dial(endpoint, timeout, dialWithCredentials, opts...)
	if err != nil {
		logrus.Panicln(err)
	}
	return conn
}

// GetClientConn ...
func GetClientConn(endpoint string, timeout int, dialWithCredentials bool, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is nil")
	}
	conn, err := dial(endpoint, timeout, dialWithCredentials, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

var (
	// connMap ...
	connMap sync.Map
	// sg ...
	sg singleflight.Group
)

func dial(endpoint string, timeout int, dialWithCredentials bool, opts ...kgrpc.ClientOption) (*grpc.ClientConn, error) {
	iConn, err, _ := sg.Do(endpoint, func() (interface{}, error) {
		var (
			err  error
			conn *grpc.ClientConn
		)
		if conn, ok := connMap.Load(endpoint); ok {
			return conn, nil
		}
		defer func() {
			if conn != nil {
				logrus.Infoln("Connecting at", endpoint)
				connMap.Store(endpoint, conn)
			}
		}()

		clientOpts := []kgrpc.ClientOption{
			kgrpc.WithEndpoint(endpoint),
			kgrpc.WithOptions(grpc.WithIdleTimeout(0)),
		}
		if timeout >= 0 {
			clientOpts = append(clientOpts, kgrpc.WithTimeout(time.Duration(timeout)*time.Second))
		}
		if d := registry.NewDiscovery(); d != nil {
			clientOpts = append(clientOpts, kgrpc.WithDiscovery(d))
		}
		clientOpts = append(clientOpts, opts...)

		if dialWithCredentials {
			clientOpts = append(clientOpts, kgrpc.WithTLSConfig(&tls.Config{}))
			conn, err = kgrpc.Dial(context.Background(), clientOpts...)
		} else {
			conn, err = kgrpc.DialInsecure(context.Background(), clientOpts...)
		}
		if err != nil {
			return nil, fmt.Errorf("grpc dial error: %s", err)
		}
		return conn, nil
	})
	if err != nil {
		return nil, err
	}
	return iConn.(*grpc.ClientConn), nil
}
