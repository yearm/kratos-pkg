package debug

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

var (
	// defaultSepN controls the number of directory separators preserved.
	defaultSepN = 3

	once sync.Once
)

// SetDefaultSepN sets the global path truncation depth for Caller() output formatting.
// Only the first call modifies defaultSepN (initial default: 3), Subsequent calls are no-ops.
func SetDefaultSepN(sepN int) {
	once.Do(func() {
		defaultSepN = sepN
	})
}

// Caller retrieves the caller's file path and line number from the runtime stack,
// then truncates the path to retain only the last N directory segments (controlled by defaultSepN).
func Caller(depth int) string {
	_, file, line, _ := runtime.Caller(depth)

	var idx int
	tmpFile := file
	for i := 0; i < defaultSepN; i++ {
		idx = strings.LastIndexByte(tmpFile, '/')
		if idx == -1 {
			return fmt.Sprintf("%s:%d", tmpFile[idx+1:], line)
		}
		tmpFile = tmpFile[:idx]
	}
	return fmt.Sprintf("%s:%d", file[idx+1:], line)
}
