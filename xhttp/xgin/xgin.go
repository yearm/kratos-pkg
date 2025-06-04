package xgin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yearm/kratos-pkg/xhttp/api"
)

// responseKey indicates a default response key.
const responseKey = "_kratos-pkg/xhttp/xgin/responsekey"

type ginKey struct{}

// NewGinContext returns a new Context that carries gin.Context value.
func NewGinContext(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ginKey{}, c)
}

// FromGinContext returns the gin.Context value stored in ctx, if any.
func FromGinContext(ctx context.Context) (c *gin.Context, ok bool) {
	c, ok = ctx.Value(ginKey{}).(*gin.Context)
	return
}

// Render formats and sends the response based on its render type.
func Render(ctx *gin.Context, resp *api.Response) {
	render(ctx, resp, false)
}

// RenderWithAbort formats and sends the response based on its render type.
func RenderWithAbort(ctx *gin.Context, resp *api.Response) {
	render(ctx, resp, true)
}

// render is the internal implementation handling common response rendering logic.
func render(ctx *gin.Context, resp *api.Response, abort bool) {
	resp = resp.WrapCallers(4)
	ctx.Set(responseKey, resp)

	if abort {
		ctx.Abort()
	}
	switch resp.GetRenderType() {
	case api.JSON:
		ctx.JSON(resp.GetHTTPCode(), resp)
	case api.DataOnlyJSON:
		ctx.JSON(resp.GetHTTPCode(), resp.Data)
	case api.String:
		ctx.String(resp.GetHTTPCode(), fmt.Sprint(resp.Data))
	}
}
