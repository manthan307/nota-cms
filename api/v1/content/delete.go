package content

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	db "github.com/manthan307/nota-cms/db/output"
	"go.uber.org/zap"
)

func DeleteContentHandler(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentID := c.Params("id")
		uuidId, err := uuid.Parse(contentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid content ID",
			})
		}
		err = queries.DeleteContent(c.Context(), uuidId)
		if err != nil {
			logger.Error("Error deleting content", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error deleting content",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Content deleted successfully",
		})
	}
}
