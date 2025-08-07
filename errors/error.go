package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yearm/kratos-pkg/utils/debug"
)

// withCaller is an error wrapper type that captures error context.
type withCaller struct {
	err     error
	message string
	caller  string
}

// Error implements error interface, returns formatted message.
func (w *withCaller) Error() string {
	var errString string
	if w.err != nil {
		errString = w.err.Error()
	}

	if w.message == "" {
		return errString
	}
	return w.message + ": " + errString
}

// Cause returns the original wrapped error.
func (w *withCaller) Cause() error {
	return w.err
}

// Unwrap supports error chain unwrapping.
func (w *withCaller) Unwrap() error {
	return w.err
}

// New creates a base error with caller location information.
func New(text string) error {
	return &withCaller{
		err:    errors.New(text),
		caller: debug.Caller(2),
	}
}

// Errorf creates a formatted error with caller location.
func Errorf(format string, a ...any) error {
	return &withCaller{
		err:    fmt.Errorf(format, a...),
		caller: debug.Caller(2),
	}
}

// Wrap adds contextual message to an existing error.
func Wrap(err error, msg ...string) error {
	var message string
	if len(msg) > 0 {
		message = msg[0]
	}
	return &withCaller{
		err:     err,
		message: message,
		caller:  debug.Caller(2),
	}
}

// Wrapf adds formatted contextual message to an existing error.
func Wrapf(err error, format string, args ...any) error {
	return &withCaller{
		err:     err,
		message: fmt.Sprintf(format, args...),
		caller:  debug.Caller(2),
	}
}

// WrapDepth adds contextual message to an existing error, depth: call stack depth.
func WrapDepth(err error, depth int) error {
	return &withCaller{
		err:    err,
		caller: debug.Caller(depth),
	}
}

// Cause retrieves the root cause error from the chain.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		c, ok := err.(causer)
		if !ok {
			break
		}
		err = c.Cause()
	}
	return err
}

// Callers extracts call location chain from error.
func Callers(err error) string {
	type unwraper interface {
		Unwrap() error
	}

	var (
		buf  strings.Builder
		flag bool
	)
	buf.Grow(256)
	buf.WriteString("[")
	for err != nil {
		if wc := new(withCaller); errors.As(err, &wc) {
			if flag {
				buf.WriteString(", ")
			}
			buf.WriteString(fmt.Sprintf(`"%v"`, wc.caller))
			flag = true
		}
		u, ok := err.(unwraper)
		if !ok {
			break
		}
		err = u.Unwrap()
	}
	buf.WriteString("]")
	return buf.String()
}

// Standard library compatibility functions:

// Is checks if target error exists in the error chain (stdlib-compatible).
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As attempts to convert error to target type (stdlib-compatible).
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Unwrap returns next error in the chain (stdlib-compatible).
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Join combines multiple errors into a single error (Go 1.20+ compatible).
func Join(errs ...error) error {
	return errors.Join(errs...)
}
