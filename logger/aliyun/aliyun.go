package aliyun

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"
	"github.com/yearm/kratos-pkg/errors"
	"github.com/yearm/kratos-pkg/utils/bytesconv"
	"github.com/yearm/kratos-pkg/utils/gjson"
	"google.golang.org/protobuf/proto"
)

// Config configuration parameters of aliyun logger client.
type Config struct {
	accessKey string
	secretKey string
	endpoint  string
	project   string
	logstore  string
	level     string
}

func NewConfig(accessKey string, secretKey string, endpoint string, project string, logstore string, level string) *Config {
	return &Config{accessKey: accessKey, secretKey: secretKey, endpoint: endpoint, project: project, logstore: logstore, level: level}
}

func (c *Config) IsValid() bool {
	if c == nil {
		return false
	}
	return c.accessKey != "" && c.secretKey != "" && c.endpoint != "" && c.project != "" && c.logstore != "" && c.level != ""
}

// Logger see more detail https://github.com/aliyun/aliyun-log-go-sdk
type Logger struct {
	project  string
	logstore string
	level    log.Level
	producer *producer.Producer
}

// NewLogger creates a aliyun logger by Config.
func NewLogger(c *Config) (log.Logger, func(), error) {
	if !c.IsValid() {
		return nil, nil, errors.Errorf("invalid config[%v]", c)
	}
	config := producer.GetDefaultProducerConfig()
	config.Endpoint = c.endpoint
	config.DisableRuntimeMetrics = true
	config.CredentialsProvider = sls.NewStaticCredentialsProvider(c.accessKey, c.secretKey, "")
	p, err := producer.NewProducer(config)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "producer.NewProducer failed, config = %v", c)
	}
	p.Start()
	l := &Logger{
		project:  c.project,
		logstore: c.logstore,
		level:    log.ParseLevel(c.level),
		producer: p,
	}
	return l, func() { l.producer.SafeClose() }, nil
}

// Log implement kratos log.Logger.
func (l *Logger) Log(level log.Level, kvs ...any) error {
	if level < l.level {
		return nil
	}
	if len(kvs) == 0 {
		return nil
	}
	if len(kvs)%2 != 0 {
		kvs = append(kvs, "")
	}

	now := time.Now()
	contents := make([]*sls.LogContent, 0, len(kvs)/2+1)
	contents = append(contents, &sls.LogContent{
		Key:   lo.ToPtr("time"),
		Value: lo.ToPtr(now.Format(time.RFC3339)),
	})
	contents = append(contents, &sls.LogContent{
		Key:   lo.ToPtr(level.Key()),
		Value: lo.ToPtr(strings.ToLower(level.String())),
	})
	for i := 0; i < len(kvs); i += 2 {
		contents = append(contents, &sls.LogContent{
			Key:   lo.ToPtr(toString(kvs[i])),
			Value: lo.ToPtr(toString(kvs[i+1])),
		})
	}
	logInst := &sls.Log{
		Time:     proto.Uint32(uint32(now.Unix())),
		Contents: contents,
	}
	return l.producer.SendLog(l.project, l.logstore, "", "", logInst)
}

// toString convert any type to string.
func toString(v any) string {
	var key string
	if v == nil {
		return key
	}
	switch v := v.(type) {
	case float64:
		key = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		key = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		key = strconv.Itoa(v)
	case uint:
		key = strconv.FormatUint(uint64(v), 10)
	case int8:
		key = strconv.Itoa(int(v))
	case uint8:
		key = strconv.FormatUint(uint64(v), 10)
	case int16:
		key = strconv.Itoa(int(v))
	case uint16:
		key = strconv.FormatUint(uint64(v), 10)
	case int32:
		key = strconv.Itoa(int(v))
	case uint32:
		key = strconv.FormatUint(uint64(v), 10)
	case int64:
		key = strconv.FormatInt(v, 10)
	case uint64:
		key = strconv.FormatUint(v, 10)
	case string:
		key = v
	case bool:
		key = strconv.FormatBool(v)
	case []byte:
		key = bytesconv.BytesToString(v)
	case fmt.Stringer:
		key = v.String()
	default:
		key = gjson.MustMarshalToString(v)
	}
	return key
}
