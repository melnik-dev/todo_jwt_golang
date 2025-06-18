package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/pkg/jwt"
	"github.com/melnik-dev/go_todo_jwt/pkg/logger"
	"github.com/melnik-dev/go_todo_jwt/pkg/response"
	"net/http"
	"strings"
)

func IsAuthed(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auLogger := logger.FromContext(c)

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			auLogger.Warn("Unauthorized: Authorization header no Bearer prefix")
			response.AbortWithStatus(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, data := jwt.NewJWT(jwtSecret).Parse(token)
		if !isValid {
			auLogger.Warn("Unauthorized: Invalid or expired JWT token provided")
			response.AbortWithStatus(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		auLogger.WithField("user_id", data.UserId).Debug("User authenticated successfully")
		c.Set("user_id", data.UserId)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (int, bool) {
	auLogger := logger.FromContext(c)

	strId, ok := c.Get("user_id")
	if !ok {
		auLogger.Warn("User ID not found in context")
		response.InternalServerError(c, "User ID not found in context")
		return 0, false
	}
	userID, ok := strId.(int)
	if !ok {
		auLogger.Error("Invalid user ID type")
		response.InternalServerError(c, "Invalid user ID type")
		return 0, false
	}
	auLogger.WithField("user_id", userID)
	return userID, true
}
