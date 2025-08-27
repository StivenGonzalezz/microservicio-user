package main

import (
	"user-service/internal/adapter/repository"
	"user-service/internal/adapter/http"
	"user-service/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables del sistema si existen.")
	}

}

func main() {
	repo := repository.NewPostgresRepo()
	user := &service.UserService{Repo: repo}

	router := gin.Default()
	http.SetupRoutes(router, user)
	router.Run(":8080")
}
