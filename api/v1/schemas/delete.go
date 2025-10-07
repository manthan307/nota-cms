package schemasRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	db "github.com/manthan307/nota-cms/db/output"
	"go.uber.org/zap"
)

func DeleteSchema(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuidID, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid UUID",
			})
		}
		// Use the request context (helps with cancellation, tracing, etc.)
		ctx := c.Context()

		err = queries.DeleteSchema(ctx, uuidID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Failed to find schemas",
			})
		}

		return c.Status(fiber.StatusOK).Send([]byte("Ok"))
	}
}
