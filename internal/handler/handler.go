package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meedoed/auth-rest/internal/config"
	"github.com/meedoed/auth-rest/internal/service"
	"github.com/meedoed/auth-rest/pkg/auth"
)

type Handler struct {
	services        *service.Services
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager, accessTokenTTL, refreshTokenTTL time.Duration) *Handler {
	return &Handler{
		services:        services,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		// auth.POST("/sign-up", h.signUp)
		auth.GET("/receive", h.receive)
		auth.GET("/refresh", h.refresh)

	}

	router.GET("ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router
}
