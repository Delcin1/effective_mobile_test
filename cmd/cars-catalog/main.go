package main

import (
	"context"
	_ "effective_mobile_test/docs" // docs is generated by Swag CLI, you have to import it.
	"effective_mobile_test/internal/config"
	carDelete "effective_mobile_test/internal/http-server/handlers/car/delete"
	carSave "effective_mobile_test/internal/http-server/handlers/car/save"
	carSearch "effective_mobile_test/internal/http-server/handlers/car/search"
	carUpdate "effective_mobile_test/internal/http-server/handlers/car/update"
	ownerDelete "effective_mobile_test/internal/http-server/handlers/owner/delete"
	ownerSave "effective_mobile_test/internal/http-server/handlers/owner/save"
	ownerUpdate "effective_mobile_test/internal/http-server/handlers/owner/update"
	mwLogger "effective_mobile_test/internal/http-server/middleware/logger"
	"effective_mobile_test/internal/lib/logger/sl"
	"effective_mobile_test/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

//	@title			Cars Catalog API
//	@version		1.0
//	@description	API for Effective Mobile Test
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Nikita Zhirnov
//	@contact.url	https://t.me/belkindelcin
//	@contact.email	naklz9@mail.ru

func main() {
	cfg := config.InitConfig()

	log := setupLogger(cfg.Env)

	log.Info("starting cars-catalog api", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)

	router.Post("/car/save", carSave.New(log, storage, cfg.HelpAPI))
	router.Get("/car/search", carSearch.New(log, storage))
	router.Delete("/car/delete", carDelete.New(log, storage))
	router.Put("/car/update", carUpdate.New(log, storage))

	router.Post("/owner/save", ownerSave.New(log, storage))
	router.Delete("/owner/delete", ownerDelete.New(log, storage))
	router.Put("/owner/update", ownerUpdate.New(log, storage))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"), //The url pointing to API definition
	))

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
