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

	"golang.org/x/sync/errgroup"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := NewRouter()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 3002),
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
