package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Config struct {
	Addr        string        `env:"addr"`
	Password    string        `env:"password"`
	DB          int           `env:"db"`
	MaxRetries  int           `env:"max_retries"`
	DialTimeout time.Duration `env:"dial_timeout"`
	Timeout     time.Duration `env:"timeout"`
}

type Redis struct {
	cfg    Config
	log    Logger
	client *redis.Client
}

func New(cfg Config, logger Logger) *Redis {
	return &Redis{cfg: cfg, log: logger}
}

func (r *Redis) Run(_ context.Context) {}

func (r *Redis) Init() error {
	options := &redis.Options{
		Addr:     r.cfg.Addr,
		Password: r.cfg.Password,
		DB:       r.cfg.DB,
	}

	client := redis.NewClient(options)

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return errors.Wrap(err, "could not connect to Redis")
	}

	r.client = client
	r.log.Info("Redis connected successfully")
	return nil
}

func (r *Redis) Stop() {
	err := r.client.Close()
	if err != nil {
		r.log.Error("r.Client.Close() err:", err)
	}
	r.log.Info("r.Client.Close() closed")
}
