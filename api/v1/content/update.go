package content

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/manthan307/nota-cms/db/output"
	"github.com/manthan307/nota-cms/utils"
	"go.uber.org/zap"
)

func UpdateContentHandler(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			ContentID string                 `json:"content_id"`
			Data      map[string]interface{} `json:"data"`
			Published bool                   `json:"published"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid body",
			})
		}

		// Validate content ID
		contentID, err := uuid.Parse(body.ContentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid content ID",
			})
		}

		// Fetch content
		content, err := queries.GetContentByID(c.Context(), contentID)
		if err != nil {
			logger.Error("Error fetching content by ID", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not fetch content",
			})
		}

		// Parse schema UUID from pgtype.UUID
		schemaID, err := uuid.FromBytes(content.SchemaID.Bytes[:])
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid schema ID",
			})
		}

		// Fetch related schema
		schema, err := queries.GetSchemaByID(c.Context(), schemaID)
		if err != nil {
			logger.Error("Error fetching schema", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not fetch schema",
			})
		}

		// Validate data with schema
		ok, err := utils.CompareSchemaWithData(schema.Definition, body.Data)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Data does not match schema: " + err.Error(),
			})
		}

		// Marshal JSON
		dataBytes, err := json.Marshal(body.Data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not encode JSON",
			})
		}

		// UPDATE Content
		updated, err := queries.UpdateContent(c.Context(), db.UpdateContentParams{
			ID:   content.ID,
			Data: dataBytes,
			Published: pgtype.Bool{
				Bool:  body.Published,
				Valid: true,
			},
		})

		if err != nil {
			logger.Error("Error updating content", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not update content",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"id":        updated.ID,
			"schemaID":  updated.SchemaID,
			"data":      updated.Data,
			"published": updated.Published,
			"createdAt": updated.CreatedAt,
			"updatedAt": updated.UpdatedAt,
		})
	}
}
