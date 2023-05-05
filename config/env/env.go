package env

import (
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/spf13/viper"
	"os"
)

// aliyun.json
// Example:
/*
{
  "aliyun": {
    "auth": {
      "soda": {
        "accessKey": "accessKey",
        "secretKey": "secretKey"
      }
    },
    "mse": {
      "endpoint": "endpoint",
      "namespaceId": "namespaceId",
      "regionId": "cn-hangzhou",
      "cacheDir": "./config/cache",
      "logDir": "./tmp/mse/log",
      "logLevel": "error",
      "dataId": "dataId",
      "group": "group"
    }
  }
}
*/
const (
	PathAliyunAccessKey = "aliyun.auth.soda.accessKey"
	PathAliyunSecretKey = "aliyun.auth.soda.secretKey"

	PathAliyunMseEndpoint    = "aliyun.mse.endpoint"
	PathAliyunMseNamespaceID = "aliyun.mse.namespaceId"
	PathAliyunMseRegionID    = "aliyun.mse.regionId"
	PathAliyunMseCacheDir    = "aliyun.mse.cacheDir"
	PathAliyunMseLogDir      = "aliyun.mse.logDir"
	PathAliyunMseLogLevel    = "aliyun.mse.logLevel"
	PathAliyunMseDataID      = "aliyun.mse.dataId"
	PathAliyunMseGroup       = "aliyun.mse.group"
)

// GetAliyunAccessKey ...
func GetAliyunAccessKey() string {
	return viper.GetString(PathAliyunAccessKey)
}

// GetAliyunSecretKey ...
func GetAliyunSecretKey() string {
	return viper.GetString(PathAliyunSecretKey)
}

// runtime.config.json
// Example:
/*
{
  "mode": "dev",
  "server": {
    "name": "name",
    "version": "1.0.0",
    "httpServer": {
      "httpServer": {
        "host": "0.0.0.0",
        "port": 9999,
        "timeout": 10
      },
      "monitorHttp": {
        "host": "0.0.0.0",
        "port": 9104
      }
    },
    "rpcServer": {
      "grpc": {
        "host": "0.0.0.0",
        "port": 9999
      }
    }
  },
  "resource": {
    "database": {
      "soda": {
        "wr": {
          "dialect": "mysql",
          "host": "",
          "port": 3306,
          "user": "",
          "password": "",
          "database": "soda",
          "maxIdle": 20,
          "maxOpen": 20,
          "maxLifetime": 300
        }
      }
    },
    "redis": {
      "default": {
        "addr": "",
        "password": "",
        "database": 10,
        "prefix": "",
        "maxIdle": 20,
        "maxActive": 50,
        "idleTimeout": 60
      },
      "mutex": {
        "expiry": 6,
        "tries": 15
      }
    },
    "aliyun": {
      "sls": {
        "endpoint": "",
        "project": "project",
        "logStore": "logStore",
        "level": "warn"
      },
      "trace": {
        "url": ""
      }
    }
  }
}
*/
const (
	PathMode = "mode"

	PathServiceName    = "server.name"
	PathServiceVersion = "server.version"

	PathMonitorHttpHost = "server.httpServer.monitorHttp.host"
	PathMonitorHttpPort = "server.httpServer.monitorHttp.port"

	PathHttpHost    = "server.httpServer.http.host"
	PathHttpPort    = "server.httpServer.http.port"
	PathHttpTimeout = "server.httpServer.http.timeout"

	PathGRPCHost = "server.rpcServer.grpc.host"
	PathGRPCPort = "server.rpcServer.grpc.port"

	PathAliyunSlsEndpoint = "resource.aliyun.sls.endpoint"
	PathAliyunSlsProject  = "resource.aliyun.sls.project"
	PathAliyunSlsLogStore = "resource.aliyun.sls.logStore"
	PathAliyunSlsLogLevel = "resource.aliyun.sls.level"

	PathAliyunTraceEndpoint = "resource.aliyun.trace.url"

	PathAliyunSlsTraceEndpoint   = "resource.aliyun.slsTrace.url"
	PathAliyunSlsTraceProject    = "resource.aliyun.slsTrace.project"
	PathAliyunSlsTraceInstanceID = "resource.aliyun.slsTrace.instanceId"

	PathRedisDefault = "resource.redis.default"

	PathRedisMutex = "resource.redis.mutex"
)

// GetMode ...
func GetMode() string {
	return viper.GetString(PathMode)
}

// GetServiceInstanceID ...
func GetServiceInstanceID() string {
	name, _ := os.Hostname()
	return name
}

// GetServiceName ...
func GetServiceName() string {
	return viper.GetString(PathServiceName)
}

// GetServiceVersion ...
func GetServiceVersion() string {
	return viper.GetString(PathServiceVersion)
}

// GetMonitorHttpHost ...
func GetMonitorHttpHost() string {
	return viper.GetString(PathMonitorHttpHost)
}

// GetMonitorHttpPort ...
func GetMonitorHttpPort() int {
	return viper.GetInt(PathMonitorHttpPort)
}

// GetHttpHost ...
func GetHttpHost() string {
	return viper.GetString(PathHttpHost)
}

// GetHttpPort ...
func GetHttpPort() int {
	return viper.GetInt(PathHttpPort)
}

// GetHttpTimeout ...
func GetHttpTimeout() int {
	return viper.GetInt(PathHttpTimeout)
}

// GetGRPCHost ...
func GetGRPCHost() string {
	return viper.GetString(PathGRPCHost)
}

// GetGRPCPort ...
func GetGRPCPort() int {
	return viper.GetInt(PathGRPCPort)
}

// GetAliyunTraceEndpoint ...
func GetAliyunTraceEndpoint() string {
	return viper.GetString(PathAliyunTraceEndpoint)
}

// GetAliyunSlsTraceEndpoint ...
func GetAliyunSlsTraceEndpoint() string {
	return viper.GetString(PathAliyunSlsTraceEndpoint)
}

// GetAliyunSlsTraceProject ...
func GetAliyunSlsTraceProject() string {
	return viper.GetString(PathAliyunSlsTraceProject)
}

// GetAliyunSlsTraceInstanceID ...
func GetAliyunSlsTraceInstanceID() string {
	return viper.GetString(PathAliyunSlsTraceInstanceID)
}

// GetRedisDefaultConfigPath ...
func GetRedisDefaultConfigPath() string {
	return PathRedisDefault
}

// GetRedisMutexConfigPath ..
func GetRedisMutexConfigPath() string {
	return PathRedisMutex
}

// modes
const (
	ModeLocal      = "local"
	ModeDevelop    = "dev"
	ModeTest       = "test"
	ModeProduction = "prod"
)

// IsDevelopment ...
func IsDevelopment() bool {
	return gstr.InArray([]string{ModeLocal, ModeDevelop, ModeTest}, GetMode())
}

// IsProduction ...
func IsProduction() bool {
	return GetMode() == ModeProduction
}

// IsLocal ...
func IsLocal() bool {
	return GetMode() == ModeLocal
}
