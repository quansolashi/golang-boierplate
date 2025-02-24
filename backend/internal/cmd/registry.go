package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/quansolashi/golang-boierplate/backend/ent"
	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	"github.com/quansolashi/golang-boierplate/backend/internal/database/mysql"
	graph "github.com/quansolashi/golang-boierplate/backend/internal/graphql/handler"
	web "github.com/quansolashi/golang-boierplate/backend/internal/web/controller"
	"github.com/quansolashi/golang-boierplate/backend/pkg/config"
	"github.com/quansolashi/golang-boierplate/backend/pkg/log"
	pmysql "github.com/quansolashi/golang-boierplate/backend/pkg/mysql"
	"github.com/quansolashi/golang-boierplate/backend/pkg/rabbitmq"
	"github.com/quansolashi/golang-boierplate/backend/pkg/redis"
	"github.com/rs/zerolog"
)

func (a *app) inject(ctx context.Context) error {
	// temporary print ctx to avoid rules by lint
	fmt.Println(ctx.Value(""))

	// load environment variables
	env, err := a.loadEnv()
	if err != nil {
		return err
	}
	a.env = env

	logger, err := log.NewLogger(
		log.WithLevel(zerolog.DebugLevel),
	)
	if err != nil {
		return err
	}
	a.logger = logger

	// connect and initialize database
	database, err := a.newDatabase()
	if err != nil {
		return err
	}

	redis := a.newRedisDatabase()

	rabbitmq, err := a.newRabbitMQ()
	if err != nil {
		return err
	}
	a.queue = rabbitmq

	// app web controller
	a.web = web.NewController(&web.Params{
		DB:               database,
		Redis:            redis,
		RabbitMQ:         rabbitmq,
		LocalTokenSecret: a.env.LocalTokenSecret,
		WebURL:           a.env.WebURL,
		GoogleAPIKey:     a.env.GoogleAPIKey,
		GoogleAPISecret:  a.env.GoogleAPISecret,
	})

	// graphql
	ent, err := a.newEntClient()
	if err != nil {
		return err
	}
	a.graph = graph.NewGraph(&graph.Params{
		LocalTokenSecret: a.env.LocalTokenSecret,
		Ent:              ent,
	})

	return nil
}

func (a *app) loadEnv() (*config.Environment, error) {
	env := &config.Environment{}

	config := config.NewClient()
	err := config.ProcessEnv("", env)

	if err != nil {
		return nil, err
	}
	return env, nil
}

func (a *app) newDatabase() (*database.Database, error) {
	params := &pmysql.Params{
		Socket:   a.env.DBSocket,
		Host:     a.env.DBHost,
		Port:     a.env.DBPort,
		Database: a.env.DBDatabase,
		Username: a.env.DBUsername,
		Password: a.env.DBPassword,
	}
	opts := []pmysql.Option{
		pmysql.WithNow(time.Now),
		pmysql.WithLocation(time.Now().Location()),
	}
	client, err := pmysql.NewClient(params, opts...)
	if err != nil {
		return nil, err
	}
	return mysql.NewDatabase(client), nil
}

func (a *app) newRedisDatabase() *redis.Client {
	params := &redis.Params{
		Address:  fmt.Sprintf("%s:%d", a.env.RedisDBHost, a.env.RedisDBPort),
		Password: a.env.RedisDBPassword,
	}
	opts := []redis.Option{
		redis.WithMaxRetries(3),
	}
	client := redis.NewClient(params, opts...)
	return client
}

func (a *app) newRabbitMQ() (rabbitmq.Client, error) {
	params := &rabbitmq.Params{
		Host:     a.env.RabbitMQHost,
		Port:     a.env.RabbitMQPort,
		Username: a.env.RabbitMQUsername,
		Password: a.env.RabbitMQPassword,
	}
	logger, err := log.NewLogger(
		log.WithLevel(zerolog.DebugLevel),
	)
	if err != nil {
		return nil, err
	}
	opts := []rabbitmq.Option{
		rabbitmq.WithLogger(logger),
	}
	client, err := rabbitmq.NewRabbitMQ(params, opts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (a *app) newEntClient() (*ent.Client, error) {
	params := &pmysql.Params{
		Socket:   a.env.DBSocket,
		Host:     a.env.DBHost,
		Port:     a.env.DBPort,
		Database: a.env.DBDatabase,
		Username: a.env.DBUsername,
		Password: a.env.DBPassword,
	}
	opts := []pmysql.Option{
		pmysql.WithNow(time.Now),
		pmysql.WithLocation(time.Now().Location()),
	}
	client, err := pmysql.NewEntClient(params, opts...)
	if err != nil {
		return nil, err
	}
	return client, nil
}
