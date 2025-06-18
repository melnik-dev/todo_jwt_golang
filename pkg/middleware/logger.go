package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/pkg/logger"
	"github.com/sirupsen/logrus"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		requestID := requestid.Get(c)
		if requestID == "" {
			requestID = "no-request-id" // Fallback
		}

		requestLogger := logger.GetLogger().WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		c.Set("logger", requestLogger)

		c.Next()

		// детали после обработки запроса
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logFields := logrus.Fields{
			"status":   statusCode,
			"latency":  latency.String(),
			"response": c.Writer.Size(),
			"error":    c.Errors.ByType(gin.ErrorTypePrivate).String(), // Логируем ошибки Gin
		}

		if raw != "" {
			path = path + "?" + raw
		}

		if statusCode >= 500 {
			requestLogger.WithFields(logFields).Error("Request completed with server error")
		} else if statusCode >= 400 {
			requestLogger.WithFields(logFields).Warn("Request completed with client error")
		} else {
			requestLogger.WithFields(logFields).Info("Request completed successfully")
		}
	}
}
