package command

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/sysdevguru/bluelabs/api"
	"github.com/sysdevguru/bluelabs/pkg"
)

type StopFn func(context.Context) error

func Start(ctx context.Context, cfg pkg.Config) StopFn {
	service, err := api.NewService(cfg)
	if err != nil {
		log.Fatalf("could not load dependency %s\n", err.Error())
	}

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      api.NewRouter(service),
		IdleTimeout:  cfg.Server.IdleTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Println("server initialized...")
		if err = server.ListenAndServe(); err != nil {
			log.Fatalf("http.ListenAndServer failed %v\n", err)
		}
	}()

	return func(stopCtx context.Context) error {
		err := server.Shutdown(stopCtx)
		if err != nil {
			log.Println("server failed to shutdown gracefully", err)
			return err
		}

		err = service.Shutdown()
		if err != nil {
			log.Println("service failed to shutdown gracefully", err)
			return err
		}

		log.Println("server has shutdown gracefully")
		return nil
	}
}
