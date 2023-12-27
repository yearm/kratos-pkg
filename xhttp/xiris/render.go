package xiris

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/gogf/gf/v2/util/gconv"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
	"github.com/yearm/kratos-pkg/ecode"
	"github.com/yearm/kratos-pkg/result"
	"github.com/yearm/kratos-pkg/util/debug"
	"github.com/yearm/kratos-pkg/xerrors"
	"net/http"
)

// Render iris render
func (c *Context) Render(status ecode.Status, data interface{}, httpCode ...int) {
	caller := debug.Caller(2, 3)
	opts := []result.Option{result.Caller(caller)}
	if len(httpCode) > 0 {
		opts = append(opts, result.HttpCode(httpCode[0]))
	}
	res := result.New(status, data, opts...)
	c.render(res)
}

// RenderError ...
func (c *Context) RenderError(status ecode.Status, data interface{}, httpCode ...int) {
	caller := debug.Caller(2, 3)
	opts := []result.Option{
		result.Caller(caller),
		result.Level(log.LevelError),
	}
	if len(httpCode) > 0 {
		opts = append(opts, result.HttpCode(httpCode[0]))
	}
	res := result.New(status, data, opts...)
	c.render(res)
}

// RenderWithMsg ...
func (c *Context) RenderWithMsg(status ecode.Status, data interface{}, msg string, httpCode ...int) {
	caller := debug.Caller(2, 3)
	opts := []result.Option{
		result.Caller(caller),
		result.Message(msg),
	}
	if len(httpCode) > 0 {
		opts = append(opts, result.HttpCode(httpCode[0]))
	}
	res := result.New(status, data, opts...)
	c.render(res)
}

// RenderErrorWithMsg ...
func (c *Context) RenderErrorWithMsg(status ecode.Status, data interface{}, msg string, httpCode ...int) {
	caller := debug.Caller(2, 3)
	opts := []result.Option{
		result.Caller(caller),
		result.Message(msg),
		result.Level(log.LevelError),
	}
	if len(httpCode) > 0 {
		opts = append(opts, result.HttpCode(httpCode[0]))
	}
	res := result.New(status, data, opts...)
	c.render(res)
}

// RenderText ...
func (c *Context) RenderText(httpCode int, text string, levels ...log.Level) {
	caller := debug.Caller(2, 3)
	opts := []result.Option{
		result.Caller(caller),
		result.HttpCode(httpCode),
		result.RenderType(result.TEXT),
	}
	if len(levels) > 0 {
		opts = append(opts, result.Level(levels[0]))
	}
	res := result.New(ecode.StatusOk, text, opts...)
	c.render(res)
}

// RenderHTML ...
func (c *Context) RenderHTML(httpCode int, text string, levels ...log.Level) {
	caller := debug.Caller(2, 3)
	opts := []result.Option{
		result.Caller(caller),
		result.HttpCode(httpCode),
		result.RenderType(result.HTML),
	}
	if len(levels) > 0 {
		opts = append(opts, result.Level(levels[0]))
	}
	res := result.New(ecode.StatusOk, text, opts...)
	c.render(res)
}

// RenderResult ...
func (c *Context) RenderResult(res *result.Result) {
	c.render(res, debug.Caller(2, 3))
}

func (c *Context) render(res *result.Result, beforeCaller ...string) {
	c.Header("Cache-Control", "no-cache")
	c.Header("X-Request-Id", tracing.TraceID()(c).(string))
	if res.HttpCode() > 0 {
		c.StatusCode(res.HttpCode())
	}
	c.log(res, beforeCaller...)
	switch res.RenderType() {
	case result.JSON:
		_, _ = c.JSON(res)
	case result.HTML:
		_, _ = c.HTML(gconv.String(res.Data))
	case result.TEXT:
		_, _ = c.Text(gconv.String(res.Data))
	}
}

// log ...
func (c *Context) log(result *result.Result, beforeCaller ...string) {
	if result == nil {
		return
	}
	body, _ := c.GetBody()
	pathParam, _ := jsoniter.MarshalToString(c.PathParams())
	queryParam, _ := jsoniter.MarshalToString(c.URLParams())
	params := map[string]interface{}{
		"clientIP":    c.RemoteAddr(),
		"method":      c.Method(),
		"path":        c.Path(),
		"handlerName": c.HandlerName(),
		"pathParam":   pathParam,
		"queryParam":  queryParam,
		"body":        string(body),
	}
	if c.Method() == http.MethodPost && len(c.Request().PostForm) > 0 {
		params["postForm"] = c.Request().PostForm
	}
	if cRouter := c.GetCurrentRoute(); cRouter != nil {
		params["path"] = cRouter.Path()
	}
	for key, valuer := range logValuers() {
		params[key] = valuer(c)
	}

	_result := map[string]interface{}{
		"status":   result.Status,
		"message":  result.Msg,
		"httpCode": c.GetStatusCode(),
	}
	if _, ok := result.Data.(string); ok {
		_result["data"] = result.Data
	}

	callers := make([]string, 0, 10)
	if len(beforeCaller) > 0 {
		callers = append(callers, beforeCaller[0])
	}
	if result.Caller() != "" {
		callers = append(callers, result.Caller())
	}
	if err, ok := result.Data.(error); ok {
		result.Data = nil
		_result["error"] = err.Error()
		callers = append(callers, xerrors.Callers(err)...)
	}

	level := lo.If(result.Level().String() != "", result.Level()).Else(log.LevelWarn)
	log.Context(c).Log(level, "field", map[string]interface{}{
		"params":      params,
		"processTime": c.ProcessTime(),
		"result":      _result,
		"callers":     callers,
	})
}
