package schemasRoutes

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/manthan307/nota-cms/db/output"
	"go.uber.org/zap"
)

func ListSchemas(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Use the request context (helps with cancellation, tracing, etc.)
		ctx := c.Context()

		schemas, err := queries.ListSchemas(ctx)
		if err != nil {
			logger.Error("Failed to fetch schemas", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch schemas",
			})
		}

		if len(schemas) == 0 {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "No schemas found",
				"data":    []interface{}{},
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"count": len(schemas),
			"data":  schemas,
		})
	}
}
