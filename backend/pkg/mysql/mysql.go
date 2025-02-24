package mysql

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	dmysql "github.com/go-sql-driver/mysql"
	"github.com/quansolashi/golang-boierplate/backend/ent"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	defaultMaxIdleConns    = 10
	defaultMaxOpenConns    = 10
	defaultConnMaxLifetime = 60 * time.Second
	defaultConnMaxIdleTime = 30 * time.Second
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
	debug    bool
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

func WithDebug(debug bool) Option {
	return func(opts *options) {
		opts.debug = debug
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

func NewEntClient(params *Params, opts ...Option) (*ent.Client, error) {
	dopts := &options{
		logger:   zap.NewNop(),
		now:      time.Now,
		location: time.UTC,
	}
	for i := range opts {
		opts[i](dopts)
	}

	cli, err := newEntClient(params, dopts)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func newEntClient(params *Params, opts *options) (*ent.Client, error) {
	var addr string
	switch params.Socket {
	case "tcp":
		addr = fmt.Sprintf("%s:%s", params.Host, params.Port)
	case "unix":
		addr = params.Host
	}
	dsn := dmysql.Config{
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

	drv, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		return nil, err
	}

	var entOptions []ent.Option
	if opts.debug {
		entOptions = append(entOptions, ent.Debug())
	}

	db := drv.DB()
	db.SetMaxIdleConns(defaultMaxIdleConns)
	db.SetMaxOpenConns(defaultMaxOpenConns)
	db.SetConnMaxLifetime(defaultConnMaxLifetime)
	db.SetConnMaxIdleTime(defaultConnMaxIdleTime)

	entOptions = append(entOptions, ent.Driver(drv))
	return ent.NewClient(entOptions...), nil
}
