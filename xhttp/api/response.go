package api

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/yearm/kratos-pkg/ecodes"
	"github.com/yearm/kratos-pkg/errors"
	"github.com/yearm/kratos-pkg/utils/debug"
	"github.com/yearm/kratos-pkg/utils/gjson"
	"github.com/yearm/kratos-pkg/xgrpc/status"
	"net/http"
	"strings"
)

// RenderType defines the response content rendering type.
type RenderType uint8

const (
	// JSON render full JSON response with code, message and data.
	JSON RenderType = iota

	// DataOnlyJSON render only data portion as JSON.
	DataOnlyJSON

	// String render response data as String content.
	String
)

// Response represents a standardized API response structure.
type Response struct {
	Code    ecodes.Code `json:"code"`
	Message string      `json:"message"`
	Data    any         `json:"data,omitempty"`

	ctx        context.Context
	httpCode   int
	err        error
	callers    string
	level      log.Level
	renderType RenderType
}

// NewSuccess creates a successful API response.
func NewSuccess(ctx context.Context, data any) *Response {
	r := &Response{
		Code:       ecodes.OK,
		Message:    ecodes.OK.Message(ctx),
		Data:       data,
		ctx:        ctx,
		httpCode:   http.StatusOK,
		level:      log.LevelInfo,
		renderType: JSON,
	}
	return r
}

// NewError creates an error API response with stack trace.
func NewError(ctx context.Context, code ecodes.Code, err error) *Response {
	err = errors.WrapDepth(err, 3)
	r := &Response{
		Code:       code,
		Message:    code.Message(ctx),
		ctx:        ctx,
		httpCode:   http.StatusOK,
		err:        err,
		callers:    errors.Callers(err),
		level:      log.LevelError,
		renderType: JSON,
	}
	return r
}

// FromError constructs a structured Response from an error.
func FromError(ctx context.Context, err error) *Response {
	err = errors.WrapDepth(err, 3)
	r := &Response{
		ctx:        ctx,
		httpCode:   http.StatusOK,
		err:        err,
		callers:    errors.Callers(err),
		level:      log.LevelError,
		renderType: JSON,
	}
	st, detail, ok := status.FromError(err)
	if !ok {
		r.Code = ecodes.UnknownError
		r.Message = ecodes.UnknownError.Message(ctx)
		return r
	}

	if detail == nil {
		code := ecodes.FromGRPCCode(st.Code())
		r.Code = code
		r.Message = code.Message(ctx)
		return r
	}

	errorDetail := status.ToErrorDetail(detail)
	r.Code = errorDetail.Code
	r.Message = errorDetail.Message
	r.level = errorDetail.Level
	return r
}

// RenderString return the string based on the render type.
func (r *Response) RenderString() string {
	switch r.renderType {
	case JSON:
		return gjson.MustMarshalToString(r)
	case DataOnlyJSON:
		return gjson.MustMarshalToString(r.Data)
	case String:
		return fmt.Sprint(r.Data)
	}
	return ""
}

// WithMessage sets custom message for the response.
func (r *Response) WithMessage(message string) *Response {
	r.Message = message
	return r
}

// WithData attaches payload data to the response.
func (r *Response) WithData(data any) *Response {
	r.Data = data
	return r
}

// WithHTTPCode overrides default HTTP status code.
func (r *Response) WithHTTPCode(code int) *Response {
	r.httpCode = code
	return r
}

// WithError attaches original error to the response.
func (r *Response) WithError(err error) *Response {
	r.err = err
	return r
}

// WithLevel sets logging level for the response.
func (r *Response) WithLevel(level log.Level) *Response {
	r.level = level
	return r
}

// WithRenderType specifies response content rendering format.
func (r *Response) WithRenderType(renderType RenderType) *Response {
	r.renderType = renderType
	return r
}

// WrapCallers appends current caller information to the stack trace.
func (r *Response) WrapCallers(depth ...int) *Response {
	d := 2
	if len(depth) > 0 {
		d = depth[0]
	}
	if r.callers == "[]" {
		r.callers = strings.Replace(r.callers, "[", fmt.Sprintf(`["%s"`, debug.Caller(d)), 1)
	} else {
		r.callers = strings.Replace(r.callers, "[", fmt.Sprintf(`["%s", `, debug.Caller(d)), 1)
	}
	return r
}

// GetContext return the HTTP status code of the response.
func (r *Response) GetContext() context.Context {
	return r.ctx
}

// GetHTTPCode return the HTTP status code of the response.
func (r *Response) GetHTTPCode() int {
	return r.httpCode
}

// GetError return the error of the response.
func (r *Response) GetError() error {
	return r.err
}

// GetCallers return the callers of the response.
func (r *Response) GetCallers() string {
	return r.callers
}

// GetLevel return the level of the response.
func (r *Response) GetLevel() log.Level {
	return r.level
}

// GetRenderType return the render type of the response.
func (r *Response) GetRenderType() RenderType {
	return r.renderType
}
