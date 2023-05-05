package xiris

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/kataras/iris/v12"
	iriscontext "github.com/kataras/iris/v12/context"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	// validate ...
	validate = validator.New()
	trans    ut.Translator
)

// init ...
func init() {
	zhT := zh.New()
	uni := ut.New(zhT, zhT)
	trans, _ = uni.GetTranslator("zh")
	if err := zhtrans.RegisterDefaultTranslations(validate, trans); err != nil {
		logrus.Panicln("validator register translations error:", err)
	}
}

// stdContext ...
type stdContext struct {
	context.Context
}

// Context ...
type Context struct {
	iris.Context
	stdContext
}

// NewContext ...
func NewContext(ctx iris.Context) *Context {
	return &Context{Context: ctx, stdContext: stdContext{ctx.Request().Context()}}
}

// ReadJSONValid ...
func (c *Context) ReadJSONValid(outPtr interface{}) error {
	if err := c.Context.ReadJSON(outPtr); err != nil {
		return err
	}
	if err := validate.Struct(outPtr); err != nil {
		if e, ok := (err).(validator.ValidationErrors); ok {
			var errBuffer bytes.Buffer
			for k, v := range e.Translate(trans) {
				errBuffer.WriteString(k)
				errBuffer.WriteString(":")
				errBuffer.WriteString(v)
				errBuffer.WriteString(",")
			}
			errText := errBuffer.String()
			if errBuffer.String() != "" {
				errText = strings.TrimSuffix(errBuffer.String(), ",")
			}
			return errors.New(errText)
		}
		return err
	}
	return nil
}

// GetLimitAndOffset ...
func (c *Context) GetLimitAndOffset(isQueryAll bool, maxLimit ...int) (limit, offset int) {
	offset = c.URLParamIntDefault("offset", 0)
	limit = c.URLParamIntDefault("limit", 0)

	if len(maxLimit) > 0 && maxLimit[0] > 0 {
		max := maxLimit[0]
		if limit > max || limit == 0 {
			limit = max
		}
	} else {
		if !isQueryAll {
			if limit <= 0 || limit > 10 {
				limit = 10
			}
		}
	}
	return
}

// ParamsString ...
func (c *Context) ParamsString() string {
	var paramsBuff bytes.Buffer
	c.Params().Visit(func(key string, value string) {
		paramsBuff.WriteString(key)
		paramsBuff.WriteString("=")
		paramsBuff.WriteString(value)
		paramsBuff.WriteString(",")
	})
	if paramsBuff.String() != "" {
		return strings.TrimSuffix(paramsBuff.String(), ",")
	}
	return paramsBuff.String()
}

// ProcessTime ...
func (c *Context) ProcessTime() (processTime int64) {
	startAt := c.Values().GetInt64Default("startAt", -1)
	if startAt > 0 {
		processTime = time.Now().UnixMilli() - startAt
	}
	return
}

// Handle ...
func Handle(handleFunc func(ctx *Context)) iriscontext.Handler {
	return func(c iriscontext.Context) {
		handleFunc(NewContext(c))
	}
}

type irisKey struct{}

// NewIrisContext ...
func NewIrisContext(ctx context.Context, c iris.Context) context.Context {
	return context.WithValue(ctx, irisKey{}, c)
}

// FromIrisContext ...
func FromIrisContext(ctx context.Context) (c iris.Context, ok bool) {
	c, ok = ctx.Value(irisKey{}).(iris.Context)
	return
}
