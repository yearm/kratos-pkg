package xiris

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/gogf/gf/v2/util/gconv"
	jsoniter "github.com/json-iterator/go"
	"github.com/yearm/kratos-pkg/ecode"
	"github.com/yearm/kratos-pkg/result"
	"github.com/yearm/kratos-pkg/util/debug"
	"github.com/yearm/kratos-pkg/xerrors"
	"net/http"
)

// Render iris render
func (c *Context) Render(status ecode.Status, data interface{}, httpCode ...int) {
	r := result.NewWithCaller(status, data, result.NilCaller)
	c.json(r, httpCode...)
}

// RenderError ...
func (c *Context) RenderError(status ecode.Status, data interface{}, httpCode ...int) {
	r := result.NewWithCaller(status, data, result.NilCaller, log.LevelError)
	c.json(r, httpCode...)
}

// RenderWithMsg ...
func (c *Context) RenderWithMsg(status ecode.Status, data interface{}, msg string, httpCode ...int) {
	r := result.NewWithCaller(status, data, result.NilCaller)
	r.SetMessage(msg)
	c.json(r, httpCode...)
}

// RenderErrorWithMsg ...
func (c *Context) RenderErrorWithMsg(status ecode.Status, data interface{}, msg string, httpCode ...int) {
	r := result.NewWithCaller(status, data, result.NilCaller, log.LevelError)
	r.SetMessage(msg)
	c.json(r, httpCode...)
}

// RenderResult ...
func (c *Context) RenderResult(r *result.Result, httpCode ...int) {
	c.json(r, httpCode...)
}

// RenderResultWithLevel ...
func (c *Context) RenderResultWithLevel(r *result.Result, level log.Level, httpCode ...int) {
	if r != nil {
		r.SetLevel(level)
	}
	c.json(r, httpCode...)
}

// RenderText ...
func (c *Context) RenderText(httpCode int, text string, levels ...log.Level) {
	level := log.LevelWarn
	if len(levels) > 0 {
		level = levels[0]
	}
	r := &result.Result{Data: text}
	r.SetLevel(level)
	c.text(r, httpCode)
}

// RenderHTML ...
func (c *Context) RenderHTML(httpCode int, text string, levels ...log.Level) {
	level := log.LevelWarn
	if len(levels) > 0 {
		level = levels[0]
	}
	r := &result.Result{Data: text}
	r.SetLevel(level)
	c.html(r, httpCode)
}

// json ...
func (c *Context) json(result *result.Result, httpCode ...int) {
	c.Header("Cache-Control", "no-cache")
	c.Header("X-Request-Id", tracing.TraceID()(c).(string))
	if len(httpCode) > 0 {
		c.StatusCode(httpCode[0])
	}
	c.log(result)
	_, _ = c.JSON(result)
}

// text ...
func (c *Context) text(result *result.Result, httpCode int) {
	c.Header("Cache-Control", "no-cache")
	c.Header("X-Request-Id", tracing.TraceID()(c).(string))
	c.StatusCode(httpCode)
	c.log(result)
	_, _ = c.Text(gconv.String(result.Data))
}

// html ...
func (c *Context) html(result *result.Result, httpCode int) {
	c.Header("Cache-Control", "no-cache")
	c.Header("X-Request-Id", tracing.TraceID()(c).(string))
	c.StatusCode(httpCode)
	c.log(result)
	_, _ = c.HTML(gconv.String(result.Data))
}

// log ...
func (c *Context) log(result *result.Result) {
	if result == nil {
		return
	}
	body, _ := c.GetBody()
	urlParams, _ := jsoniter.MarshalToString(c.URLParams())
	params := map[string]interface{}{
		"method":    c.Method(),
		"path":      c.Path(),
		"param":     c.ParamsString(),
		"urlParams": urlParams,
		"body":      string(body),
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

	callers := make([]string, 0, 10)
	callers = append(callers, debug.Caller(4, 3))
	if result.Caller() != "" {
		callers = append(callers, result.Caller())
	}

	_result := map[string]interface{}{
		"status":   result.Status,
		"message":  result.Msg,
		"httpCode": c.GetStatusCode(),
	}
	if _, ok := result.Data.(string); ok {
		_result["data"] = result.Data
	}
	if err, ok := result.Data.(error); ok {
		result.Data = nil
		_result["error"] = err.Error()
		callers = append(callers, xerrors.Callers(err)...)
	}

	level := log.LevelWarn
	if result.Level().String() != "" {
		level = result.Level()
	}
	log.Context(c).Log(level,
		"@field", map[string]interface{}{
			"params":      params,
			"processTime": c.ProcessTime(),
			"result":      _result,
			"callers":     callers,
		},
	)
}
