package auth

import (
	"os"

	"github.com/gofiber/fiber/v2"
	db "github.com/manthan307/nota-cms/db/output"
	"github.com/manthan307/nota-cms/utils"
	"go.uber.org/zap"
)

func RegisterHandler(queries *db.Queries, logger *zap.Logger) fiber.Handler {

	env := os.Getenv("ENV")

	return func(c *fiber.Ctx) error {

		exist, err := queries.AdminExists(c.Context())
		if err != nil {
			logger.Error("failed to check admin existence", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
		}

		if exist {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin already exists"})
		}

		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		hash, err := utils.HashPassword(body.Password)
		if err != nil {
			logger.Error("failed to hash password", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
		}

		user, err := queries.CreateUser(c.Context(), db.CreateUserParams{
			Email:        body.Email,
			PasswordHash: hash,
			Role:         body.Role,
		})
		if err != nil {
			logger.Error("failed to create user", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
		}

		token, err := utils.GenerateJWT(user.ID, user.Role)
		if err != nil {
			logger.Error("failed to generate jwt", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create token"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    token,
			HTTPOnly: true,
			Secure:   env == "PRODUCTION",
			SameSite: "Lax",
			Path:     "/",
			MaxAge:   3600 * 24,
		})

		return c.JSON(fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"role":      user.Role,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		})
	}
}

func LoginHandler(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	env := os.Getenv("ENV")

	return func(c *fiber.Ctx) error {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		user, err := queries.GetUserByEmail(c.Context(), body.Email)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}

		if !utils.CheckPasswordHash(body.Password, user.PasswordHash) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}

		token, err := utils.GenerateJWT(user.ID, user.Role)
		if err != nil {
			logger.Error("failed to generate jwt", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create token"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    token,
			HTTPOnly: true,
			Secure:   env == "PRODUCTION",
			SameSite: "Lax",
			Path:     "/",
			MaxAge:   3600 * 24,
		})

		return c.JSON(fiber.Map{"token": token})
	}
}
