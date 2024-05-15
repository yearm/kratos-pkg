package xgormv2

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"
	"github.com/yearm/kratos-pkg/config/env"
	semconv "go.opentelemetry.io/otel/semconv/v1.14.0"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"time"

	"github.com/hoisie/mustache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// NewDBClient ...
func NewDBClient(configPath string, opts ...gorm.Option) *gorm.DB {
	user := viper.GetString(fmt.Sprintf("%s.user", configPath))
	password := viper.GetString(fmt.Sprintf("%s.password", configPath))
	host := viper.GetString(fmt.Sprintf("%s.host", configPath))
	port := viper.GetString(fmt.Sprintf("%s.port", configPath))
	database := viper.GetString(fmt.Sprintf("%s.database", configPath))
	maxIdle := viper.GetInt(fmt.Sprintf("%s.maxIdle", configPath))
	maxOpen := viper.GetInt(fmt.Sprintf("%s.maxOpen", configPath))
	maxLifetime := viper.GetInt(fmt.Sprintf("%s.maxLifetime", configPath))
	charset := viper.GetString(fmt.Sprintf("%s.charset", configPath))
	slowThreshold := viper.GetInt(fmt.Sprintf("%s.slowThreshold", configPath))

	url := mustache.Render("{{user}}:{{password}}@tcp({{host}}:{{port}})/{{database}}?charset={{charset}}&parseTime=True&loc=Local", map[string]interface{}{
		"user":     user,
		"password": password,
		"host":     host,
		"port":     port,
		"database": database,
		"charset":  lo.If(charset == "", "utf8").Else(charset),
	})

	ormConfigs := []gorm.Option{
		&gorm.Config{
			SkipDefaultTransaction: true,
			NamingStrategy:         schema.NamingStrategy{SingularTable: true},
			Logger:                 getLogger(slowThreshold),
		},
	}
	ormConfigs = append(ormConfigs, opts...)
	db, err := gorm.Open(mysql.Open(url), ormConfigs...)
	if err != nil {
		logrus.Panicln(fmt.Sprintf("failed to connect database:[%s], error:%s", url, err))
	}
	err = db.Use(tracing.NewPlugin(
		tracing.WithAttributes(semconv.DBUserKey.String(user)),
		tracing.WithDBName(database),
	))
	if err != nil {
		logrus.Panicln(fmt.Sprintf("failed to use tracing plugin:[%s], error:%s", url, err))
	}
	sdb, _ := db.DB()
	sdb.SetMaxIdleConns(maxIdle)
	sdb.SetMaxOpenConns(maxOpen)
	sdb.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)
	return db
}

func getLogger(slowThreshold int) logger.Interface {
	if env.IsLocal() {
		return logger.Default.LogMode(logger.Info)
	}
	logLevel := logger.Warn
	if env.IsDevelopment() {
		logLevel = logger.Info
	}
	return &aliyunLogger{
		Config: logger.Config{
			SlowThreshold:             time.Duration(lo.If(slowThreshold > 0, slowThreshold).Else(200)) * time.Millisecond,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logLevel,
		},
	}
}

type aliyunLogger struct {
	logger.Config
}

func (a *aliyunLogger) LogMode(level logger.LogLevel) logger.Interface {
	al := *a
	al.LogLevel = level
	return &al
}

func (a *aliyunLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if a.LogLevel >= logger.Info {
		log.Context(ctx).Infof(s, i)
	}
}

func (a *aliyunLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if a.LogLevel >= logger.Warn {
		log.Context(ctx).Warnf(s, i)
	}
}

func (a *aliyunLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if a.LogLevel >= logger.Error {
		log.Context(ctx).Errorf(s, i)
	}
}

func (a *aliyunLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if a.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && a.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !a.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			if errors.Is(err, logger.ErrRecordNotFound) {
				a.Info(ctx, "%s [%.3fms] [rows:%v] %s", err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				a.Error(ctx, "%s [%.3fms] [rows:%v] %s", err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			}
		} else {
			if errors.Is(err, logger.ErrRecordNotFound) {
				a.Info(ctx, "%s [%.3fms] [rows:%v] %s", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			} else {
				a.Error(ctx, "%s [%.3fms] [rows:%v] %s", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	case elapsed > a.SlowThreshold && a.SlowThreshold != 0 && a.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", a.SlowThreshold)
		if rows == -1 {
			a.Warn(ctx, "%s [%.3fms] [rows:%v] %s", slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			a.Warn(ctx, "%s [%.3fms] [rows:%v] %s", slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case a.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			a.Info(ctx, "[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			a.Info(ctx, "[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
