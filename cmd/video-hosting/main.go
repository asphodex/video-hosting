package main

import (
	"log/slog"
	"net/http"
	"os"
	"video-hosting/internal/config"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	//http.HandleFunc("/video", getVideo)у
	cfg := config.MustLoad()
	//fmt.Println(cfg)

	log := setupLogger(cfg.Env)

	log.Info("LOG: starting server", slog.String("env", cfg.Env))
	log.Debug("LOG: debug mode enabled")



	if err := http.ListenAndServe(":"+"8000", nil); err != nil {
		panic(err)
	}
}

func setupLogger(env string)  *slog.Logger {
	var log *slog.Logger
	switch env {
		case envLocal:
			log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case envDev:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		case envProd:
			log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
