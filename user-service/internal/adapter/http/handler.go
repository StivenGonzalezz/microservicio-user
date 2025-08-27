package http

import (
	"net/http"
	"user-service/internal/domain/model"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userService *service.UserService) {
	router.POST("/register", func(c *gin.Context){
		var req model.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		err := userService.Register(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, could not register user", "error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": &req})
	})

	router.POST("/login", func(c *gin.Context){
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

	router.PUT("/update", func(c *gin.Context){
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

	router.DELETE("/delete", func(c *gin.Context){
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

	router.GET("/user", func(c *gin.Context){
		users, err := userService.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error, could not get users", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Users retrieved successfully", "users": users})
	})	

	router.POST("/recover-password", func(c *gin.Context){
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
}
