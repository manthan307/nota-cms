package media

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	db "github.com/manthan307/nota-cms/db/output"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func DeleteMediaHandler(queries *db.Queries, logger *zap.Logger, minioClient *minio.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//get user id
		userID := c.Locals("claims").(jwt.MapClaims)["user_id"].(string)
		_, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
		}

		// Delete file path from request body
		type DeleteMediaRequest struct {
			FilePath string `json:"file_path"`
		}
		var req DeleteMediaRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		bucket := os.Getenv("MINIO_BUCKET_NAME")
		if bucket == "" {
			bucket = "media"
		}
		ctx := context.Background()
		// Delete the file from MinIO
		err = minioClient.RemoveObject(ctx, bucket, req.FilePath, minio.RemoveObjectOptions{})
		if err != nil {
			logger.Error("Error deleting file from MinIO", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete file"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "file deleted successfully"})
	}
}
