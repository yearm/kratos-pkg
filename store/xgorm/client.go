package xgorm

import (
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yearm/kratos-pkg/config/env"
	"time"
)

// NewDBClient ...
func NewDBClient(configPath string) *gorm.DB {
	dialect := viper.GetString(fmt.Sprintf("%s.dialect", configPath))
	user := viper.GetString(fmt.Sprintf("%s.user", configPath))
	password := viper.GetString(fmt.Sprintf("%s.password", configPath))
	host := viper.GetString(fmt.Sprintf("%s.host", configPath))
	port := viper.GetString(fmt.Sprintf("%s.port", configPath))
	database := viper.GetString(fmt.Sprintf("%s.database", configPath))
	maxIdle := viper.GetInt(fmt.Sprintf("%s.maxIdle", configPath))
	maxOpen := viper.GetInt(fmt.Sprintf("%s.maxOpen", configPath))
	maxLifetime := viper.GetInt(fmt.Sprintf("%s.maxLifetime", configPath))

	url := mustache.Render("{{user}}:{{password}}@tcp({{host}}:{{port}})/{{database}}?charset=utf8&parseTime=True&loc=Local", map[string]interface{}{
		"user":     user,
		"password": password,
		"host":     host,
		"port":     port,
		"database": database,
	})
	db, err := gorm.Open(dialect, url)
	if err != nil {
		logrus.Panicln(fmt.Sprintf("failed to connect database:[%s], error:%s", url, err))
	}

	if env.IsDevelopment() {
		db.LogMode(true)
	}
	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxOpen)
	db.DB().SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)
	db.SingularTable(true)

	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	return db
}

// updateTimeStampForCreateCallback ...
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		now := time.Now().Local()

		if createdTimestampField, ok := scope.FieldByName("CreatedTimestamp"); ok {
			if createdTimestampField.IsBlank {
				_ = createdTimestampField.Set(now.Unix())
			}
		}
		if createdAtField, ok := scope.FieldByName("CreatedAt"); ok {
			if createdAtField.IsBlank {
				_ = createdAtField.Set(now)
			}
		}
		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			if updatedAtField.IsBlank {
				_ = updatedAtField.Set(now)
			}
		}
	}
}
