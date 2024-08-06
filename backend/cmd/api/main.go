package main

import (
	"log"
	"sasa-elterminali-service/cmd/api/config"
	_ "sasa-elterminali-service/docs"
	"sasa-elterminali-service/internal/adapters"
	"sasa-elterminali-service/internal/adapters/postgres"
	"sasa-elterminali-service/internal/adapters/repository"

	"github.com/gofiber/fiber/v2"
)

// @title Sasa El Terminali API
// @version 1.0
// @description This is a sample server for a microservice.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:45013
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer' followed by a space and then your token."
func main() {
	// Yapılandırmayı yükle
	config.LoadConfig()

	app := fiber.New()

	// Veritabanı bağlantısını başlat

	postgresDB, err := postgres.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	// Redis bağlantısını başlat
	// if redisCache, err := cache.NewRedisCache("", ""); err != nil {
	//  log.Fatalf("Failed to connect to Redis: %v", err)
	// }

	defer postgresDB.Pool.Close()
	//defer redisCache.Close()

	DB := repository.NewDB(postgresDB.Pool, nil)

	// RabbitMQ bağlantısını başlat
	// rabbitMQ, err := messaging.NewRabbitMQ(config.AppConfig.RabbitMQ.URL)
	// if err != nil {
	//  log.Println("Failed to connect to RabbitMQ: ", err)
	// }
	// defer rabbitMQ.Close()

	adapters.SetupRoutes(app, DB, nil)

	log.Fatal(app.Listen(":45013"))
}
