package schemasRoutes

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/manthan307/nota-cms/db/output"
	"go.uber.org/zap"
)

func SchemasCreateHandler(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			Name       string          `json:"name"`
			Defination json.RawMessage `json:"definition"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		if body.Name == "" || len(body.Defination) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name and definition are required"})
		}

		claims := c.Locals("claims").(jwt.MapClaims)
		userIDStr := claims["user_id"].(string)

		// Parse string UUID
		parsedUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
		}

		userID := pgtype.UUID{Bytes: parsedUUID, Valid: true}
		schema, err := queries.CreateSchema(c.Context(), db.CreateSchemaParams{
			CreatedBy:  userID,
			Name:       body.Name,
			Definition: body.Defination,
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create schema"})
		}

		return c.JSON(fiber.Map{
			"id":         schema.ID,
			"name":       schema.Name,
			"createdBy":  schema.CreatedBy,
			"definition": schema.Definition,
		})
	}
}
