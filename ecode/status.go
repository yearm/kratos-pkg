package ecode

// Status ...
type Status string

// String ...
func (s Status) String() string {
	return string(s)
}

// Message ...
func (s Status) Message() string {
	return StatusMap[s]
}
