package file

import (
	"os"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
)

// NewConfigSource creates a file config source.
func NewConfigSource(path string) config.Source {
	return file.NewSource(path)
}

// NewConfigSourceIfExists creates a file config source if exists.
func NewConfigSourceIfExists(path string) config.Source {
	_, err := os.Stat(path)
	if err != nil {
		log.Warnf("config file[%s] is not exists", path)
		return nil
	}
	return file.NewSource(path)
}
