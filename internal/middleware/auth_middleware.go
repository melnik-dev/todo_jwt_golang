package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func ValidateToken(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "Token required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		decoded, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !decoded.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "Invalid or expired token"})
			return
		}

		claims, ok := decoded.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "Invalid token claims"})
			return
		}

		userID, ok := claims["user_id"].(float64) // JWT числа — float64
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "User unauthorized"})
			return
		}

		c.Set("user_id", int(userID))

		c.Next()
	}
}
