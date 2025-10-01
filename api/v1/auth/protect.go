package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	db "github.com/manthan307/nota-cms/db/output"
	"github.com/manthan307/nota-cms/utils"
	"go.uber.org/zap"
)

var roleHierarchy = map[string]int{
	"viewer": 1,
	"editor": 2,
	"admin":  3,
}

func ProtectedRoute(logger *zap.Logger, queries *db.Queries, privilage string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("token")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		// Verify token
		token, err := utils.VerifyJWT(tokenStr)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid claims"})
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user_id"})
		}

		// Convert to UUID
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user_id format"})
		}

		// Check existence in DB
		exists, err := queries.UserExists(c.Context(), userID)
		if err != nil {
			logger.Error("failed to check user existence", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user does not exist"})
		}

		//get role
		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid role"})
		}

		requiredLevel := roleHierarchy[privilage]
		userLevel := roleHierarchy[role]

		if userLevel < requiredLevel {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		// Attach claims for downstream handlers
		c.Locals("claims", claims)

		return c.Next()
	}
}
