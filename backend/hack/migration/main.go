// データベースのマイグレーション実行
// go run ./hack/migration/main.go -db-socket='tcp' -db-port='3306' -db-name='msx-mysql' -db-username='root' -db-password='password'
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/quansolashi/message-extractor/backend/internal/database/migrate"
	"github.com/quansolashi/message-extractor/backend/pkg/mysql"
	"go.uber.org/zap"
)

var (
	dbsocket   string
	dbhost     string
	dbport     string
	dbname     string
	dbusername string
	dbpassword string
)

func main() {
	startedAt := time.Now()
	if err := run(); err != nil {
		panic(err)
	}
	const format = "2006-01-02 15:04:05"
	fmt.Printf("Done: %s -> %s\n", startedAt.Format(format), time.Now().Format(format))
}

func run() error {
	// from arguments
	flag.StringVar(&dbsocket, "db-socket", "tcp", "mysql server protocol")
	flag.StringVar(&dbhost, "db-host", "mysql", "mysql server host")
	flag.StringVar(&dbport, "db-port", "3306", "mysql server port")
	flag.StringVar(&dbname, "db-name", "msx-mysql", "mysql database name")
	flag.StringVar(&dbusername, "db-username", "root", "mysql auth username")
	flag.StringVar(&dbpassword, "db-password", "root", "mysql auth password")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	db, err := newDBClient(dbsocket, dbhost, dbport, dbname, dbusername, dbpassword, &zap.Logger{})
	if err != nil {
		return err
	}

	log.Println("database migration will begin")
	if err := migrate.Run(ctx, db, &zap.Logger{}); err != nil {
		return err
	}
	log.Println("database migration has been completed")
	return nil
}

func newDBClient(socket, host, port, database, username, password string, logger *zap.Logger) (*mysql.Client, error) {
	params := &mysql.Params{
		Socket:   socket,
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}
	return mysql.NewClient(params, mysql.WithLocation(time.Now().Location()), mysql.WithLogger(logger))
}
