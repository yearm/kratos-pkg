package gconfig

import (
	"github.com/yearm/kratos-pkg/errors"
	"sync"
)

var (
	// defaultServerGRPCConfigKey default key for grpc server configuration.
	defaultServerGRPCConfigKey = "server.grpc"
	serverGRPCConfigKeyOnce    sync.Once
)

// SetServerGRPCConfigKey customizes the global config key for grpc server.
func SetServerGRPCConfigKey(key string) {
	serverGRPCConfigKeyOnce.Do(func() {
		defaultServerGRPCConfigKey = key
	})
}

// ServerGRPCConfig grpc server config.
type ServerGRPCConfig struct {
	Host                        string `json:"host"`
	Port                        int    `json:"port"`
	EnableHandlingTimeHistogram bool   `json:"enableHandlingTimeHistogram"`
}

func (s *ServerGRPCConfig) isValid() bool {
	return s.Host != "" && s.Port != 0
}

// GetServerGRPCConfig retrieves grpc server configuration from global settings.
func GetServerGRPCConfig() (*ServerGRPCConfig, error) {
	var config *ServerGRPCConfig
	err := Value(defaultServerGRPCConfigKey).Scan(&config)
	if err != nil {
		return nil, errors.Wrapf(err, "config.Scan[%v] faild", defaultServerGRPCConfigKey)
	}
	if !config.isValid() {
		return nil, errors.Errorf("server grpc config[%v] is invalid", defaultServerGRPCConfigKey)
	}
	return config, nil
}

var (
	// defaultServerHTTPConfigKey default key for http server configuration.
	defaultServerHTTPConfigKey = "server.http"
	serverHTTPConfigKeyOnce    sync.Once
)

// SetServerHTTPConfigKey customizes the global config key for http server.
func SetServerHTTPConfigKey(key string) {
	serverHTTPConfigKeyOnce.Do(func() {
		defaultServerHTTPConfigKey = key
	})
}

// CorsConfig http cors config.
type CorsConfig struct {
	AllowOrigins     []string `json:"allowOrigins"`
	AllowMethods     []string `json:"allowMethods"`
	AllowHeaders     []string `json:"allowHeaders"`
	ExposeHeaders    []string `json:"exposeHeaders"`
	MaxAge           int      `json:"maxAge"`
	AllowCredentials bool     `json:"allowCredentials"`
}

// ServerHTTPConfig http server config.
type ServerHTTPConfig struct {
	Host    string      `json:"host"`
	Port    int         `json:"port"`
	Timeout int         `json:"timeout"`
	Cors    *CorsConfig `json:"cors"`
}

func (s *ServerHTTPConfig) isValid() bool {
	return s.Host != "" && s.Port != 0 && s.Timeout >= 0
}

// GetServerHTTPConfig retrieves http server configuration from global settings.
func GetServerHTTPConfig() (*ServerHTTPConfig, error) {
	var config *ServerHTTPConfig
	err := Value(defaultServerHTTPConfigKey).Scan(&config)
	if err != nil {
		return nil, errors.Wrapf(err, "config.Scan[%v] faild", defaultServerHTTPConfigKey)
	}
	if !config.isValid() {
		return nil, errors.Errorf("server http config[%v] is invalid", defaultServerHTTPConfigKey)
	}
	return config, nil
}

var (
	// defaultServerMonitorHTTPConfigKey default key for monitor http server configuration.
	defaultServerMonitorHTTPConfigKey = "server.monitorHttp"
	serverMonitorHTTPConfigKeyOnce    sync.Once
)

// SetServerMonitorHTTPConfigKey customizes the global config key for monitor http server.
func SetServerMonitorHTTPConfigKey(key string) {
	serverMonitorHTTPConfigKeyOnce.Do(func() {
		defaultServerMonitorHTTPConfigKey = key
	})
}

// ServerMonitorHTTPConfig monitor http server config.
type ServerMonitorHTTPConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (s *ServerMonitorHTTPConfig) isValid() bool {
	return s.Host != "" && s.Port != 0
}

// GetServerMonitorHTTPConfig retrieves monitor http server configuration from global settings.
func GetServerMonitorHTTPConfig() (*ServerMonitorHTTPConfig, error) {
	var config *ServerMonitorHTTPConfig
	err := Value(defaultServerMonitorHTTPConfigKey).Scan(&config)
	if err != nil {
		return nil, errors.Wrapf(err, "config.Scan[%v] faild", defaultServerMonitorHTTPConfigKey)
	}
	if !config.isValid() {
		return nil, errors.Errorf("server monitor http config[%v] is invalid", defaultServerMonitorHTTPConfigKey)
	}
	return config, nil
}
