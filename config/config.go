package config

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/iancoleman/orderedmap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	nconfig "github.com/yearm/kratos-pkg/config/nacos"
	"github.com/yearm/kratos-pkg/util/mergeomap"
	"os"
)

var (
	runtimeConfPath = flag.String("runtime", "./config/runtime.config.json", "runtime config file path")
	defaultConfPath = flag.String("default", "./config/default.json", "default config file path")
	localConfPath   = flag.String("local", "./config/local.json", "local config file path")
	aliyunConfPath  = flag.String("aliyun", "./config/aliyun.json", "aliyun mse config file path")
)

// Init ...
func Init() {
	flag.Parse()
	aliyunViper := viper.New()
	aliyunViper.AddConfigPath("./")
	aliyunViper.SetConfigFile(*aliyunConfPath)
	if err := aliyunViper.ReadInConfig(); err != nil {
		logrus.Errorln("read aliyun config error:", err)
	}

	defer func() {
		viper.AddConfigPath("./")
		viper.SetConfigFile(*runtimeConfPath)
		if err := viper.ReadInConfig(); err != nil {
			logrus.Panicln("read runtime config error:", err)
		}
		_ = viper.MergeConfigMap(aliyunViper.AllSettings())
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			logrus.Println("config file changed:", e.Name)
		})
	}()

	content, contentCh, _ := nconfig.Load(aliyunViper)
	writeFile(*runtimeConfPath, *defaultConfPath, *localConfPath, content)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorln("[recovery] watch config error:", err)
			}
		}()
		for content := range contentCh {
			writeFile(*runtimeConfPath, *defaultConfPath, *localConfPath, content)
		}
	}()
}

// writeFile ...
func writeFile(runtimeConfPath, defaultConfPath, localConfPath, content string) {
	var (
		runtimeConfMap = orderedmap.New()
		localConfMap   = orderedmap.New()
		contentMap     = orderedmap.New()
	)

	// NOTE: 优先级 local > nacos > default
	defaultConf, err := os.ReadFile(defaultConfPath)
	if err == nil {
		if err = json.Unmarshal(defaultConf, &runtimeConfMap); err != nil {
			logrus.Panicln("default config unmarshal error:", err)
		}
	}

	if content != "" {
		if err = json.Unmarshal([]byte(content), &contentMap); err != nil {
			logrus.Panicln("content unmarshal error:", err)
		}
		runtimeConfMap = mergeomap.Merge(runtimeConfMap, contentMap)
	}

	localConf, err := os.ReadFile(localConfPath)
	if err == nil {
		if err = json.Unmarshal(localConf, &localConfMap); err != nil {
			logrus.Panicln("local config unmarshal error:", err)
		}
		runtimeConfMap = mergeomap.Merge(runtimeConfMap, localConfMap)
	}

	runtimeConf := new(bytes.Buffer)
	encoder := json.NewEncoder(runtimeConf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(runtimeConfMap); err != nil {
		logrus.Panicln("runtime config encode error:", err)
	}
	err = os.WriteFile(runtimeConfPath, runtimeConf.Bytes(), os.ModePerm)
	if err != nil {
		logrus.Panicln("runtime config write file error:", err)
	}
	logrus.Println(runtimeConfPath)
}
