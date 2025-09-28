package main

import (
	"log"
	"os"
	"time"
	_ "user-service/documentation/docs"
	"user-service/internal/adapter/http"
	"user-service/internal/adapter/repository"
	"user-service/internal/service"
	"user-service/pkg/rabbitmq"

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
	publisher := connectRabbitWithRetry("user.events")
	defer publisher.Close()

	repo := repository.NewPostgresRepo()
	user := &service.UserService{Repo: repo, Publisher: publisher}

	router := gin.Default()
	http.SetupRoutes(router, user)

	router.StaticFile("/swagger.json", "documentation/docs/myswagger.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger.json"),
	))

	router.Run(":8080")
}


func connectRabbitWithRetry(exchange string) *rabbitmq.Publisher {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@rabbitmq:5672/"
	}

	var publisher *rabbitmq.Publisher
	var err error

	for i := 1; i <= 10; i++ {
		publisher, err = rabbitmq.NewPublisher(url, exchange)
		if err == nil {
			log.Printf("✅ Conectado a RabbitMQ en intento %d", i)
			return publisher
		}

		log.Printf("⚠️  Intento %d: no se pudo conectar a RabbitMQ: %v", i, err)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("❌ No se pudo conectar a RabbitMQ después de 10 intentos: %v", err)
	return nil
}
