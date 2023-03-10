package api

import (
	"github.com/NeverovDS/playlist/logger"

	"github.com/gin-gonic/gin"
)

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	logger.Errorf(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
