package logger

import (
	log "log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func GinSlogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Info("HTTP request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"latency", duration.String(),
		)
	}
}
