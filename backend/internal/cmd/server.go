package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	web "github.com/quansolashi/golang-boierplate/backend/internal/web/controller"
	"github.com/quansolashi/golang-boierplate/backend/pkg/config"
	"github.com/quansolashi/golang-boierplate/backend/pkg/http"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type app struct {
	logger zerolog.Logger
	web    web.Controller
	env    *config.Environment
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
	eg.Go(func() error {
		if err := server.Serve(); err != nil {
			app.logger.Warn().AnErr("Failed to run http server", err).Send()
			return nil
		}
		return nil
	})

	app.logger.Info().Int64("port", app.env.Port).Send()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ectx.Done():
		app.logger.Warn().Err(ectx.Err()).Send()
	case signal := <-quit:
		app.logger.Info().Any("Shutdown Server ...", signal).Send()
		delay := time.Duration(5) * time.Second
		time.Sleep(delay)
	}

	if err := server.Stop(ctx); err != nil {
		app.logger.Error().AnErr("Failed to stopped http server", err).Send()
		return err
	}
	return eg.Wait()
}
