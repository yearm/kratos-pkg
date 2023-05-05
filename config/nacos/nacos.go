package nconfig

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yearm/kratos-pkg/config/env"
)

// Load ...
func Load(aliyunViper *viper.Viper) (content string, contentCh chan string, err error) {
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		constant.KEY_SERVER_CONFIGS: []constant.ServerConfig{{
			IpAddr: aliyunViper.GetString(env.PathAliyunMseEndpoint),
			Port:   8848,
		}},
		constant.KEY_CLIENT_CONFIG: constant.ClientConfig{
			TimeoutMs:   5000,
			NamespaceId: aliyunViper.GetString(env.PathAliyunMseNamespaceID),
			RegionId:    aliyunViper.GetString(env.PathAliyunMseRegionID),
			AccessKey:   aliyunViper.GetString(env.PathAliyunAccessKey),
			SecretKey:   aliyunViper.GetString(env.PathAliyunSecretKey),
			CacheDir:    aliyunViper.GetString(env.PathAliyunMseCacheDir),
			LogDir:      aliyunViper.GetString(env.PathAliyunMseLogDir),
			LogLevel:    aliyunViper.GetString(env.PathAliyunMseLogLevel),
		},
	})
	if err != nil {
		logrus.Errorln("load nacos config client error:", err)
		return
	}

	contentCh = make(chan string, 1)
	configParam := vo.ConfigParam{
		DataId: aliyunViper.GetString(env.PathAliyunMseDataID),
		Group:  aliyunViper.GetString(env.PathAliyunMseGroup),
		OnChange: func(namespace, group, dataId, data string) {
			contentCh <- data
		},
	}
	content, err = configClient.GetConfig(configParam)
	if err != nil {
		logrus.Errorln("get nacos config error:", err)
		return
	}
	if err = configClient.ListenConfig(configParam); err != nil {
		logrus.Errorln("listen nacos config error:", err)
		return
	}
	return
}
