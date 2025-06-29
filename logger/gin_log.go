package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

func GinLogrusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		Log.WithFields(map[string]interface{}{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"latency":    duration.String(),
		}).Info("HTTP request")
	}
}
