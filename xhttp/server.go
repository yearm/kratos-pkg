package xhttp

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yearm/kratos-pkg/config/gconfig"
	"github.com/yearm/kratos-pkg/errors"
	"github.com/yearm/kratos-pkg/utils/net"
)

// NewHTTPServer creates an http server.
func NewHTTPServer(handler http.Handler, opts ...khttp.ServerOption) (*khttp.Server, error) {
	return newHTTPServer(handler, "", opts...)
}

// NewHTTPServerByConfigKey creates an http server by config key.
func NewHTTPServerByConfigKey(handler http.Handler, configKey string, opts ...khttp.ServerOption) (*khttp.Server, error) {
	return newHTTPServer(handler, configKey, opts...)
}

// newHTTPServer creates an http server by config key.
func newHTTPServer(handler http.Handler, configKey string, opts ...khttp.ServerOption) (*khttp.Server, error) {
	c, err := gconfig.GetServerHTTPConfig()
	if err != nil {
		return nil, errors.Wrap(err, "gconfig.GetServerHTTPConfig failed")
	}
	if configKey != "" {
		err = gconfig.Value(configKey).Scan(&c)
		if err != nil {
			return nil, errors.Wrapf(err, "scan config[%v] failed", configKey)
		}
	}

	baseOptions := []khttp.ServerOption{
		khttp.Address(fmt.Sprintf("%s:%d", c.Host, c.Port)),
		khttp.Timeout(time.Duration(c.Timeout) * time.Second),
	}
	if c.Cors != nil {
		opts := []handlers.CORSOption{
			handlers.AllowedOrigins(c.Cors.AllowOrigins),
			handlers.AllowedMethods(c.Cors.AllowMethods),
			handlers.AllowedHeaders(c.Cors.AllowHeaders),
			handlers.ExposedHeaders(c.Cors.ExposeHeaders),
			handlers.MaxAge(c.Cors.MaxAge),
		}
		if c.Cors.AllowCredentials {
			opts = append(opts, handlers.AllowCredentials())
		}
		baseOptions = append(baseOptions, khttp.Filter(handlers.CORS(opts...)))
	}
	serverOptions := append(baseOptions, opts...)
	httpSrv := khttp.NewServer(serverOptions...)

	switch app := handler.(type) {
	case *gin.Engine:
		httpSrv.HandlePrefix("/", app)
		return httpSrv, nil
	case nil:
		return httpSrv, nil
	default:
		return nil, errors.Errorf("unsupported http.Handler")
	}
}

// NewMonitorHTTPServer creates a monitor HTTP server.
func NewMonitorHTTPServer() (*khttp.Server, error) {
	c, err := gconfig.GetServerMonitorHTTPConfig()
	if err != nil {
		return nil, errors.Wrap(err, "gconfig.NewMonitorHTTPServer failed")
	}
	if c.Port <= 0 {
		port, err := net.GetFreePort()
		if err != nil {
			return nil, errors.Wrap(err, "net.GetFreePort failed")
		}
		c.Port = port
	}
	opts := []khttp.ServerOption{
		khttp.Address(fmt.Sprintf("%s:%d", c.Host, c.Port)),
		khttp.Timeout(0),
	}
	httpServer := khttp.NewServer(opts...)
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
