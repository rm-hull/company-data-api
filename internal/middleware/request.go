package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger is a middleware that logs details about every incoming request.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		slog.Info("Request handled",
			"status", status,
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"latency", latency,
			"client_ip", c.ClientIP(),
		)
	}
}
