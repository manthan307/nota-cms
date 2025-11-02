package auth

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
			Role:         "admin",
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
			Secure:   env == "PRODUCTION", // use HTTPS in prod only
			SameSite: "Lax",               // allow cross-origin
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
		if err != nil || !utils.CheckPasswordHash(body.Password, user.PasswordHash) {
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

		return c.JSON(fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		})
	}
}

func CheckAuthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("token")

		if tokenStr == "" {
			return c.Status(401).JSON(fiber.Map{"auth": false})
		}

		token, err := utils.VerifyJWT(tokenStr)
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"auth": false})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"auth": false})
		}

		role, _ := claims["role"].(string)

		return c.Status(200).JSON(fiber.Map{
			"auth": true,
			"role": role,
		})
	}
}
