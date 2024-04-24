package alogger

import (
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yearm/kratos-pkg/config/env"
	"google.golang.org/protobuf/proto"
	"time"
)

// Logger ...
type Logger struct {
	project  string
	logStore string
	level    klog.Level
	producer *producer.Producer
}

// NewLogger ...
func NewLogger() (klog.Logger, func()) {
	accessKey := env.GetAliyunAccessKey()
	secretKey := env.GetAliyunSecretKey()
	if accessKey == "" || secretKey == "" {
		logrus.Panicln(fmt.Sprintf("aliyun logger config error, accessKey[%s], secretKey[%s]", accessKey, secretKey))
	}
	endpoint := viper.GetString(env.PathAliyunSlsEndpoint)
	project := viper.GetString(env.PathAliyunSlsProject)
	logStore := viper.GetString(env.PathAliyunSlsLogStore)
	levelStr := viper.GetString(env.PathAliyunSlsLogLevel)
	if endpoint == "" || project == "" || logStore == "" {
		logrus.Panicln(fmt.Sprintf("aliyun logger config error, endpoint:[%s], project:[%s], logStore:[%s]", endpoint, project, logStore))
	}

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = endpoint
	producerConfig.CredentialsProvider = sls.NewStaticCredentialsProvider(accessKey, secretKey, "")
	producerInst := producer.InitProducer(producerConfig)
	producerInst.Start()

	return &Logger{
			project:  project,
			logStore: logStore,
			level:    klog.ParseLevel(levelStr),
			producer: producerInst,
		}, func() {
			if err := producerInst.Close(5 * 1000); err != nil {
				logrus.Errorln("aliyun logger close error:", err)
				return
			}
			logrus.Println("aliyun logger graceful close")
		}
}

// Log ...
func (a *Logger) Log(level klog.Level, keyValues ...interface{}) error {
	if level < a.level {
		return nil
	}
	if len(keyValues) == 0 {
		return nil
	}
	if len(keyValues)%2 != 0 {
		keyValues = append(keyValues, "")
	}

	contents := make([]*sls.LogContent, 0, len(keyValues)/2+1)
	contents = append(contents, &sls.LogContent{
		Key:   proto.String("level"),
		Value: proto.String(level.String()),
	})
	for i := 0; i < len(keyValues); i += 2 {
		contents = append(contents, &sls.LogContent{
			Key:   proto.String(gconv.String(keyValues[i])),
			Value: proto.String(gconv.String(keyValues[i+1])),
		})
	}

	logInst := &sls.Log{
		Time:     proto.Uint32(uint32(time.Now().Unix())),
		Contents: contents,
	}
	err := a.producer.SendLog(a.project, a.logStore, "", "", logInst)
	return err
}
