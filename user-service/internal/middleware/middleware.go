package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"user-service/internal/domain/model"
	"user-service/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func ValidateUserPayload() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var raw map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &raw); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
			c.Abort()
			return
		}

		allowed := map[string]bool{
			"name":     true,
			"lastName": true,
			"email":    true,
			"password": true,
		}

		for key := range raw {
			if !allowed[key] {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Unexpected field: " + key})
				c.Abort()
				return
			}
		}

		var user model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user data", "error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing Authorization header"})
			c.Abort()
			return
		}

		// Formato: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header"})
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Guardamos el email en contexto
		if email, ok := claims["user_email"].(string); ok {
			c.Set("jwtEmail", email)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token does not contain user_email"})
			c.Abort()
			return
		}

		c.Next()
	}
}