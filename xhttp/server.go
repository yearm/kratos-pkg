package xhttp

import (
	"fmt"
	thttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"github.com/yearm/kratos-pkg/config/env"
	"net/http"
	"net/http/pprof"
	"time"
)

// NewHTTPServer ...
func NewHTTPServer(handler http.Handler, configPath ...string) (*thttp.Server, error) {
	opts := make([]thttp.ServerOption, 0)
	var (
		httpHost              string
		httpPort, httpTimeout int
	)
	if len(configPath) > 0 {
		_configPath := configPath[0]
		httpHost = viper.GetString(fmt.Sprintf("%s.host", _configPath))
		httpPort = viper.GetInt(fmt.Sprintf("%s.port", _configPath))
		httpTimeout = viper.GetInt(fmt.Sprintf("%s.timeout", _configPath))
	} else {
		httpHost = env.GetHttpHost()
		httpPort = env.GetHttpPort()
		httpTimeout = env.GetHttpTimeout()
	}
	if httpHost != "" && httpPort > 0 {
		opts = append(opts, thttp.Address(fmt.Sprintf("%s:%d", httpHost, httpPort)))
	}
	if httpTimeout >= 0 {
		// NOTE: context 的超时时间
		opts = append(opts, thttp.Timeout(time.Duration(httpTimeout)*time.Second))
	}
	httpSrv := thttp.NewServer(opts...)

	switch app := handler.(type) {
	case *iris.Application:
		httpSrv.HandlePrefix("/", app)
		return httpSrv, app.Build()
	default:
		return nil, fmt.Errorf("unsupported http.Handler")
	}
}

// NewMonitorHTTPServer ...
func NewMonitorHTTPServer() (*thttp.Server, error) {
	httpHost := env.GetMonitorHttpHost()
	httpPort := env.GetMonitorHttpPort()
	if httpPort <= 0 {
		httpPort = gtcp.MustGetFreePort()
	}
	opts := make([]thttp.ServerOption, 0)
	if env.GetMonitorHttpHost() != "" && httpPort > 0 {
		opts = append(opts, thttp.Address(fmt.Sprintf("%s:%d", httpHost, httpPort)))
	}
	opts = append(opts, thttp.Timeout(0))
	httpServer := thttp.NewServer(opts...)
	httpServer.HandleFunc("/debug/pprof/", pprof.Index)
	httpServer.HandleFunc("/debug/pprof/allocs", pprof.Index)
	httpServer.HandleFunc("/debug/pprof/block", pprof.Index)
	httpServer.HandleFunc("/debug/pprof/goroutine", pprof.Index)
	httpServer.HandleFunc("/debug/pprof/heap", pprof.Index)
	httpServer.HandleFunc("/debug/pprof/mutex", pprof.Index)
	httpServer.HandleFunc("/debug/pprof/threadcreate", pprof.Index)
	httpServer.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpServer.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpServer.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpServer.HandleFunc("/debug/pprof/trace", pprof.Trace)
	httpServer.Handle("/metrics", promhttp.Handler())
	return httpServer, nil
}
