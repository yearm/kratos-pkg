package ecode

import "golang.org/x/exp/maps"

// Status ...
type Status string

// String ...
func (s Status) String() string {
	return string(s)
}

// Message ...
func (s Status) Message() string {
	return statusMap[s]
}

var (
	statusMap = map[Status]string{}
)

func init() {
	maps.Copy(statusMap, commonStatus)
	maps.Copy(statusMap, resourceStatus)
	maps.Copy(statusMap, authStatus)
	maps.Copy(statusMap, paymentStatus)
}
