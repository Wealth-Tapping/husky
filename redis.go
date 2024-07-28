package husky

import (
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Url      string  `toml:"host"`
	Password string  `toml:"password"`
	Prefix   *string `toml:"prefix"`
}

type RedisIns struct {
	*redis.Client
	Prefix *string
}

var _RedisIns map[string]*RedisIns

func init() {
	_RedisIns = make(map[string]*RedisIns)
}

func InitRedis(config *RedisConfig, key ...string) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Url,
		Password: config.Password,
	})
	r := RedisIns{
		Client: client,
		Prefix: config.Prefix,
	}
	if len(key) == 0 {
		_RedisIns[""] = &r
	} else {
		_RedisIns[key[0]] = &r
	}
}

func Redis(key ...string) *RedisIns {
	if len(key) == 0 {
		return _RedisIns[""]
	} else {
		return _RedisIns[key[0]]
	}
}
