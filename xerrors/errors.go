package xerrors

import (
	"fmt"
	"github.com/yearm/kratos-pkg/util/debug"
)

// New ...
func New(message string) error {
	return &fundamental{
		msg:    message,
		caller: caller(),
	}
}

// Errorf ...
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:    fmt.Sprintf(format, args...),
		caller: caller(),
	}
}

type fundamental struct {
	msg    string
	caller string
}

func (f *fundamental) Error() string { return f.msg }

// Wrap ...
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withCaller{
		cause:  err,
		caller: caller(),
	}
}

// Wrapf ...
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &withCaller{
		cause:  err,
		caller: caller(),
	}
}

// WithCaller ...
func WithCaller(err error) error {
	if err == nil {
		return nil
	}
	return &withCaller{
		cause:  err,
		caller: caller(),
	}
}

type withCaller struct {
	cause  error
	caller string
}

func (w *withCaller) Error() string { return w.cause.Error() }
func (w *withCaller) Cause() error  { return w.cause }
func (w *withCaller) Unwrap() error { return w.cause }

// WithMessage ...
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

// WithMessagef ...
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withMessage) Cause() error  { return w.cause }
func (w *withMessage) Unwrap() error { return w.cause }

// caller ...
func caller() string {
	return debug.Caller(3, 3)
}

// Cause ...
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// Callers ...
func Callers(err error) []string {
	if err == nil {
		return []string{}
	}
	type unwrap interface {
		Unwrap() error
	}
	callers := make([]string, 0, 10)
	for err != nil {
		switch e := err.(type) {
		case *fundamental:
			callers = append(callers, e.caller)
		case *withCaller:
			callers = append(callers, e.caller)
		}
		unwrap, ok := err.(unwrap)
		if !ok {
			break
		}
		err = unwrap.Unwrap()
	}
	return callers
}
