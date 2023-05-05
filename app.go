package kratospkg

import (
	"github.com/go-kratos/kratos/v2"
	kregistry "github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/yearm/kratos-pkg/config/env"
	_ "github.com/yearm/kratos-pkg/log"
	_ "go.uber.org/automaxprocs"
)

// NewHttpApp ...
func NewHttpApp(hs []*http.Server) *kratos.App {
	ts := make([]transport.Server, 0)
	for _, h := range hs {
		ts = append(ts, h)
	}
	return kratos.New(
		kratos.ID(env.GetServiceInstanceID()),
		kratos.Name(env.GetServiceName()),
		kratos.Version(env.GetServiceVersion()),
		kratos.Server(ts...),
	)
}

// NewGRPCApp ...
func NewGRPCApp(reg kregistry.Registrar, gs *grpc.Server, ms *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(env.GetServiceInstanceID()),
		kratos.Name(env.GetServiceName()),
		kratos.Version(env.GetServiceVersion()),
		kratos.Server(gs, ms),
		kratos.Registrar(reg),
	)
}
