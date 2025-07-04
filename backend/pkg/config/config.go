package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	Port                  int64  `envconfig:"PORT" default:"8080"`
	DBSocket              string `envconfig:"DB_SOCKET" default:"tcp"`
	DBHost                string `envconfig:"DB_HOST" default:"127.0.0.1"`
	DBPort                string `envconfig:"DB_PORT" default:"3306"`
	DBDatabase            string `envconfig:"DB_DATABASE" default:""`
	DBUsername            string `envconfig:"DB_USERNAME" default:"root"`
	DBPassword            string `envconfig:"DB_PASSWORD" default:""`
	WebURL                string `envconfig:"WEB_URL" default:""`
	LocalTokenSecret      string `envconfig:"LOCAL_TOKEN_SECRET" default:""`
	PublicTokenPublicKey  string `envconfig:"PUBLIC_TOKEN_PUBLIC_KEY" default:""`
	PublicTokenPrivateKey string `envconfig:"PUBLIC_TOKEN_PRIVATE_KEY" default:""`
	GoogleAPIKey          string `envconfig:"GOOGLE_API_KEY" default:""`
	GoogleAPISecret       string `envconfig:"GOOGLE_API_SECRET" default:""`
	RedisDBHost           string `envconfig:"REDIS_DB_HOST" default:"127.0.0.1"`
	RedisDBPort           int64  `envconfig:"REDIS_DB_PORT" default:"6379"`
	RedisDBPassword       string `envconfig:"REDIS_DB_PASSWORD" default:""`
	RabbitMQUsername      string `envconfig:"RBM_USERNAME" default:"rabbitmq"`
	RabbitMQPassword      string `envconfig:"RBM_PASSWORD" default:""`
	RabbitMQHost          string `envconfig:"RBM_HOST" default:"127.0.0.1"`
	RabbitMQPort          string `envconfig:"RBM_PORT" default:"5672"`
	S3BucketName          string `envconfig:"S3_BUCKET_NAME" default:""`
	AwsRegion             string `envconfig:"AWS_REGION" default:""`
	AwsAccessKey          string `envconfig:"AWS_ACCESS_KEY" default:""`
	AwsSecretKey          string `envconfig:"AWS_SECRET_KEY" default:""`
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
