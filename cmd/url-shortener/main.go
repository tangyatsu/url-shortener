package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var flagConfigPath = flag.String("config-path", "./config/local.yaml", "absolute path to config file")

func main() {
	flag.Parse()

	cfg := config.MustLoad(*flagConfigPath) // Done: init config: cleanenv (вместо вайпера или кобры)

	logger := setupLogger(cfg.Env)

	logger.Info("starting", slog.String("env", cfg.Env))
	logger.Debug("currently in debug mode") // Done: init logger: slog (самая акутальная библиотека, но кажется в 1.21 добавилась в станлратную библиотеку)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to initialize storage", sl.Err(err)) // Done: init storage: sqlite чтоб не парится сильно, база данных в файлике (посмотреть видос у николая)
		os.Exit(1)
	}

	router := chi.NewRouter() // DONE: init router: chi, chi render (минамилистичный и полностью совместим с net/http)
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat) // дает удобство, но из-за этого привязываешься в chi (см реализацию хэндлеров)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(logger, storage))
	}) // внутри имеющегося роутера создаем еще роутер, чтоб добавить авторизацию на группу модифицирующих запросов

	router.Get("/{alias}", redirect.New(logger, storage))

	logger.Info("starting server", slog.String("addres", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}

	logger.Error("server has been stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
