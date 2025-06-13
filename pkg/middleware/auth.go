package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/pkg/jwt"
	"net/http"
	"strings"
)

func writeUnAuthed(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": http.StatusText(http.StatusUnauthorized)})
}

func IsAuthed(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnAuthed(c)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, data := jwt.NewJWT(jwtSecret).Parse(token)
		if !isValid {
			writeUnAuthed(c)
			return
		}

		c.Set("user_id", data.UserId)

		c.Next()
	}
}
