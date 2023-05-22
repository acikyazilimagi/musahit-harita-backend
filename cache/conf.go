package cache

import "github.com/caarlos0/env/v8"

type Config struct {
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:"eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81"`
}

func NewConfig() *Config {
	conf := &Config{}
	err := env.Parse(conf)
	if err != nil {
		panic(err)
	}

	return conf
}
