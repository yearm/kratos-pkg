package kratospkg

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/yearm/kratos-pkg/env"
	_ "go.uber.org/automaxprocs"
)

// NewApp creates a kratos application.
func NewApp(ss []transport.Server, opts []kratos.Option) *kratos.App {
	options := []kratos.Option{
		kratos.ID(env.GetServiceID()),
		kratos.Name(env.GetServiceName()),
		kratos.Version(env.GetServiceVersion()),
		kratos.Metadata(env.GetServiceMetadata()),
		kratos.Server(ss...),
	}
	options = append(options, opts...)
	return kratos.New(options...)
}
