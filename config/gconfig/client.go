package gconfig

import (
	"fmt"
	"sync"

	"github.com/yearm/kratos-pkg/errors"
)

var (
	// defaultClientGRPCConfigKey default key for grpc client configuration.
	defaultClientGRPCConfigKey = "client.grpc"
	clientGRPCConfigKeyOnce    sync.Once
)

// SetClientGRPCConfigKey customizes the global config key for grpc client.
func SetClientGRPCConfigKey(key string) {
	clientGRPCConfigKeyOnce.Do(func() {
		defaultClientGRPCConfigKey = key
	})
}

// ClientGRPCConfig grpc client config.
type ClientGRPCConfig struct {
	Endpoint string `json:"endpoint"`
	Timeout  int    `json:"timeout"`
	TLS      bool   `json:"tls"`
}

func (c *ClientGRPCConfig) isValid() bool {
	return c.Endpoint != "" && c.Timeout >= 0
}

// GetClientGRPCConfig retrieves grpc client configuration from global settings.
func GetClientGRPCConfig(name string) (*ClientGRPCConfig, error) {
	key := fmt.Sprintf("%s.%s", defaultClientGRPCConfigKey, name)
	var config *ClientGRPCConfig
	err := Value(key).Scan(&config)
	if err != nil {
		return nil, errors.Wrapf(err, "config.Scan[%v] faild", key)
	}
	if !config.isValid() {
		return nil, errors.Errorf("client grpc config[%v] is invalid", key)
	}
	return config, nil
}
