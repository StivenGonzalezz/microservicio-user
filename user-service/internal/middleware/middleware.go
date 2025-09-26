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
		// Leer el cuerpo de la petición
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Invalid request body"})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var raw map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &raw); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Invalid JSON format"})
			c.Abort()
			return
		}

		allowedFields := map[string]bool{
			"name":     true,
			"lastName": true,
			"email":    true,
			"password": true,
		}

		requiredFields := []string{"name", "lastName", "email", "password"}

		for key := range raw {
			if !allowedFields[key] {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "400",
					"error":   "Bad Request",
					"message": "Unexpected field: " + key,
				})
				c.Abort()
				return
			}
		}

		for _, field := range requiredFields {
			if value, exists := raw[field]; !exists || value == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "400",
					"error":   "Bad Request",
					"message": "Field '" + field + "' is required and cannot be empty",
				})
				c.Abort()
				return
			}
		}

		if email, ok := raw["email"].(string); ok {
			if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "400",
					"error":   "Bad Request",
					"message": "Invalid email format",
				})
				c.Abort()
				return
			}
		}

		if password, ok := raw["password"].(string); ok {
			if len(password) < 6 {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "400",
					"error":   "Bad Request",
					"message": "Password must be at least 6 characters long",
				})
				c.Abort()
				return
			}
		}

		var user model.User
		if err := json.Unmarshal(bodyBytes, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "400",
				"error":   "Bad Request",
				"message": "Invalid user data: " + err.Error(),
			})
			c.Abort()
			return
		}

		if strings.TrimSpace(user.Name) == "" || strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "400",
				"error":   "Bad Request",
				"message": "Fields cannot be empty or contain only whitespace",
			})
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
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401","error": "Unauthorized","message": "Missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401","error": "Unauthorized","message": "Invalid Authorization header"})
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401","error": "Unauthorized","message": "Invalid or expired token"})
			c.Abort()
			return
		}
		if email, ok := claims["user_email"].(string); ok {
			c.Set("jwtEmail", email)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401","error": "Unauthorized","message": "Token does not contain user_email"})
			c.Abort()
			return
		}

		if id, ok := claims["user_id"].(float64); ok {
			c.Set("jwtId", int(id))
		}

		c.Next()
	}
}