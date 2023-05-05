package ecode

import "google.golang.org/grpc/codes"

const (
	StatusOk                     Status = "OK"
	StatusUnknownError           Status = "UNKNOWN_ERROR"
	StatusApiNotFount            Status = "API_NOT_FOUND"
	StatusNotFound               Status = "NOT_FOUND"
	StatusInvalidRequest         Status = "INVALID_REQUEST"
	StatusRequestTimeout         Status = "REQUEST_TIMEOUT"
	StatusTooManyRequests        Status = "TOO_MANY_REQUESTS"
	StatusTemporarilyUnavailable Status = "TEMPORARILY_UNAVAILABLE"
	StatusInternalServerError    Status = "INTERNAL_SERVER_ERROR"
	StatusExternalServerError    Status = "EXTERNAL_SERVER_ERROR"

	StatusInvalidResource Status = "INVALID_RESOURCE"
	StatusLocked          Status = "LOCKED"
	StatusGone            Status = "GONE"
	StatusCancelled       Status = "CANCELLED"
	StatusCompleted       Status = "COMPLETED"

	StatusUnauthorized     Status = "UNAUTHORIZED"
	StatusUnauthorizedUser Status = "UNAUTHORIZED_USER"
	StatusAccessDenied     Status = "ACCESS_DENIED"
)

var StatusMap = map[Status]string{
	StatusOk:                     "",
	StatusUnknownError:           "未知错误",
	StatusApiNotFount:            "请求的API不存在",
	StatusNotFound:               "资源不存在",
	StatusInvalidRequest:         "请求参数错误",
	StatusRequestTimeout:         "请求超时",
	StatusTooManyRequests:        "请求过于频繁",
	StatusTemporarilyUnavailable: "服务暂时无法访问",
	StatusInternalServerError:    "服务内部异常",
	StatusExternalServerError:    "外部服务异常",

	StatusInvalidResource: "资源无效",
	StatusLocked:          "资源已被锁定",
	StatusGone:            "资源已失效",
	StatusCancelled:       "操作已被取消",
	StatusCompleted:       "操作已完成",

	StatusUnauthorized:     "会话信息无效",
	StatusUnauthorizedUser: "当前会话未登录",
	StatusAccessDenied:     "拒绝访问",
}

// RPCBusinessError ...
const RPCBusinessError codes.Code = 101
