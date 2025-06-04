package env

import (
	"os"
	"sync"
)

var (
	once             sync.Once
	_serviceID       string
	_serviceName     string
	_serviceVersion  string
	_serviceMetadata map[string]string
)

// Init initialize the global service parameters for identifying the service identity and metadata.
func Init(serviceName, serviceVersion string, serviceMetadata map[string]string) {
	once.Do(func() {
		_serviceID, _ = os.Hostname()
		_serviceName = serviceName
		_serviceVersion = serviceVersion
		_serviceMetadata = serviceMetadata
	})
}

// GetServiceID returns the service id.
func GetServiceID() string {
	return _serviceID
}

// GetServiceName returns the service name.
func GetServiceName() string {
	return _serviceName
}

// GetServiceVersion returns the service version.
func GetServiceVersion() string {
	return _serviceVersion
}

// GetServiceMetadata return the service metadata.
func GetServiceMetadata() map[string]string {
	return _serviceMetadata
}
