package mysql

import (
	"fmt"
	"time"

	dmysql "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Client struct {
	DB *gorm.DB
}

type Params struct {
	Socket   string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type options struct {
	logger   *zap.Logger
	now      func() time.Time
	location *time.Location
}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithNow(now func() time.Time) Option {
	return func(opts *options) {
		opts.now = now
	}
}

func WithLocation(location *time.Location) Option {
	return func(opts *options) {
		opts.location = location
	}
}

func NewClient(params *Params, opts ...Option) (*Client, error) {
	dopts := &options{
		logger:   zap.NewNop(),
		now:      time.Now,
		location: time.UTC,
	}
	for i := range opts {
		opts[i](dopts)
	}

	db, err := newDBClient(params, dopts)
	if err != nil {
		return nil, err
	}

	c := &Client{
		DB: db,
	}
	return c, nil
}

func newDBClient(params *Params, opts *options) (*gorm.DB, error) {
	conf := &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{
			SingularTable: true,
		},
		NowFunc: opts.now,
	}
	var addr string
	switch params.Socket {
	case "tcp":
		addr = fmt.Sprintf("%s:%s", params.Host, params.Port)
	case "unix":
		addr = params.Host
	}
	dsn := &dmysql.Config{
		User:                 params.Username,
		Passwd:               params.Password,
		Net:                  params.Socket,
		Addr:                 addr,
		DBName:               params.Database,
		Loc:                  opts.location,
		ParseTime:            true,
		Collation:            "utf8mb4_general_ci",
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		Params:               map[string]string{"charset": "utf8mb4"},
		MaxAllowedPacket:     4194304, // 4MiB
	}
	return gorm.Open(mysql.Open(dsn.FormatDSN()), conf)
}
