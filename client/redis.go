package client

import (
	"context"
	"fmt"
	"github.com/CalvinDjy/iteaGo/constant"
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/CalvinDjy/iteaGo/system"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

const (
	REDIS_KEY 						= "redis"
	REDIS_HOST 						= ""
	REDIS_PORT 						= 6379
	REDIS_DATABASE 					= 0
	REDIS_PASSWORD 					= ""
	REDIS_POOL_MAX_IDLE 			= 10
	REDIS_POOL_MAX_ACTIVE 			= 100
	REDIS_POOL_IDLE_TIMEOUT 		= 300
	REDIS_POOL_MAX_CONN_LIFETIME 	= 0
	REDIS_POOL_IDLE_CHECK_FREQUENCY = 60
)

type RedisConf struct {
	Host 			string
	Port 			int
	Database 		int
	Password 		string
	MaxIdle 		int
	MaxActive 		int
	IdleTimeout 	int
	MaxConnLifetime int
	IdleCheck 		int
}

type Redis struct {
	pool 			*redis.Client
	Ctx 			context.Context
	debug 			bool
}

func (p *Redis) Construct() {
	if d, ok := p.Ctx.Value(constant.DEBUG).(bool); ok {
		p.debug = d
	}

	c := system.Conf.GetStruct(fmt.Sprintf("%s.%s", constant.DATABASE_KEY, REDIS_KEY), RedisConf{})
	if c == nil {
		panic("Can not find database config of redis!")
	}

	p.pool = redis.NewClient(p.initOpt(c.(*RedisConf)))

	ilog.Info("redis pool create success")
	//go func() {
	//	for {
	//		time.Sleep(time.Second)
	//		fmt.Printf("PoolStats, TotalConns: %d, FreeConns: %d\n", p.pool.PoolStats().TotalConns, p.pool.PoolStats().IdleConns)
	//	}
	//}()
}

func (p *Redis) initOpt(conf *RedisConf) *redis.Options {
	host, port := REDIS_HOST, REDIS_PORT
	if !strings.EqualFold(conf.Host, "") {
		host = conf.Host
	}
	if conf.Port == 0 {
		port = conf.Port
	}

	opt := &redis.Options{
		Addr:     			fmt.Sprintf("%s:%d", host, port),
		Password: 			REDIS_PASSWORD,
		DB:       			REDIS_DATABASE,
		PoolSize: 			REDIS_POOL_MAX_ACTIVE,
		MinIdleConns:		REDIS_POOL_MAX_IDLE,
		IdleTimeout: 		time.Duration(REDIS_POOL_IDLE_TIMEOUT) * time.Second,
		MaxConnAge:			time.Duration(REDIS_POOL_MAX_CONN_LIFETIME) * time.Second,
		IdleCheckFrequency: time.Duration(REDIS_POOL_IDLE_CHECK_FREQUENCY) * time.Second,
	}

	if !strings.EqualFold(conf.Password, "") {
		opt.Password = conf.Password
	}
	if conf.Database > 0 {
		opt.DB = conf.Database
	}
	if conf.MaxIdle > 0 {
		opt.MinIdleConns = conf.MaxIdle
	}
	if conf.MaxActive > 0 {
		opt.PoolSize = conf.MaxActive
	}
	if conf.IdleTimeout > 0 {
		opt.IdleTimeout = time.Duration(conf.IdleTimeout) * time.Second
	}
	if conf.MaxConnLifetime > 0 {
		opt.MaxConnAge = time.Duration(conf.MaxConnLifetime) * time.Second
	}
	if conf.IdleCheck > 0 {
		opt.IdleCheckFrequency = time.Duration(conf.IdleCheck) * time.Second
	}

	return opt
}

func (p *Redis) Setex(key string, value string, expire int) (string, error) {
	if p.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Redis Setex】耗时：", time.Since(start))
		}()
	}
	return p.pool.Set(key, value, time.Duration(expire) * time.Second).Result()
}

func (p *Redis) Get(key string) (string, error) {
	if p.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Redis Get】耗时：", time.Since(start))
		}()
	}
	return p.pool.Get(key).Val(), nil
}

func (p *Redis) Expire(key string, expire int) (bool, error) {
	if p.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Redis Expire】耗时：", time.Since(start))
		}()
	}
	return p.pool.Expire(key, time.Duration(expire) * time.Second).Result()
}

func (p *Redis) Delete(key string) (int64, error) {
	if p.debug {
		start := time.Now()
		defer func() {
			ilog.Info("【Redis Delete】耗时：", time.Since(start))
		}()
	}
	return p.pool.Del(key).Result()
}