package ecode

import (
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
)

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

var commonStatus = map[Status]statusResult{
	StatusOk:                     {level: log.LevelInfo, message: ""},
	StatusUnknownError:           {level: log.LevelError, message: "未知错误"},
	StatusApiNotFount:            {level: log.LevelWarn, message: "请求的API不存在"},
	StatusNotFound:               {level: log.LevelWarn, message: "资源不存在"},
	StatusRecordNotFound:         {level: log.LevelWarn, message: "记录不存在"},
	StatusInvalidRequest:         {level: log.LevelWarn, message: "请求参数错误"},
	StatusRequestTimeout:         {level: log.LevelWarn, message: "请求超时"},
	StatusTooManyRequests:        {level: log.LevelWarn, message: "请求过于频繁"},
	StatusTemporarilyUnavailable: {level: log.LevelError, message: "服务暂时无法访问"},
	StatusInternalServerError:    {level: log.LevelError, message: "服务内部异常"},
	StatusExternalServerError:    {level: log.LevelError, message: "外部服务异常"},
}

const (
	StatusInvalidResource Status = "INVALID_RESOURCE" // 资源无效
	StatusLocked          Status = "LOCKED"           // 资源已被锁定
	StatusGone            Status = "GONE"             // 资源已失效
	StatusCancelled       Status = "CANCELLED"        // 操作已取消
	StatusCompleted       Status = "COMPLETED"        // 操作已完成
)

var resourceStatus = map[Status]statusResult{
	StatusInvalidResource: {level: log.LevelWarn, message: "资源无效"},
	StatusLocked:          {level: log.LevelWarn, message: "资源已被锁定"},
	StatusGone:            {level: log.LevelWarn, message: "资源已失效"},
	StatusCancelled:       {level: log.LevelWarn, message: "操作已取消"},
	StatusCompleted:       {level: log.LevelWarn, message: "操作已完成"},
}

const (
	StatusUnauthorized     Status = "UNAUTHORIZED"      // 会话信息无效
	StatusUnauthorizedUser Status = "UNAUTHORIZED_USER" // 当前会话未登录
	StatusAccessDenied     Status = "ACCESS_DENIED"     // 拒绝访问
)

var authStatus = map[Status]statusResult{
	StatusUnauthorized:     {level: log.LevelInfo, message: "会话信息无效"},
	StatusUnauthorizedUser: {level: log.LevelInfo, message: "当前会话未登录"},
	StatusAccessDenied:     {level: log.LevelWarn, message: "拒绝访问"},
}

const (
	StatusPaymentValueNotEnough Status = "PAYMENT_VALUE_NOT_ENOUGH" // 余额不足
	StatusPayTimeout            Status = "PAY_TIMEOUT"              // 支付超时
	StatusPaymentLocked         Status = "PAYMENT_LOCKED"           // 支付方式已被锁定
	StatusPaymentUnauthorized   Status = "PAYMENT_UNAUTHORIZED"     // 支付方式未被授权使用
	StatusAlipayUnauthorized    Status = "ALIPAY_UNAUTHORIZED"      // 未授权支付宝
	StatusWechatUnauthorized    Status = "WECHAT_UNAUTHORIZED"      // 未授权微信
	StatusInvalidPayment        Status = "INVALID_PAYMENT"          // 无效的支付方式
	StatusUnsupportedPayment    Status = "UNSUPPORTED_PAYMENT"      // 不支持的支付方式
	StatusInvalidCallbackURL    Status = "INVALID_CALLBACK_URL"     // 无效的回调地址
	StatusExternalPaymentError  Status = "EXTERNAL_PAYMENT_ERROR"   // 外部支付错误
	StatusPaymentUnknownError   Status = "PAYMENT_UNKNOWN_ERROR"    // 支付未知错误
	StatusPointValueNotEnough   Status = "POINT_VALUE_NOT_ENOUGH"   // 积分余额不足
)

var paymentStatus = map[Status]statusResult{
	StatusPaymentValueNotEnough: {level: log.LevelWarn, message: "余额不足"},
	StatusPayTimeout:            {level: log.LevelWarn, message: "支付超时"},
	StatusPaymentLocked:         {level: log.LevelWarn, message: "支付方式已被锁定"},
	StatusPaymentUnauthorized:   {level: log.LevelWarn, message: "支付方式未授权"},
	StatusAlipayUnauthorized:    {level: log.LevelWarn, message: "未授权支付宝"},
	StatusWechatUnauthorized:    {level: log.LevelWarn, message: "未授权微信"},
	StatusInvalidPayment:        {level: log.LevelWarn, message: "无效的支付方式"},
	StatusUnsupportedPayment:    {level: log.LevelWarn, message: "不支持的支付方式"},
	StatusInvalidCallbackURL:    {level: log.LevelWarn, message: "无效的回调地址"},
	StatusExternalPaymentError:  {level: log.LevelError, message: "外部支付错误"},
	StatusPaymentUnknownError:   {level: log.LevelError, message: "支付未知错误"},
	StatusPointValueNotEnough:   {level: log.LevelWarn, message: "积分余额不足"},
}

const (
	// RPCBusinessError ...
	RPCBusinessError codes.Code = 101
)
