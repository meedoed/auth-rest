package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meedoed/auth-rest/internal/config"
	"github.com/meedoed/auth-rest/internal/server"
	"github.com/meedoed/auth-rest/pkg/auth"
	"github.com/meedoed/auth-rest/pkg/database/mongodb"
	"github.com/meedoed/auth-rest/pkg/hash"
	"github.com/meedoed/auth-rest/pkg/logger"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	// Dependencies
	mongoClient, err := mongodb.ConnectToDatabase(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logger.Error(err)
		return
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}

	handlers := delivery.NewHandler(services, tokenManager)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err.Error())
	}
}
