package config

import "github.com/kelseyhightower/envconfig"

type Environment struct {
	Port       int64  `envconfig:"PORT" default:"8080"`
	DBSocket   string `envconfig:"DB_SOCKET" default:"tcp"`
	DBHost     string `envconfig:"DB_HOST" default:"127.0.0.1"`
	DBPort     string `envconfig:"DB_PORT" default:"3306"`
	DBDatabase string `envconfig:"DB_DATABASE" default:""`
	DBUsername string `envconfig:"DB_USERNAME" default:"root"`
	DBPassword string `envconfig:"DB_PASSWORD" default:""`
}

type Client interface {
	ProcessEnv(prefix string, spec interface{}) error
}

type client struct{}

func NewClient() Client {
	return &client{}
}

func (c *client) ProcessEnv(prefix string, spec interface{}) error {
	return envconfig.Process(prefix, spec)
}
