package gconfig

import (
	"sync"

	"github.com/yearm/kratos-pkg/errors"
)

var (
	// defaultTraceConfigKey default key for trace configuration.
	defaultTraceConfigKey = "trace"
	traceConfigKeyOnce    sync.Once
)

type TraceExporterConfig struct {
	Endpoint string `json:"endpoint"`
}

type TraceConfig struct {
	Exporter *TraceExporterConfig `json:"exporter"`
}

// SetTraceConfigKey customizes the global config key for trace.
func SetTraceConfigKey(key string) {
	traceConfigKeyOnce.Do(func() {
		defaultTraceConfigKey = key
	})
}

// GetTraceConfig retrieves trace configuration from global settings.
func GetTraceConfig() (*TraceConfig, error) {
	var config *TraceConfig
	err := Value(defaultTraceConfigKey).Scan(&config)
	if err != nil {
		return nil, errors.Wrapf(err, "config.Scan[%v] faild", defaultTraceConfigKey)
	}
	return config, nil
}
