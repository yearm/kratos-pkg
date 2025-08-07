package nacos

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/contrib/config/nacos/v2"
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/yearm/kratos-pkg/errors"
)

type Nacos struct {
	AccessKey   string `yaml:"accessKey" json:"accessKey" xml:"accessKey"`
	SecretKey   string `yaml:"secretKey" json:"secretKey" xml:"secretKey"`
	Endpoint    string `yaml:"endpoint" json:"endpoint" xml:"endpoint"`
	NamespaceId string `yaml:"namespaceId" json:"namespaceId" xml:"namespaceId"`
	RegionId    string `yaml:"regionId" json:"regionId" xml:"regionId"`
	Group       string `yaml:"group" json:"group" xml:"group"`
	DataId      string `yaml:"dataId" json:"dataId" xml:"dataId"`
	CacheDir    string `yaml:"cacheDir" json:"cacheDir" xml:"cacheDir"`
	LogDir      string `yaml:"logDir" json:"logDir" xml:"logDir"`
	LogLevel    string `yaml:"logLevel" json:"logLevel" xml:"logLevel"`
}

// Config configuration parameters of nacos client.
type Config struct {
	Nacos Nacos `yaml:"nacos" json:"nacos" xml:"nacos"`
}

func NewConfig(accessKey string, secretKey string, endpoint string, namespaceId string, regionId string, group string, dataId string, cacheDir string, logDir string, logLevel string) *Config {
	return &Config{
		Nacos: Nacos{
			AccessKey:   accessKey,
			SecretKey:   secretKey,
			Endpoint:    endpoint,
			NamespaceId: namespaceId,
			RegionId:    regionId,
			Group:       group,
			DataId:      dataId,
			CacheDir:    cacheDir,
			LogDir:      logDir,
			LogLevel:    logLevel,
		},
	}
}

func (c *Config) IsValid() bool {
	if c == nil {
		return false
	}
	if c.Nacos.AccessKey == "" ||
		c.Nacos.SecretKey == "" ||
		c.Nacos.Endpoint == "" ||
		c.Nacos.NamespaceId == "" ||
		c.Nacos.RegionId == "" ||
		c.Nacos.Group == "" ||
		c.Nacos.DataId == "" ||
		c.Nacos.CacheDir == "" ||
		c.Nacos.LogDir == "" ||
		c.Nacos.LogLevel == "" {
		return false
	}
	return true
}

// NewConfigSource creates a nacos config source.
func NewConfigSource(c *Config) (kconfig.Source, error) {
	if !c.IsValid() {
		return nil, errors.Errorf("invalid nacos config[%v]", c)
	}

	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			TimeoutMs:           5 * 1000,
			NamespaceId:         c.Nacos.NamespaceId,
			RegionId:            c.Nacos.RegionId,
			AccessKey:           c.Nacos.AccessKey,
			SecretKey:           c.Nacos.SecretKey,
			CacheDir:            c.Nacos.CacheDir,
			NotLoadCacheAtStart: true,
			LogDir:              c.Nacos.LogDir,
			LogLevel:            c.Nacos.LogLevel,
		},
		ServerConfigs: []constant.ServerConfig{
			*constant.NewServerConfig(c.Nacos.Endpoint, 8848),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "clients.NewConfigClient failed")
	}
	return config.NewConfigSource(client, config.WithGroup(c.Nacos.Group), config.WithDataID(c.Nacos.DataId)), nil
}

// NewConfigSourceFormFile creates a nacos config source from the file.
func NewConfigSourceFormFile(path string) (kconfig.Source, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open failed, path = %v", path)
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "io.ReadAll failed")
	}

	format := strings.TrimPrefix(filepath.Ext(path), ".")
	codec := encoding.GetCodec(format)
	if codec == nil {
		return nil, errors.Errorf("unknown encoding format: %s", format)
	}

	var c *Config
	if err := codec.Unmarshal(data, &c); err != nil {
		return nil, errors.Wrap(err, "codec.Unmarshal failed")
	}
	return NewConfigSource(c)
}
