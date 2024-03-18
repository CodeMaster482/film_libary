package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"films_library/config"
	actorDelivery "films_library/internal/actor/delivery/http"
	actorRep "films_library/internal/actor/repository/postgresql"
	actorUsecase "films_library/internal/actor/usecase"
	filmDelivery "films_library/internal/film/delivery/http"
	filmRep "films_library/internal/film/repository/postgresql"
	filmUsecase "films_library/internal/film/usecase"
	"films_library/internal/middlware"
	"films_library/pkg/httpserver"
	"films_library/pkg/logger"
	"films_library/pkg/postgres"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(
		cfg.PG.Host,
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.Name,
		cfg.PG.Port,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Usecase
	actorRepo := actorRep.NewRepository(pg.Pool)
	actorUsecase := actorUsecase.NewActorUsecase(actorRepo, l)

	filmRepo := filmRep.NewRepository(pg.Pool)
	filmUsecase := filmUsecase.NewFilmUsecase(filmRepo, l)

	// Middleware

	recoveryMW := middlware.NewRecoveryMiddleware(l)
	logMW := middlware.NewLoggingMiddleware(l)

	// HTTP Server
	mux := http.NewServeMux()

	filmDelivery.NewFilmHandler(mux, filmUsecase, l)
	actorDelivery.NewActorHandler(mux, actorUsecase, l)

	r := recoveryMW.Recoverer(mux)
	r = logMW.LoggingMiddleware(r)
	r = middlware.AllowedMethod(r)
	r = middlware.Authentication(r)

	httpServer := httpserver.New(r, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
