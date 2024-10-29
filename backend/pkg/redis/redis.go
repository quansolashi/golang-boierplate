package redis

import (
	"context"
	"time"

	gredis "github.com/redis/go-redis/v9"
)

type Client struct {
	DB       *gredis.Client
	Duration int64
}

type Params struct {
	Address  string
	Password string
}

type options struct {
	retries int
}

type Option func(opts *options)

func WithMaxRetries(retries int) Option {
	return func(opts *options) {
		opts.retries = retries
	}
}

func NewClient(params *Params, opts ...Option) *Client {
	rdb := newRedisDB(params, &options{})

	dopts := &options{
		retries: 3,
	}
	for i := range opts {
		opts[i](dopts)
	}

	return &Client{
		DB: rdb,
	}
}

func (c *Client) Set(ctx context.Context, key string, value interface{}) error {
	exp := time.Duration(c.Duration) * time.Second
	return c.DB.Set(ctx, key, value, exp).Err()
}

func (c *Client) Get(ctx context.Context, key string) (interface{}, error) {
	res, err := c.DB.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) HSet(ctx context.Context, key string, value interface{}) error {
	exp := time.Duration(c.Duration) * time.Second
	return c.DB.HSet(ctx, key, value, exp).Err()
}

func (c *Client) HGet(ctx context.Context, key, field string) (interface{}, error) {
	res, err := c.DB.HGet(ctx, key, field).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func newRedisDB(params *Params, opts *options) *gredis.Client {
	rdb := gredis.NewClient(&gredis.Options{
		Addr:       params.Address,
		Password:   params.Password,
		MaxRetries: opts.retries,
	})

	return rdb
}
