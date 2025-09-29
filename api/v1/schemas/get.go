package schemasRoutes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	db "github.com/manthan307/nota-cms/db/output"
	"go.uber.org/zap"
)

func GetSchemaByID(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		parsed, err := uuid.Parse(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		schema, err := queries.GetSchemaByID(c.Context(), parsed)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "schema not found"})
			}
			logger.Error("failed to fetch schema", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not fetch schema"})
		}

		return c.JSON(fiber.Map{
			"id":         schema.ID,
			"name":       schema.Name,
			"createdBy":  schema.CreatedBy,
			"definition": schema.Definition,
			"createdAt":  schema.CreatedAt,
			"updatedAt":  schema.UpdatedAt,
		})
	}
}

func GetSchemaByName(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.Params("name")
		schema, err := queries.GetSchemaByName(c.Context(), name)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "schema not found"})
			}
			logger.Error("failed to fetch schema", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not fetch schema"})
		}

		return c.JSON(fiber.Map{
			"id":         schema.ID,
			"name":       schema.Name,
			"createdBy":  schema.CreatedBy,
			"definition": schema.Definition,
			"createdAt":  schema.CreatedAt,
			"updatedAt":  schema.UpdatedAt,
		})
	}
}
