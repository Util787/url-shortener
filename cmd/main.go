package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Util787/url-shortener/internal/adapters/rest"
	"github.com/Util787/url-shortener/internal/adapters/storage"
	"github.com/Util787/url-shortener/internal/adapters/tgbot"
	"github.com/Util787/url-shortener/internal/config"
	"github.com/Util787/url-shortener/internal/shortener-usecase"
)

func main() {
	cfg := config.MustLoadConfig()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	pgStorage := storage.MustInitPostgres(context.Background(), cfg.PostgresConfig)

	ShortenerUsecase := shortener.NewShortenerUsecase(&pgStorage, cfg.RedirectBaseURL)
	// Start Telegram bot if token provided
	if cfg.TgBotConfig.Token != "" {
		bot, err := tgbot.NewBot(cfg.TgBotConfig.Token, ShortenerUsecase)
		if err != nil {
			log.Error("failed to init telegram bot", slog.String("error", err.Error()))
		} else {
			go func() {
				log.Info("Telegram bot started")
				bot.Start()
			}()
		}
	}

	server := rest.NewRestServer(log, cfg.HTTPServerConfig, ShortenerUsecase, cfg.RedirectBaseURL)

	go func() {
		log.Info("HTTP server start", slog.String("host", cfg.HTTPServerConfig.Host), slog.Int("port", cfg.HTTPServerConfig.Port))
		if err := server.Run(); err != nil {
			log.Error("HTTP server error", slog.String("error", err.Error()))
		}
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Info("Shutting down gracefully...")

	log.Info("Shutting down server")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error("HTTP server shutdown error", slog.String("error", err.Error()))
	}

	log.Info("Shutting down postgres")
	pgStorage.Shutdown()

	log.Info("Shutdown complete")

}
