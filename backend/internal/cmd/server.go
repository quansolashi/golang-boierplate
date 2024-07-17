package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/sync/errgroup"
)

type Environment struct {
	Port int64 `envconfig:"PORT" default:"8080"`
}

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := NewRouter()

	var env Environment
	if err := envconfig.Process("", &env); err != nil {
		return fmt.Errorf("cmd: failed to load environment: %w", err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", env.Port),
		Handler: router.Handler(),
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		return
	})

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

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
		return err
	}
	return eg.Wait()
}
