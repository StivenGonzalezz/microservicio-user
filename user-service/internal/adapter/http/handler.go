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
				c.JSON(http.StatusConflict, gin.H{"status": "409", "error": "Conflict", "message": "The email is already in use"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"status": "500", "error": "Internal Server Error", "message": "Internal server error, could not register user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "201", "message": "User created successfully", "user": req})
	})

	//Login de usuario
	router.POST("/login", func(c *gin.Context) {
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Invalid request"})
			return
		}
		token, err := userService.Login(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "404", "error": "Not Found", "message": "User not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	//Generacion de URL para recuperacion de contraseña
	router.POST("/recovery", func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Invalid request"})
			return
		}

		url, err := userService.RecoverPassword(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "500", "error": "Internal Server Error", "message": "Internal server error, could not generate recovery link"})
			return
		}

		// En el futuro se envia por email, de momento lo devolvemos por JSON
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
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Invalid user ID"})
			return
		}

		var req struct {
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Invalid request"})
			return
		}

		if req.Password ==""{
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Password cannot be empty"})
			return
		}
		err = userService.UpdatePassword(uint(id), req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "500", "error": "Internal Server Error", "message": "Could not update password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "200", "message": "Password updated successfully"})
	})

	//Obtencion de un usuario
	router.GET("/users/:id", middleware.VerifyToken(), func(c *gin.Context) {
		id := c.Param("id")
		userId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Invalid request"})
			return
		}
		user, err := userService.GetId(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "500", "error": "Internal Server Error", "message": "Internal server error, could not get user"})
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "404", "error": "Not Found", "message": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "200", "message": "User retrieved successfully", "user": user})
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
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Invalid user ID"})
			return
		}
		if URLIdInt != int(req.ID) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401", "error": "Unauthorized", "message": "Unauthorized"})
			return
		}

		if email != req.Email {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401", "error": "Unauthorized", "message": "Unauthorized"})
			return
		}

		err := userService.Update(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "500", "error": "Internal Server Error", "message": "Internal server error, could not update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "200", "message": "User updated successfully", "user": &req})
	})

	//Eliminacion de usuario
	router.DELETE("/users/:id", middleware.VerifyToken(), func(c *gin.Context) {
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Invalid request"})
			return
		}

		userId, _ := c.Get("jwtId")
		if userId == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401", "error": "Unauthorized", "message": "Unauthorized"})
			return
		}

		URLId := c.Param("id")
		URLIdInt, err2 := strconv.Atoi(URLId)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "400", "error": "Bad Request", "message": "Invalid user ID"})
			return
		}
		if URLIdInt != int(userId.(int)) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401", "error": "Unauthorized", "message": "Unauthorized"})
			return
		}

		err := userService.Delete(userId.(int), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "500", "error": "Internal Server Error", "message": "Internal server error, could not delete user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "200", "message": "User deleted successfully", "user": &req})
	})

	// Obtención de usuarios con paginación y filtrado
	router.GET("/users", middleware.VerifyToken(), func(c *gin.Context) {
		email, _ := c.Get("jwtEmail")
		if email != "admin@admin.com" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "401", "error": "Unauthorized", "message": "Only admin can access this endpoint"})
			return
		}

		name := c.Query("name")
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit < 1 || limit > 100 {
			limit = 10
		}

		sort := c.DefaultQuery("sort", "asc")
		if sort != "asc" && sort != "desc" {
			sort = "asc"
		}

		result, err := userService.GetUsersWithPagination(name, page, limit, sort)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "500", "error": "Internal Server Error", "message": "Internal server error, could not get users with pagination"})
			return
		}

		c.JSON(http.StatusOK, result)
	})
}
