package xiris

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/kataras/iris/v12"
	iriscontext "github.com/kataras/iris/v12/context"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

var (
	// validate ...
	validate *validator.Validate
	trans    ut.Translator
)

// init ...
func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("errMsg")
	})
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
// Deprecated: Use ReadJSONValidate
func (c *Context) ReadJSONValid(outPtr interface{}) error {
	if err := c.ReadJSON(outPtr); err != nil {
		return err
	}
	_, err := c.validateStruct(outPtr)
	return err
}

// ReadJSONValidate true: error is errMsg tag message
func (c *Context) ReadJSONValidate(outPtr interface{}) (customized bool, err error) {
	if err = c.ReadJSON(outPtr); err != nil {
		return
	}
	return c.validateStruct(outPtr)
}

// ReadQueryValidate true: error is errMsg tag message
func (c *Context) ReadQueryValidate(ptr interface{}) (customized bool, err error) {
	if err = c.ReadQuery(ptr); err != nil {
		return
	}
	return c.validateStruct(ptr)
}

// ValidateStruct ...
func (c *Context) validateStruct(ptr interface{}) (customized bool, err error) {
	if err = validate.Struct(ptr); err != nil {
		if fieldErrors, ok := (err).(validator.ValidationErrors); ok {
			for _, fieldError := range fieldErrors {
				var errMsg string
				translateValue := fieldError.Translate(trans)
				// NOTE: Field() 和 StructField() 不相等说明取到了 errMsg tag 值
				if fieldError.Field() != fieldError.StructField() {
					// NOTE: 翻译时取的值是 Field()，由于前面 RegisterTagNameFunc 取的是 errMsg tag 对应的值，所以这里翻译后要替换成 StructField()
					translateValue = strings.Replace(translateValue, fieldError.Field(), fieldError.StructField(), 1)
					errMsg = fieldError.Field()
				}
				if errMsg != "" {
					return true, errors.New(errMsg)
				}
				return false, errors.New(fmt.Sprintf("%s:%s", fieldError.StructNamespace(), translateValue))
			}
		}
		return
	}
	return
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
