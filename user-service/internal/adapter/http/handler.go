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
	router.POST("/user", middleware.ValidateUserPayload(), func(c *gin.Context) {
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
	router.POST("/auth/login", func(c *gin.Context){
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

	//Recuperacion de contraseña
	router.PATCH("/user/password", func(c *gin.Context){
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}
		err := userService.RecoverPassword(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, could not recover password", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Password recovered successfully", "user": &req})
	})

	//Obtencion de un usuario
	router.GET("/user/:id", middleware.AuthMiddleware(), func(c *gin.Context) {
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
	router.PUT("/user", func(c *gin.Context){
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		err := userService.Update(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, could not update user", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": &req})
	})

	//Eliminacion de usuario
	router.DELETE("/user", func(c *gin.Context){
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		err := userService.Delete(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, could not delete user", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "user": &req})
	})

	//Obtencion de todos los usuarios
	router.GET("/users", func(c *gin.Context){
		users, err := userService.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, could not get users", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Users retrieved successfully", "users": users})
	})	
}
