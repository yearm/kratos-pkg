package gconfig

import (
	"github.com/yearm/kratos-pkg/errors"
	"sync"
)

var (
	// defaultLogFileConfigKey default key for file log configuration.
	defaultLogFileConfigKey = "log.file"
	logFileConfigKeyOnce    sync.Once
)

// LogFileConfig file log config.
type LogFileConfig struct {
	Path  string `json:"path"`
	Level string `json:"level"`
}

func (l *LogFileConfig) isValid() bool {
	return l.Path != "" && l.Level != ""
}

// SetLogFileConfigKey customizes the global config key for file log.
func SetLogFileConfigKey(key string) {
	logFileConfigKeyOnce.Do(func() {
		defaultLogFileConfigKey = key
	})
}

// GetLogFileConfig retrieves file logging configuration from global settings.
func GetLogFileConfig() (*LogFileConfig, error) {
	var config *LogFileConfig
	err := Value(defaultLogFileConfigKey).Scan(&config)
	if err != nil {
		return nil, errors.Wrapf(err, "config.Scan[%v] faild", defaultLogFileConfigKey)
	}
	if !config.isValid() {
		return nil, errors.Errorf("log file config[%v] is invalid", config)
	}
	return config, nil
}

var (
	// defaultLogAliyunConfigKey default key for aliyun log configuration.
	defaultLogAliyunConfigKey = "log.aliyun"
	logAliyunConfigKeyOnce    sync.Once
)

// LogAliyunConfig aliyun log config.
type LogAliyunConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Endpoint  string `json:"endpoint"`
	Project   string `json:"project"`
	Logstore  string `json:"logstore"`
	Level     string `json:"level"`
}

func (l *LogAliyunConfig) isValid() bool {
	return l.AccessKey != "" && l.SecretKey != "" && l.Endpoint != "" && l.Project != "" && l.Logstore != "" && l.Level != ""
}

// SetLogAliyunConfigKey customizes the global config key for aliyun log.
func SetLogAliyunConfigKey(key string) {
	logAliyunConfigKeyOnce.Do(func() {
		defaultLogAliyunConfigKey = key
	})
}

// GetLogAliyunConfig retrieves aliyun logging configuration from global settings.
func GetLogAliyunConfig() (*LogAliyunConfig, error) {
	var config *LogAliyunConfig
	err := Value(defaultLogAliyunConfigKey).Scan(&config)
	if err != nil {
		return nil, errors.Wrapf(err, "config.Scan[%v] faild", defaultLogAliyunConfigKey)
	}
	if !config.isValid() {
		return nil, errors.Errorf("log aliyun config[%v] is invalid", defaultLogAliyunConfigKey)
	}
	return config, nil
}
