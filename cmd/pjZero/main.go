package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"main/internal/config"
	"main/internal/converter"
	"main/internal/http-server/handlers/video/download"
	"main/internal/http-server/handlers/video/upload"
	mwLogger "main/internal/http-server/middleware/logger"
	"main/internal/lib/logger/el"
	"main/internal/provider"
	"main/internal/repository"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//TODO: init config: cleanenv

	cfg := config.MustLoad()
	fmt.Println(cfg)
	var connDbStr = fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable", cfg.StoragePath.DbType, cfg.StoragePath.DbUser, cfg.StoragePath.DbUserPassword, cfg.StoragePath.DbPort, cfg.StoragePath.DbName)
	//TODO: init logger: slog

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))                          // к каждому сообщению будет добавляться поле с информацией о текущем окружении
	log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")

	//TODO: inti storage: PostgreSQL

	pg, err := provider.NewPgConn(connDbStr)
	if err != nil {
		log.Error("failed to init storage: ", el.Err(err))
		os.Exit(1)
	}
	defer pg.Close()
	//defer func(pg *sqlx.DB) {
	//	err := pg.Close()
	//	if err != nil {
	//
	//	}
	//}(pg)
	pgRepo := repository.NewPgRepo(pg)
	fmt.Println(pgRepo)

	//TODO: init converter
	conv := converter.NewConv()

	//TODO: init router: chi, "chi render"
	router := chi.NewRouter()

	//проверка что случилось у конкретного запроса
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(mwLogger.New(log))

	router.Post("/add-file", upload.UploadVid(log, conv, pgRepo))
	router.Post("/get-file", download.DownloadVid(log, conv, pgRepo))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped!!")
	//TODO: init server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
