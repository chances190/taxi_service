package middlewares

import (
	"log"
	"taxi_service/internal/apperrors"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler middleware captura erros retornados pelos handlers e aplica o formato padronizado.
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			status := apperrors.HTTPStatus(err)
			payload := apperrors.ToPayload(err)
			if status >= 500 { // logar erros de servidor
				log.Printf("internal error: %v", err)
			}
			return c.Status(status).JSON(payload)
		}
		return nil
	}
}
