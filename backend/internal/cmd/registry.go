package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/quansolashi/golang-boierplate/backend/internal/database"
	"github.com/quansolashi/golang-boierplate/backend/internal/database/mysql"
	web "github.com/quansolashi/golang-boierplate/backend/internal/web/controller"
	"github.com/quansolashi/golang-boierplate/backend/pkg/config"
	"github.com/quansolashi/golang-boierplate/backend/pkg/log"
	pmysql "github.com/quansolashi/golang-boierplate/backend/pkg/mysql"
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
	// app web controller
	a.web = web.NewController(&web.Params{
		DB:               database,
		LocalTokenSecret: a.env.LocalTokenSecret,
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
