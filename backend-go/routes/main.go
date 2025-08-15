package routes

import (
	"taxi_service/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes inicializa todas as rotas da aplicação.
func SetupRoutes(app *fiber.App) {
	// Middlewares
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(middlewares.ErrorHandler())

	// Grupo de rotas da API
	api := app.Group("/")

	// Rota de Health Check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
	})

	SetupMotoristaRoutes(api)
}
