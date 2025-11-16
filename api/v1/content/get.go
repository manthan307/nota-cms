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

func GetAllContentsBySchemaHandler(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schemaName := c.Params("schema_name")

		// Fetch schema
		schema, err := queries.GetSchemaByName(c.Context(), schemaName)
		if err != nil {
			logger.Error("Error fetching schema", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching schema",
			})
		}

		// Check published query param: true | false | all
		p := c.Query("published", "all")

		var (
			contents []db.Content
		)

		pgID := pgtype.UUID{Bytes: schema.ID, Valid: true}

		switch p {
		case "true":
			b := true
			contents, err = queries.GetContentsBySchema(
				c.Context(),
				db.GetContentsBySchemaParams{
					SchemaID: pgID,
					Published: pgtype.Bool{
						Bool:  b,
						Valid: true,
					},
				},
			)
		case "false":
			b := false
			contents, err = queries.GetContentsBySchema(
				c.Context(),
				db.GetContentsBySchemaParams{
					SchemaID: pgID,
					Published: pgtype.Bool{
						Bool:  b,
						Valid: true,
					},
				},
			)
		default: // "all"
			contents, err = queries.GetAllContentsBySchema(c.Context(), pgID)
		}

		if err != nil {
			logger.Error("Error fetching contents", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching contents",
			})
		}

		// Formatting output
		var result []map[string]interface{}
		for _, content := range contents {
			var data map[string]interface{}
			if err := json.Unmarshal(content.Data, &data); err != nil {
				logger.Warn("Invalid JSON in content.Data", zap.Error(err))
				continue
			}

			item := map[string]interface{}{
				"id":        content.ID,
				"schemaID":  content.SchemaID,
				"data":      data,
				"published": content.Published.Bool,
				"createdAt": content.CreatedAt,
				"updatedAt": content.UpdatedAt,
			}

			result = append(result, item)
		}

		return c.JSON(result)
	}
}
