package xiris

import (
	"github.com/kataras/iris/v12"
	iriscontext "github.com/kataras/iris/v12/context"
	"github.com/yearm/kratos-pkg/ecode"
	"net/http"
	"sync"
)

var (
	// WithRemoteAddrHeaders ...
	WithRemoteAddrHeaders = func(app *iris.Application) {
		app.Configure(iris.WithRemoteAddrHeader("X-Forwarded-For"), iris.WithRemoteAddrHeader("X-Real-Ip"))
	}
	// DefaultConfigures ...
	DefaultConfigures = []iris.Configurator{iris.WithoutBodyConsumptionOnUnmarshal, iris.WithoutPathCorrectionRedirection, WithRemoteAddrHeaders}
)

// RegisterOnErrorCode ...
func RegisterOnErrorCode(app *iris.Application, codes ...int) {
	for _, code := range codes {
		switch code {
		case http.StatusInternalServerError:
			app.OnErrorCode(http.StatusInternalServerError, func(c iriscontext.Context) {
				NewContext(c).Render(ecode.StatusInternalServerError, nil)
			})
		case http.StatusNotFound:
			app.OnErrorCode(http.StatusNotFound, func(c iriscontext.Context) {
				NewContext(c).Render(ecode.StatusApiNotFount, nil)
			})
		case http.StatusTooManyRequests:
			app.OnErrorCode(http.StatusTooManyRequests, func(c iriscontext.Context) {
				NewContext(c).Render(ecode.StatusTooManyRequests, nil)
			})
		}
	}
}

var (
	defaultLogValuers = make(map[string]LogValuer)
	logValuerOnce     = new(sync.Once)
)

// LogValuer ...
type LogValuer func(ctx *Context) interface{}

// RegisterLogValuers ...
func RegisterLogValuers(ms map[string]LogValuer) {
	logValuerOnce.Do(func() {
		defaultLogValuers = ms
	})
}

// logValuers ...
func logValuers() map[string]LogValuer {
	return defaultLogValuers
}
