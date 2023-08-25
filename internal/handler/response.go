package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/meedoed/auth-rest/pkg/logger"
)

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
