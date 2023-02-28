package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"

	"github.com/AntonBraer/urlShorter/config"
	v1 "github.com/AntonBraer/urlShorter/internal/controller/http/v1"
	"github.com/AntonBraer/urlShorter/internal/usecase"
	"github.com/AntonBraer/urlShorter/internal/usecase/repo"
	"github.com/AntonBraer/urlShorter/pkg/httpserver"
	"github.com/AntonBraer/urlShorter/pkg/logger"
	"github.com/AntonBraer/urlShorter/pkg/postgres"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Connect to database
	url := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		cfg.PG.User, cfg.PG.Password, cfg.PG.ServiceName, cfg.PG.DB)
	pg, err := postgres.New(url, l, context.Background())
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Pool.Close()

	linkUseCase := usecase.NewLinkUseCase(repo.NewLinkRepo(pg))
	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, linkUseCase)
	httpServer := httpserver.New(handler, cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-interrupt
		l.Info("app - Run - signal: syscall.SIGINT")

		if err = httpServer.Shutdown(); err != nil {
			l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}()

	// Waiting for errors
	if err = httpServer.NotifyError(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.NotifyError: %w", err))
	}
}
