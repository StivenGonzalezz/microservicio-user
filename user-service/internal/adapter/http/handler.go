package http

import (
	"net/http"
	"strconv"
	"user-service/internal/domain/model"
	"user-service/internal/middleware"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userService *service.UserService) {
	//Creacion de usuario
	router.POST("/users", middleware.ValidateUserPayload(), func(c *gin.Context) {
		userData, _ := c.Get("user")
		req := userData.(model.User)

		err := userService.Register(&req)
		if err != nil {
			if err.Error() == "email already registered" {
				c.JSON(http.StatusConflict, gin.H{
					"message": "The email is already in use",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error, could not register user",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user":    req,
		})
	})

	//Login de usuario
	router.POST("/login", func(c *gin.Context) {
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}
		token, err := userService.Login(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, could not login user", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": token})
	})

	//Generacion de URL para recuperacion de contraseña
	router.POST("/recovery/password", func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		url, err := userService.RecoverPassword(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error, could not generate recovery link",
				"error":   err.Error(),
			})
			return
		}

		// En el futuro se envía por email, de momento lo devolvemos
		c.JSON(http.StatusOK, gin.H{
			"message": "Recovery link generated successfully",
			"url":     url,
		})
	})

	//Recuperacion de contraseña
	router.PATCH("/users/password/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
			return
		}

		var req struct {
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		err = userService.UpdatePassword(uint(id), req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update password", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	})

	//Obtencion de un usuario
	router.GET("/users/:id", middleware.VerifyToken(), func(c *gin.Context) {
		id := c.Param("id")
		userId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}
		user, err := userService.GetId(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error, could not get user",
				"error":   err.Error(),
			})
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User retrieved successfully",
			"user":    user,
		})
	})

	//Actualizacion de usuario
	router.PUT("/users/:id", middleware.ValidateUserPayload(), middleware.VerifyToken(), func(c *gin.Context) {
		val, _ := c.Get("jwtEmail")
		email := val.(string)

		valUser, _ := c.Get("user")
		req := valUser.(model.User)

		userId, _ := c.Get("jwtId")
		req.ID = uint(userId.(int))

		URLId := c.Param("id")
		URLIdInt, err2 := strconv.Atoi(URLId)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
			return
		}
		if URLIdInt != int(req.ID) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		if email != req.Email {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		err := userService.Update(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, could not update user", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": &req})
	})

	//Eliminacion de usuario
	router.DELETE("/users/:id", middleware.VerifyToken(), func(c *gin.Context) {
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		userId, _ := c.Get("jwtId")
		if userId == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		URLId := c.Param("id")
		URLIdInt, err2 := strconv.Atoi(URLId)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
			return
		}
		if URLIdInt != int(userId.(int)) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		err := userService.Delete(userId.(int), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, could not delete user", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "user": &req})
	})

	//Obtencion de todos los usuarios
	router.GET("/users", middleware.VerifyToken(), func(c *gin.Context) {
		email, _ := c.Get("jwtEmail")
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}
		nameOrEmail := req.Name
		if email != "admin@admin.com" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		users, err := userService.GetAll(nameOrEmail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, could not get users", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Users retrieved successfully", "users": users})
	})
}
