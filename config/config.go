package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/samber/lo"
	"github.com/yearm/kratos-pkg/config/gconfig"
	"github.com/yearm/kratos-pkg/errors"
)

// Load initializes and merges configurations from multiple sources.
func Load(cs []config.Source, opts ...config.Option) (func(), error) {
	cs = lo.Filter(cs, func(item config.Source, index int) bool {
		return item != nil
	})
	options := []config.Option{config.WithSource(cs...)}
	c := config.New(append(options, opts...)...)
	if err := c.Load(); err != nil {
		return nil, errors.Wrap(err, "config.Load failed")
	}
	gconfig.SetConfig(c)
	return func() { _ = c.Close() }, nil
}
