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

func CreateContentHandler(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//get body
		var body struct {
			SchemaID  string                 `json:"schema_id"`
			Data      map[string]interface{} `json:"data"`
			Published bool                   `json:"published"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid body",
			})
		}

		//get schema from db using schemaID
		uuidId, err := uuid.Parse(body.SchemaID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid schema ID",
			})
		}
		schema, err := queries.GetSchemaByID(c.Context(), uuidId)
		if err != nil {
			logger.Error("Error fetching schema", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching schema",
			})
		}

		ok, err := utils.CompareSchemaWithData(schema.Definition, body.Data)
		if !ok {
			logger.Error("Data does not match schema", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Data does not match schema: " + err.Error(),
			})
		}

		// Marshal data into JSON for insertion
		dataBytes, err := json.Marshal(body.Data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error encoding data",
			})
		}

		pguuid := pgtype.UUID{Bytes: uuidId, Valid: true}

		content, err := queries.CreateContent(c.Context(), db.CreateContentParams{
			SchemaID:  pguuid,
			Data:      dataBytes,
			Published: pgtype.Bool{Bool: body.Published, Valid: true},
		})
		if err != nil {
			logger.Error("Error creating content", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating content",
			})
		}

		return c.JSON(fiber.Map{
			"id":        content.ID,
			"schemaID":  content.SchemaID,
			"data":      content.Data,
			"published": content.Published,
			"createdAt": content.CreatedAt,
			"updateAt":  content.UpdatedAt,
		})
	}
}
