package ecode

import "google.golang.org/grpc/codes"

const (
	StatusOk                     Status = "OK"                      // 请求成功
	StatusUnknownError           Status = "UNKNOWN_ERROR"           // 未知错误
	StatusApiNotFount            Status = "API_NOT_FOUND"           // 请求的API不存在
	StatusNotFound               Status = "NOT_FOUND"               // 资源不存在
	StatusRecordNotFound         Status = "RECORD_NOT_FOUND"        // 记录不存在
	StatusInvalidRequest         Status = "INVALID_REQUEST"         // 请求参数错误
	StatusRequestTimeout         Status = "REQUEST_TIMEOUT"         // 请求超时
	StatusTooManyRequests        Status = "TOO_MANY_REQUESTS"       // 请求过于频繁
	StatusTemporarilyUnavailable Status = "TEMPORARILY_UNAVAILABLE" // 服务暂时无法访问
	StatusInternalServerError    Status = "INTERNAL_SERVER_ERROR"   // 服务内部异常
	StatusExternalServerError    Status = "EXTERNAL_SERVER_ERROR"   // 外部服务异常
)

var commonStatus = map[Status]string{
	StatusOk:                     "",
	StatusUnknownError:           "未知错误",
	StatusApiNotFount:            "请求的API不存在",
	StatusNotFound:               "资源不存在",
	StatusRecordNotFound:         "记录不存在",
	StatusInvalidRequest:         "请求参数错误",
	StatusRequestTimeout:         "请求超时",
	StatusTooManyRequests:        "请求过于频繁",
	StatusTemporarilyUnavailable: "服务暂时无法访问",
	StatusInternalServerError:    "服务内部异常",
	StatusExternalServerError:    "外部服务异常",
}

const (
	StatusInvalidResource Status = "INVALID_RESOURCE" // 资源无效
	StatusLocked          Status = "LOCKED"           // 资源已被锁定
	StatusGone            Status = "GONE"             // 资源已失效
	StatusCancelled       Status = "CANCELLED"        // 操作已被取消
	StatusCompleted       Status = "COMPLETED"        // 操作已完成
)

var resourceStatus = map[Status]string{
	StatusInvalidResource: "资源无效",
	StatusLocked:          "资源已被锁定",
	StatusGone:            "资源已失效",
	StatusCancelled:       "操作已被取消",
	StatusCompleted:       "操作已完成",
}

const (
	StatusUnauthorized     Status = "UNAUTHORIZED"      // 会话信息无效
	StatusUnauthorizedUser Status = "UNAUTHORIZED_USER" // 当前会话未登录
	StatusAccessDenied     Status = "ACCESS_DENIED"     // 拒绝访问
)

var authStatus = map[Status]string{
	StatusUnauthorized:     "会话信息无效",
	StatusUnauthorizedUser: "当前会话未登录",
	StatusAccessDenied:     "拒绝访问",
}

const (
	StatusPaymentValueNotEnough Status = "PAYMENT_VALUE_NOT_ENOUGH" // 余额不足
	StatusPaymentLocked         Status = "PAYMENT_LOCKED"           // 支付方式已被锁定
	StatusPaymentUnauthorized   Status = "PAYMENT_UNAUTHORIZED"     // 支付方式未被授权使用
	StatusAlipayUnauthorized    Status = "ALIPAY_UNAUTHORIZED"      // 未授权支付宝
	StatusWechatUnauthorized    Status = "WECHAT_UNAUTHORIZED"      // 未授权微信
	StatusInvalidPayment        Status = "INVALID_PAYMENT"          // 无效的支付方式
	StatusUnsupportedPayment    Status = "UNSUPPORTED_PAYMENT"      // 不支持的支付方式
	StatusInvalidCallbackURL    Status = "INVALID_CALLBACK_URL"     // 无效的回调地址
	StatusPaymentUnknownError   Status = "PAYMENT_UNKNOWN_ERROR"    // 支付未知错误
	StatusPointValueNotEnough   Status = "POINT_VALUE_NOT_ENOUGH"   // 积分余额不足
)

var paymentStatus = map[Status]string{
	StatusPaymentValueNotEnough: "余额不足",
	StatusPaymentLocked:         "支付方式已被锁定",
	StatusPaymentUnauthorized:   "支付方式未授权",
	StatusAlipayUnauthorized:    "未授权支付宝",
	StatusWechatUnauthorized:    "未授权微信",
	StatusInvalidPayment:        "无效的支付方式",
	StatusUnsupportedPayment:    "不支持的支付方式",
	StatusInvalidCallbackURL:    "无效的回调地址",
	StatusPaymentUnknownError:   "支付未知错误",
	StatusPointValueNotEnough:   "积分余额不足",
}

const (
	// RPCBusinessError ...
	RPCBusinessError codes.Code = 101
)
