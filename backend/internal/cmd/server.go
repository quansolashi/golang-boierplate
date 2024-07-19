package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	web "github.com/quansolashi/message-extractor/backend/internal/web/controller"
	"github.com/quansolashi/message-extractor/backend/pkg/config"
	"github.com/quansolashi/message-extractor/backend/pkg/http"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type app struct {
	logger *zap.Logger
	web    web.Controller
	env    *config.Environment
}

type environment struct {
	Port int64 `envconfig:"PORT" default:"8080"`
}

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := &app{}
	// inject dependencies
	if err := app.inject(ctx); err != nil {
		return fmt.Errorf("cmd: failed to new registry: %w", err)
	}

	// initialize router and http server with port
	router := app.newRouter()
	server := http.NewHTTPServer(router.Handler(), app.env.Port)

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if err = server.Serve(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
		return
	})

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ectx.Done():
		log.Println("Done context")
	case signal := <-quit:
		log.Println("Shutdown Server ...", signal)
		delay := time.Duration(5) * time.Second
		time.Sleep(delay)
	}

	if err := server.Stop(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
		return err
	}
	return eg.Wait()
}
