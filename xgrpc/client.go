package xgrpc

import (
	"context"
	"crypto/tls"
	"fmt"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yearm/kratos-pkg/registry"
	"google.golang.org/grpc"
	"sync"
	"time"
)

var (
	// connMap ...
	connMap sync.Map
	// connLock ...
	connLock sync.Mutex
)

// GetRPCClientConn 配置结构规则 config/env/env.go:58
func GetRPCClientConn(configPath string, opts ...kgrpc.ClientOption) *grpc.ClientConn {
	var (
		err  error
		conn *grpc.ClientConn
	)

	connLock.Lock()
	defer connLock.Unlock()
	if v, ok := connMap.Load(configPath); ok {
		return v.(*grpc.ClientConn)
	}
	defer func() {
		if conn != nil {
			connMap.Store(configPath, conn)
		}
	}()

	endpoint := viper.GetString(fmt.Sprintf("%s.endpoint", configPath))
	if endpoint == "" {
		logrus.Panicln("endpoint is nil, config path:", configPath)
	}
	timeout := viper.GetInt(fmt.Sprintf("%s.timeout", configPath))
	dialWithCredentials := viper.GetBool(fmt.Sprintf("%s.dialWithCredentials", configPath))

	clientOpts := []kgrpc.ClientOption{kgrpc.WithEndpoint(endpoint)}
	clientOpts = append(clientOpts, opts...)
	if d := registry.NewDiscovery(); d != nil {
		clientOpts = append(clientOpts, kgrpc.WithDiscovery(d))
	}

	if timeout >= 0 {
		clientOpts = append(clientOpts, kgrpc.WithTimeout(time.Duration(timeout)*time.Second))
	}

	if dialWithCredentials {
		clientOpts = append(clientOpts, kgrpc.WithTLSConfig(&tls.Config{}))
		conn, err = kgrpc.Dial(context.Background(), clientOpts...)
	} else {
		conn, err = kgrpc.DialInsecure(context.Background(), clientOpts...)
	}

	if err != nil {
		logrus.Panicln("grpc dial error:", err)
	}
	logrus.Infoln("Connecting at", endpoint)
	return conn
}
