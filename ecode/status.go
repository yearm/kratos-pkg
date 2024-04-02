package ecode

import (
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/exp/maps"
)

type (
	// Status ...
	Status string
	// statusResult ...
	statusResult struct {
		level   log.Level
		message string
	}
)

// String ...
func (s Status) String() string {
	return string(s)
}

// Message ...
func (s Status) Message() string {
	return statusMap[s].message
}

// Level ...
func (s Status) Level() log.Level {
	return statusMap[s].level
}

var (
	statusMap = map[Status]statusResult{}
)

func init() {
	maps.Copy(statusMap, commonStatus)
	maps.Copy(statusMap, resourceStatus)
	maps.Copy(statusMap, authStatus)
	maps.Copy(statusMap, paymentStatus)
}
