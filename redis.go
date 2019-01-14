package itea

import (
	"github.com/go-redis/redis"
	"time"
	"log"
	"fmt"
	"context"
	"encoding/json"
	"strings"
)

const (
	REDIS_HOST = ""
	REDIS_PORT = "6379"
	REDIS_DATABASE = 0
	REDIS_PASSWORD = ""
	REDIS_POOL_MAX_IDLE = 10
	REDIS_POOL_MAX_ACTIVE = 100
	REDIS_POOL_IDLE_TIMEOUT = 300
	REDIS_POOL_MAX_CONN_LIFETIME = 0
	REDIS_POOL_IDLE_CHECK_FREQUENCY = 60
)

type Redis struct {
	pool *redis.Client
	Ctx context.Context
}

func (p *Redis) Construct() {
	var conf RedisConf
	if c := p.Ctx.Value("redis").(*json.RawMessage); c != nil {
		err := json.Unmarshal(*c, &conf)
		if err != nil {
			panic(err)
		}
	} else {
		panic("Can not find database config of redis!")
	}

	redisPool := redis.NewClient(p.initOpt(conf))

	p.pool = redisPool

	//go func() {
	//	for {
	//		time.Sleep(time.Second)
	//		fmt.Printf("PoolStats, TotalConns: %d, FreeConns: %d\n", p.pool.PoolStats().TotalConns, p.pool.PoolStats().IdleConns)
	//	}
	//}()
}

func (p *Redis) initOpt(conf RedisConf) *redis.Options {
	host, port := REDIS_HOST, REDIS_PORT
	if !strings.EqualFold(conf.Host, "") {
		host = conf.Host
	}
	if !strings.EqualFold(conf.Port, "") {
		port = conf.Port
	}

	opt := &redis.Options{
		Addr:     			fmt.Sprintf("%s:%s", host, port),
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

func (p *Redis) Setex(key string, value string, expire int) (reply interface{}, err error) {
	start := time.Now()
	defer func() {
		log.Println("【Redis Setex】耗时：", time.Since(start))
	}()
	p.pool.Set(key, value, time.Duration(expire) * time.Second)
	return nil, nil
}

func (p *Redis) Get(key string) (value string, err error) {
	start := time.Now()
	defer func() {
		log.Println("【Redis Get】耗时：", time.Since(start))
	}()
	return p.pool.Get(key).String(), nil
}

func (p *Redis) Expire(key string, expire int) (reply interface{}, err error) {
	start := time.Now()
	defer func() {
		log.Println("【Redis Expire】耗时：", time.Since(start))
	}()
	p.pool.Expire(key, time.Duration(expire) * time.Second)
	return nil, nil
}