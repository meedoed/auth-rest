package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/meedoed/auth-rest/internal/config"
	"github.com/meedoed/auth-rest/internal/service"
	"github.com/meedoed/auth-rest/pkg/auth"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Inint(cfg *config.Config) *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sigh-up", h.signUp)
		auth.GET("/sign-in", h.signIn)
		auth.GET("/refresh", h.refresh)
	}

	return router
}
