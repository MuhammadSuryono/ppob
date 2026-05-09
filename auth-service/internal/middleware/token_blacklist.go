package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yontech/ppob/auth-service/internal/services"
)

func TokenBlacklistMiddleware(blacklistService *services.TokenBlacklistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, nil)
		if err != nil {
			c.Next()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if jti, ok := claims["jti"].(string); ok && jti != "" {
				isBlacklisted, err := blacklistService.IsAccessTokenBlacklisted(c.Request.Context(), jti)
				if err == nil && isBlacklisted {
					c.JSON(http.StatusUnauthorized, gin.H{
						"error":   "token_revoked",
						"message": "Token has been revoked",
					})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}