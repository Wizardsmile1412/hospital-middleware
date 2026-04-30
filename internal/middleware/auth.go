package middleware

import (
	"net/http"
	"strings"

	"github.com/Wizardsmile1412/hospital-middleware/internal/token"
	"github.com/gin-gonic/gin"
)

const ClaimsKey = "claims"

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := token.ParseAccessToken(tokenStr, jwtSecret)
		if err != nil {
			msg := "invalid token"
			if err == token.ErrTokenExpired {
				msg = "token has expired"
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		c.Set(ClaimsKey, claims)
		c.Next()
	}
}
