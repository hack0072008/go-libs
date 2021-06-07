package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Options struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
	PoolSize int    `yaml:"pool_size" mapstructure:"pool_size"`
	// 读写超时时间
	TimeOut int `yaml:"timeout" mapstructure:"timeout"`
	// 缓存 key 统一前缀
	Prefix string `yaml:"prefix" mapstructure:"prefix"`
}

type Client struct {
	redisClient *redis.Client
	prefix      string
}

func NewCacheClient(ops Options) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         ops.Addr,
		Password:     ops.Password,
		DB:           ops.DB,
		PoolSize:     ops.PoolSize,
		ReadTimeout:  time.Duration(ops.TimeOut) * time.Second,
		WriteTimeout: time.Duration(ops.TimeOut) * time.Second,

		TLSConfig: nil,
	})
	if err := redisClient.Ping().Err(); err != nil {
		fmt.Printf("创建 redis client 失败: %v", err)
		return nil, err
	}

	cacheClient := &Client{
		redisClient: redisClient,
		prefix:      ops.Prefix,
	}
	return cacheClient, nil
}

func (c *Client) Close() error {
	return c.redisClient.Close()
}

func (c *Client) realKey(part string) string {
	if len(c.prefix) == 0 {
		return part
	}
	return c.prefix + "." + part
}

func (c *Client) Set(key, value string, expireTime ...int) error {
	// 设置过期时间
	expire := 0
	if len(expireTime) != 0 {
		expire = expireTime[0]
	}

	realKey := c.realKey(key)
	if err := c.redisClient.Set(realKey, value, time.Duration(expire)*time.Second).Err(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Get(key string) (string, error) {
	realKey := c.realKey(key)
	value, err := c.redisClient.Get(realKey).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (c *Client) Remove(key string) error {
	realKey := c.realKey(key)
	if err := c.redisClient.Del(realKey).Err(); err != nil {
		return err
	}
	return nil
}
