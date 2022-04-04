package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sysdevguru/bluelabs/cmd/web/command"
	"github.com/sysdevguru/bluelabs/pkg"
)

func main() {
	ctx := context.Background()
	cfg, err := pkg.Load()
	if err != nil {
		log.Fatalf("could not load configuration %s\n", err.Error())
	}
	stopFn := command.Start(ctx, cfg)

	sig := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal.
	<-sig

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if stopFn(ctx) != nil {
		os.Exit(1)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	os.Exit(0)
}
