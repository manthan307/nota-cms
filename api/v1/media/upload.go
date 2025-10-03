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
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

// Configure MinIO/S3 client
func newMinioClient() (*minio.Client, error) {
	return minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false, // set true if using HTTPS
	})
}

func UploadMediaHandler(queries *db.Queries, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//get user id
		userID := c.Locals("claims").(jwt.MapClaims)["user_id"].(string)
		UUIDuserID, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
		}
		pgUserId := pgtype.UUID{
			Bytes: UUIDuserID,
			Valid: true,
		}

		// Parse file from form-data
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file is required"})
		}

		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			logger.Error("Error opening uploaded file", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot open file"})
		}
		defer file.Close()

		// Init MinIO client
		minioClient, err := newMinioClient()
		if err != nil {
			logger.Error("Error creating MinIO client", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "storage not available"})
		}

		bucket := os.Getenv("MINIO_BUCKET_NAME")
		if bucket == "" {
			bucket = "media"
		}
		ctx := context.Background()

		// Ensure bucket exists
		exists, err := minioClient.BucketExists(ctx, bucket)
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed checking bucket"})
		}
		if !exists {
			err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create bucket"})
			}
		}

		// Generate unique object key
		objectKey := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))

		// Upload file to MinIO
		uploadInfo, err := minioClient.PutObject(
			ctx,
			bucket,
			objectKey,
			file,
			fileHeader.Size,
			minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")},
		)
		if err != nil {
			logger.Error("Error uploading file to MinIO", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "upload failed"})
		}

		// Construct URL (depends on how your S3/MinIO is exposed)
		fileURL := fmt.Sprintf("http://localhost:9000/%s/%s", bucket, objectKey)

		// Save metadata in DB
		media, err := queries.CreateMedia(ctx, db.CreateMediaParams{
			Key:        objectKey,
			Url:        fileURL,
			Bucket:     bucket,
			Type:       detectFileType(fileHeader), // image/video/file
			UploadedBy: pgUserId,
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
			"uploaded":  media.UploadedBy,
			"size":      uploadInfo.Size,
			"createdAt": time.Now(),
		})
	}
}

// simple MIME type detection
func detectFileType(fileHeader *multipart.FileHeader) string {
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		return "file"
	}
	if strings.HasPrefix(contentType, "image/") {
		return "image"
	}
	if strings.HasPrefix(contentType, "video/") {
		return "video"
	}
	return "file"
}
