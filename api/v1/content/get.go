package content

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/manthan307/nota-cms/db/output"
	"go.uber.org/zap"
)

func GetContentHandler(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		parsedID, err := uuid.Parse(id)
		if err != nil {
			logger.Error("Error parsing UUID", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		content, err := queries.GetContentByID(c.Context(), parsedID)
		if err != nil {
			logger.Error("Error fetching content by ID", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching content",
			})
		}

		// Unmarshal JSON data
		var data map[string]interface{}
		if err := json.Unmarshal(content.Data, &data); err != nil {
			logger.Error("Error unmarshalling content data", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error parsing content data",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"id":        content.ID,
			"schemaID":  content.SchemaID,
			"data":      data,
			"createdAt": content.CreatedAt,
		})
	}
}

func GetAllContentsBySchemaHandler(queries *db.Queries, logger *zap.Logger, published bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		SchemaName := c.Params("schema_name")
		schema, err := queries.GetSchemaByName(c.Context(), SchemaName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching schema",
			})
		}
		pgID := pgtype.UUID{
			Bytes: schema.ID,
			Valid: true,
		}
		contents, err := queries.GetContentsBySchema(c.Context(), pgID)
		if err != nil {
			logger.Error("Error fetching all contents", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching contents",
			})
		}
		var result []map[string]interface{}
		for _, content := range contents {
			var data map[string]interface{}
			if err := json.Unmarshal(content.Data, &data); err != nil {
				logger.Error("Error unmarshalling content data", zap.Error(err))
				continue
			}
			item := map[string]interface{}{
				"id":        content.ID,
				"schemaID":  content.SchemaID,
				"data":      data,
				"createdAt": content.CreatedAt,
			}
			result = append(result, item)
		}
		return c.Status(fiber.StatusOK).JSON(result)
	}
}
