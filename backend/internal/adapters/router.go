package adapters

import (
	"sasa-elterminali-service/internal/adapters/handler"
	"sasa-elterminali-service/internal/adapters/repository"
	"sasa-elterminali-service/internal/core/services"
	"sasa-elterminali-service/internal/messaging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App, db *repository.DB, rabbitMQ *messaging.RabbitMQ) {

	// CORS
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "http://127.0.0.1:45013",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Swagger endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Create a new Messaging instance
	webSocketMessaging := messaging.NewMessaging()
	// WebSocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	// WebSocket endpoint
	app.Get("/ws", websocket.New(webSocketMessaging.GetEmployeeID))
	// Send message to all connected devices
	app.Post("/api/v1/send", webSocketMessaging.SendMessageToAll)
	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(repository.SecretKey),
	}))

	// Ariza
	arizaService := services.NewArizaService(db)
	arizaHandler := handler.NewArizaHandler(*arizaService, *webSocketMessaging)
	app.Post("/api/v1/ariza/get", arizaHandler.GetAriza)

}
