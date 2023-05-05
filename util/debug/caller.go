package debug

import (
	"fmt"
	"runtime"
	"strings"
)

// Caller ...
func Caller(skip, sepN int) string {
	_, file, line, _ := runtime.Caller(skip)
	var (
		_file = file
		idx   int
	)
	for i := 0; i < sepN; i++ {
		idx = strings.LastIndexByte(_file, '/')
		if idx == -1 {
			return fmt.Sprintf("%s:%d", _file[idx+1:], line)
		}
		_file = _file[:idx]
	}
	return fmt.Sprintf("%s:%d", file[idx+1:], line)
}

// CallerByFrame ...
func CallerByFrame(frame *runtime.Frame, sepN int) string {
	var (
		_file = frame.File
		idx   int
	)
	for i := 0; i < sepN; i++ {
		idx = strings.LastIndexByte(_file, '/')
		if idx == -1 {
			return fmt.Sprintf("%s:%d", _file[idx+1:], frame.Line)
		}
		_file = _file[:idx]
	}
	return fmt.Sprintf("%s:%d", frame.File[idx+1:], frame.Line)
}
