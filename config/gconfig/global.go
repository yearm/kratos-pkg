package gconfig

import (
	"github.com/go-kratos/kratos/v2/config"
	"sync"
)

// globalConfig holds the application's global configuration state.
type globalConfig struct {
	once sync.Once
	config.Config
}

// global is the singleton instance maintaining configuration state.
var global = &globalConfig{
	once:   sync.Once{},
	Config: config.New(),
}

// SetConfig initializes the global configuration (single-shot operation).
func (g *globalConfig) SetConfig(c config.Config) {
	g.once.Do(func() {
		g.Config = c
	})
}

// GetConfig retrieves the initialized configuration instance.
func (g *globalConfig) GetConfig() config.Config {
	return g.Config
}

// SetConfig sets the global configuration singleton.
func SetConfig(c config.Config) {
	global.SetConfig(c)
}

// GetConfig returns the active configuration instance.
func GetConfig() config.Config {
	return global.GetConfig()
}

// Value retrieves a configuration value by dot-delimited key path.
func Value(key string) config.Value {
	return GetConfig().Value(key)
}

// Scan unmarshal entire configuration into target struct.
func Scan(v any) error {
	return GetConfig().Scan(v)
}

// Watch registers observer for configuration changes.
func Watch(key string, o config.Observer) error {
	return GetConfig().Watch(key, o)
}
