package gconfig

import (
	"github.com/samber/lo"
	"github.com/yearm/kratos-pkg/errors"
	"sync"
)

const (
	Local       Mode = "local"
	Development Mode = "dev"
	Test        Mode = "test"
	Staging     Mode = "staging"
	Production  Mode = "prod"
)

type Mode string

func (m Mode) String() string {
	return string(m)
}

func (m Mode) isValid() bool {
	return lo.Contains([]Mode{Local, Development, Test, Staging, Production}, m)
}

var (
	// defaultModeKey default key for mode.
	defaultModeKey = "mode"
	modeKeyOnce    sync.Once
)

// SetModeKey customizes the global key for mode.
func SetModeKey(key string) {
	modeKeyOnce.Do(func() {
		defaultModeKey = key
	})
}

// GetMode retrieves mode configuration from global settings.
func GetMode() (Mode, error) {
	var m Mode
	err := Value(defaultModeKey).Scan(&m)
	if err != nil {
		return "", errors.Wrapf(err, "config.Scan[%v] faild", defaultModeKey)
	}
	if !m.isValid() {
		return "", errors.Errorf("mode config[%v] is invalid", defaultModeKey)
	}
	return m, nil
}

// IsLocalMode checks if the current runtime mode is a local environment.
func IsLocalMode() (bool, error) {
	m, err := GetMode()
	if err != nil {
		return false, errors.Wrap(err, "GetMode failed")
	}
	return m == Local, nil
}

// IsDevelopmentMode checks if the current runtime mode is a development environment.
func IsDevelopmentMode() (bool, error) {
	m, err := GetMode()
	if err != nil {
		return false, errors.Wrap(err, "GetMode failed")
	}
	return lo.Contains([]Mode{Local, Development, Test}, m), nil
}

// IsProductionMode checks if the current runtime mode is a production environment.
func IsProductionMode() (bool, error) {
	m, err := GetMode()
	if err != nil {
		return false, errors.Wrap(err, "GetMode failed")
	}
	return lo.Contains([]Mode{Staging, Production}, m), nil
}
