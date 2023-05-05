package xiris

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	iriscontext "github.com/kataras/iris/v12/context"
	"github.com/yearm/kratos-pkg/ecode"
	"net/http"
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

// LogParam ...
type LogParam func(ctx *Context) string

// RegisterRenderLogParams ...
func RegisterRenderLogParams(app *iris.Application, ms map[string]LogParam) {
	for key, val := range ms {
		app.Configure(iris.WithOtherValue(key, val))
	}
}

// LogParamHeaders ...
func LogParamHeaders(names []string) LogParam {
	return func(ctx *Context) string {
		hsMap := make(map[string]interface{})
		for _, name := range names {
			hsMap[name] = ctx.GetHeader(name)
		}
		hs, _ := jsoniter.MarshalToString(hsMap)
		return hs
	}
}

// logParams ...
func logParams(app iriscontext.Application) map[string]LogParam {
	ms := make(map[string]LogParam)
	others := app.ConfigurationReadOnly().GetOther()
	for key, val := range others {
		if param, ok := val.(LogParam); ok {
			ms[key] = param
		}
	}
	return ms
}
