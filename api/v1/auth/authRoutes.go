package auth

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/manthan307/nota-cms/db/output"
	"github.com/manthan307/nota-cms/utils"
	"go.uber.org/zap"
)

func Register(app *fiber.App, queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		hash, _ := utils.HashPassword(body.Password)
		user, err := queries.CreateUser(c.Context(), db.CreateUserParams{
			Email:        body.Email,
			PasswordHash: hash,
			Role:         body.Role,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"role":      user.Role,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		})
	}
}

func Login(app *fiber.App, queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		user, err := queries.GetUserByEmail(c.Context(), body.Email)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		if !utils.CheckPasswordHash(body.Password, user.PasswordHash) {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		token, _ := utils.GenerateJWT(int(user.ID), user.Role)
		return c.JSON(fiber.Map{"token": token})
	}
}
