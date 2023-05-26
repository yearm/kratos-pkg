# kratos-pkg

## 简介
kratos-pkg 整合 [kratos](https://github.com/go-kratos/kratos)、[iris](https://github.com/kataras/iris) 的工具包。包含配置、错误码、序列化、日志、监控、服务注册发现、链路追踪、存储等工具

- **config**：配置模块，默认数据源为 nacos
- **ecode**：错误码定义
- **encoding**：序列化模块，支持 proto to map，默认使用 protojson
- **log**：内部使用的日志库
- **logger**：日志模块，默认对接 aliyun sls
- **metrics**：指标模块，默认对接 prometheus
- **registry**：服务注册发现模块，默认使用 k8s 服务发现
- **result**：http 响应结构定义
- **status**：重新包装 grpc error 类型
- **store**：存储层，支持 gorm、go-redis
- **trace**：链路追踪，支持 trace、slstrace，默认使用 trace
- **util**：内部工具包
- **xgrpc**：提供创建 grpc server client 以及默认 middlewares
- **xhttp**：提供创建 http server 和 monitor server 以及默认 middlewares，自定义 iris Context

## 组件说明（部分说明）

### xhttp

#### 自定义 Context
因为 iris.Context ([v12.1.8](https://github.com/kataras/iris/tree/v12.1.8)) 未实现 context.Context 接口，为了 trace 信息的传递、timeout ctx 传递以及打印日志所需的 Context 等，需自定义 Context 来组合 context.Context
```go
type stdContext struct {
	context.Context
}

type Context struct {
	iris.Context
	stdContext
}
```

#### 转换 kratos middleware 为 iris middleware
c.ResetRequest(c.Request().WithContext(ctx)) 为重点，目的是为了将 context.Context 继续往下传递，否则自定义的 Context 获取不到经过 tracing.Server() middleware 后的 trace 信息
```go
func Middlewares(m ...middleware.Middleware) iriscontext.Handler {
	chain := middleware.Chain(m...)
	return func(c iris.Context) {
		next := func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// NOTE: !!!
			c.ResetRequest(c.Request().WithContext(ctx))
			c.Next()
			return
		}
		next = chain(next)
		ctx := NewIrisContext(c.Request().Context(), c)
		if irisCtx, ok := FromIrisContext(ctx); ok {
			thttp.SetOperation(ctx, irisCtx.GetCurrentRoute().Path())
		}
		_, _ = next(ctx, c.Request())
	}
}
```

#### Render 时记录日志
log 中除了记录基本的一些参数之外，提供以注册 LogValuer function 的方式记录业务所需要的日志字段
```go
func (c *Context) log(result *result.Result) {
	// ...
	for key, valuer := range logValuers() {
		params[key] = valuer(c)
	}
	// ...
}

var (
	defaultLogValuers = make(map[string]LogValuer)
	logValuerOnce     = new(sync.Once)
)

type LogValuer func(ctx *Context) interface{}

func RegisterLogValuers(ms map[string]LogValuer) {
	logValuerOnce.Do(func() {
		defaultLogValuers = ms
	})
}

func logValuers() map[string]LogValuer {
	return defaultLogValuers
}
```

### xgrpc

#### 服务端
服务端 timeout 设置为 0，具体原因见注释。默认集成 grpc_prometheus
```go
func NewGRPCServer(opts ...kgrpc.ServerOption) *kgrpc.Server {
	var serverOpts = []kgrpc.ServerOption{
		kgrpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		kgrpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	}
	serverOpts = append(serverOpts, opts...)
	grpcHost := env.GetGRPCHost()
	grpcPort := env.GetGRPCPort()
	if grpcHost != "" && grpcPort > 0 {
		serverOpts = append(serverOpts, kgrpc.Address(fmt.Sprintf("%s:%d", grpcHost, grpcPort)))
	}
	// NOTE: kgrpc.Timeout 暂不用设置，因为 unaryServerInterceptor 中并没有 select case <-ctx.Done()
	// NOTE: 设置了反倒会改变 context 的超时传递时间，一般情况 client 的 context 带有超时时间
	// NOTE: 正常情况下 client 调用超时是为了避免链路阻塞堆积，server 端继续处理请求也属正常
	// NOTE: 未自定义设置服务端 timeout 时，kratos 框架默认设置为1秒，导致服务端调用时间过长或者链路较长时服务超时中断
	// NOTE: 所以此处设置 timeout 为0，即使用客户端调用传来的 ctx 中的超时控制
	serverOpts = append(serverOpts, kgrpc.Timeout(0))
	srv := kgrpc.NewServer(serverOpts...)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(srv.Server)
	return srv
}
```

#### 客户端
使用 sync.Map 保存 grpc 连接
```go
var (
	connMap sync.Map
	connLock sync.Mutex
)

func GetRPCClientConn(configPath string, opts ...kgrpc.ClientOption) *grpc.ClientConn {
	var (
		err  error
		conn *grpc.ClientConn
	)

	connLock.Lock()
	defer connLock.Unlock()
	if v, ok := connMap.Load(configPath); ok {
		return v.(*grpc.ClientConn)
	}
	defer func() {
		if conn != nil {
			connMap.Store(configPath, conn)
		}
	}()

	endpoint := viper.GetString(fmt.Sprintf("%s.endpoint", configPath))
	if endpoint == "" {
		logrus.Panicln("endpoint is nil, config path:", configPath)
	}
	timeout := viper.GetInt(fmt.Sprintf("%s.timeout", configPath))
	dialWithCredentials := viper.GetBool(fmt.Sprintf("%s.dialWithCredentials", configPath))

	clientOpts := []kgrpc.ClientOption{kgrpc.WithEndpoint(endpoint)}
	clientOpts = append(clientOpts, opts...)
	if d := registry.NewDiscovery(); d != nil {
		clientOpts = append(clientOpts, kgrpc.WithDiscovery(d))
	}

	if timeout >= 0 {
		clientOpts = append(clientOpts, kgrpc.WithTimeout(time.Duration(timeout)*time.Second))
	}

	if dialWithCredentials {
		clientOpts = append(clientOpts, kgrpc.WithTLSConfig(&tls.Config{}))
		conn, err = kgrpc.Dial(context.Background(), clientOpts...)
	} else {
		conn, err = kgrpc.DialInsecure(context.Background(), clientOpts...)
	}

	if err != nil {
		logrus.Panicln("grpc dial error:", err)
	}
	logrus.Infoln("Connecting at", endpoint)
	return conn
}
```

#### 参数校验中间件
使用 [protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag) 生成 proto struct 的时候注入 struct tag，再通过 [validator](https://github.com/go-playground/validator/v10) 参数校验改 proto struct，支持返回自定义 errMsg tag 
```go
func Validator() middleware.Middleware {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("errMsg")
	})
	zhT := zh.New()
	uni := ut.New(zhT, zhT)
	trans, _ := uni.GetTranslator("zh")
	if err := zhtrans.RegisterDefaultTranslations(validate, trans); err != nil {
		logrus.Panicln("validator middleware register translations error:", err)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if _err := validate.Struct(req); _err != nil {
				if fieldErrors, ok := (_err).(validator.ValidationErrors); ok {
					for _, fieldError := range fieldErrors {
						var errMsg string
						translateValue := fieldError.Translate(trans)
						// NOTE: Field() 和 StructField() 不相等说明取到了 errMsg tag 值
						if fieldError.Field() != fieldError.StructField() {
							// NOTE: 翻译时取的值是 Field()，由于前面 RegisterTagNameFunc 取的是 errMsg tag 对应的值，所以这里翻译后要替换成 StructField()
							translateValue = strings.Replace(translateValue, fieldError.Field(), fieldError.StructField(), 1)
							errMsg = fieldError.Field()
						}
						_fieldError := errors.New(fmt.Sprintf("%s:%s", fieldError.StructNamespace(), translateValue))
						if errMsg != "" {
							return nil, status.ErrorWithMsg(_fieldError, ecode.StatusInvalidRequest, errMsg)
						}
						return nil, status.Error(_fieldError, ecode.StatusInvalidRequest)
					}
				}
				return nil, status.Error(_err, ecode.StatusInvalidRequest)
			}
			return handler(ctx, req)
		}
	}
}
```

### status

#### 包装 grpc error
rpc 层业务错误需要再进行包装，带上错误码以及其他信息再对其抛出
```go
func Error(err error, status ecode.Status, levels ...log.Level) error {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	_struct, _ := structpb.NewStruct(map[string]interface{}{
		"status": status.String(),
		"msg":    status.Message(),
		"errorf": debug.Caller(2, 3),
		"level":  level(levels...).String(),
	})
	st, _ := gstatus.New(ecode.RPCBusinessError, fmt.Sprintf("[%s]%s", env.GetServiceName(), errMsg)).WithDetails(_struct)
	return st.Err()
}

func FromError(err error) (*gstatus.Status, *structpb.Struct) {
	if err == nil {
		return nil, nil
	}
	if st, ok := gstatus.FromError(err); ok {
		for _, detail := range st.Details() {
			if _struct, ok := detail.(*structpb.Struct); ok {
				return st, _struct
			}
		}
		return st, nil
	}
	return nil, nil
}
```

### result

#### 解析 rpc 层抛出的错误
Result 为 http 响应结构定义，fromError 作用是将 rpc error 解析转换成 Result
```go
func fromError(err error, useErrMsg bool) *Result {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return nil
	}
	errorf := func() string { return debug.Caller(5, 3) }
	switch st.Code() {
	case codes.Canceled:
		return NewWithErrorf(ecode.StatusCancelled, err, errorf, log.LevelWarn)
	case codes.Unknown:
		return NewWithErrorf(ecode.StatusUnknownError, err, errorf)
	case codes.DeadlineExceeded:
		return NewWithErrorf(ecode.StatusRequestTimeout, err, errorf)
	case codes.Internal:
		return NewWithErrorf(ecode.StatusInternalServerError, err, errorf)
	case codes.Unavailable:
		return NewWithErrorf(ecode.StatusTemporarilyUnavailable, err, errorf)
	case ecode.RPCBusinessError:
		var (
			_struct *structpb.Struct
			_ok     bool
		)
		for _, detail := range st.Details() {
			if _struct, _ok = detail.(*structpb.Struct); _ok {
				break
			}
		}
		if _struct != nil {
			structMap := _struct.AsMap()
			result := NewWithErrorf(ecode.Status(gconv.String(structMap["status"])), err, errorf, log.ParseLevel(gconv.String(structMap["level"])))
			if useErrMsg {
				result.SetMessage(gconv.String(structMap["msg"]))
			}
			return result
		}
		return NewWithErrorf(ecode.StatusInternalServerError, err, errorf)
	default:
		return NewWithErrorf(ecode.StatusInternalServerError, err, errorf)
	}
}
```
