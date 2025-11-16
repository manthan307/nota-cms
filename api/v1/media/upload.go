package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/manthan307/nota-cms/db/output"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func UploadMediaHandler(queries *db.Queries, logger *zap.Logger, minioClient *minio.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(jwt.MapClaims)
		userID := claims["user_id"].(string)

		uid, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
		}

		pgUserID := pgtype.UUID{Bytes: uid, Valid: true}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file is required"})
		}

		file, err := fileHeader.Open()
		if err != nil {
			logger.Error("Error opening uploaded file", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot open file"})
		}
		defer file.Close()

		bucket := os.Getenv("MINIO_BUCKET_NAME")
		if bucket == "" {
			bucket = "media"
		}

		ctx := context.Background()

		exists, err := minioClient.BucketExists(ctx, bucket)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "bucket check failed"})
		}

		if !exists {
			err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create bucket"})
			}
		}

		objectKey := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))

		uploadInfo, err := minioClient.PutObject(
			ctx,
			bucket,
			objectKey,
			file,
			fileHeader.Size,
			minio.PutObjectOptions{
				ContentType: fileHeader.Header.Get("Content-Type"),
			},
		)
		if err != nil {
			logger.Error("Error uploading to MinIO", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "upload failed"})
		}

		fileURL := fmt.Sprintf("http://%s/%s/%s",
			os.Getenv("MINIO_ENDPOINT"),
			bucket,
			objectKey,
		)

		media, err := queries.CreateMedia(ctx, db.CreateMediaParams{
			Key:        objectKey,
			Url:        fileURL,
			Bucket:     bucket,
			Type:       detectFileType(fileHeader),
			UploadedBy: pgUserID,
		})
		if err != nil {
			logger.Error("Error saving media to DB", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db insert failed"})
		}

		return c.JSON(fiber.Map{
			"id":        media.ID,
			"key":       media.Key,
			"url":       media.Url,
			"type":      media.Type,
			"bucket":    media.Bucket,
			"uploaded":  media.UploadedBy.Bytes, // FIXED
			"size":      uploadInfo.Size,
			"createdAt": time.Now(),
		})
	}
}

func detectFileType(fileHeader *multipart.FileHeader) string {
	t := fileHeader.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(t, "image/"):
		return "image"
	case strings.HasPrefix(t, "video/"):
		return "video"
	default:
		return "file"
	}
}
