package cmd

import (
	"context"
	"time"

	"github.com/quansolashi/message-extractor/backend/internal/database"
	"github.com/quansolashi/message-extractor/backend/internal/database/mysql"
	web "github.com/quansolashi/message-extractor/backend/internal/web/controller"
	"github.com/quansolashi/message-extractor/backend/pkg/config"
	pmysql "github.com/quansolashi/message-extractor/backend/pkg/mysql"
)

func (a *app) inject(ctx context.Context) error {
	// load environment variables
	env, err := a.loadEnv()
	if err != nil {
		return err
	}
	a.env = env

	// connect and initialize database
	database, err := a.newDatabase()
	if err != nil {
		return err
	}
	// app web controller
	a.web = web.NewController(&web.Params{
		DB: database,
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
	// if err := client.DB.Use(telemetry.NewNrTracer(a.DBDatabase, a.DBHost, string(newrelic.DatastoreMySQL))); err != nil {
	// 	return nil, err
	// }
	return mysql.NewDatabase(client), nil
}
