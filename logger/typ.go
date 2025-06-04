package logger

const (
	RequestKey  = "request"
	ResponseKey = "response"
)

// RequestLog contains request-specific information
type RequestLog struct {
	Kind     string `json:"kind"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
	Header   string `json:"header"`
	Req      string `json:"req"`
	ClientIP string `json:"clientIP"`
	Extra    string `json:"extra"`
}

// ResponseLog contains response-specific information
type ResponseLog struct {
	Code        uint32          `json:"code"`
	Error       string          `json:"error,omitempty"`
	ErrorDetail *ErrorDetailLog `json:"errorDetail,omitempty"`
	Reply       string          `json:"reply"`
	Latency     int64           `json:"latency"`
}

// ErrorDetailLog contains error detail information
type ErrorDetailLog struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
	Level   string `json:"level"`
	Callers string `json:"callers"`
}
