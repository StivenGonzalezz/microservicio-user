package main

import (
	"log"
	_ "user-service/documentation/docs"
	"user-service/internal/adapter/http"
	"user-service/internal/adapter/repository"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	//router.StaticFile("/swagger.json", "documentation/docs/myswagger.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger.json"),
	))

	router.Run(":8080")
}
