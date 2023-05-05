package xredis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yearm/kratos-pkg/config/env"
	"time"
)

// Redsync ...
type Redsync struct {
	*redsync.Redsync
	opts []redsync.Option
}

// NewDefaultRedsync ...
func NewDefaultRedsync(redisDB RedisDB) *Redsync {
	return NewRedsync(redisDB, env.GetRedisMutexConfigPath())
}

// NewRedsync ...
func NewRedsync(redisDB RedisDB, configPath string) *Redsync {
	return newRedsync(redisDB, configPath)
}

// NewDefaultMutex ...
func (r *Redsync) NewDefaultMutex(name string, options ...redsync.Option) *redsync.Mutex {
	return r.NewMutex(name, append(r.opts, options...)...)
}

// newRedsync ...
func newRedsync(client *redis.Client, configPath string) *Redsync {
	pool := goredis.NewPool(client)
	expiry := viper.GetInt(fmt.Sprintf("%s.expiry", configPath))
	tries := viper.GetInt(fmt.Sprintf("%s.tries", configPath))
	if expiry == 0 || tries == 0 {
		logrus.Panicln("redsync option params not set")
	}
	return &Redsync{
		Redsync: redsync.New(pool),
		opts: []redsync.Option{
			redsync.WithExpiry(time.Duration(expiry) * time.Second),
			redsync.WithTries(tries),
		},
	}
}
