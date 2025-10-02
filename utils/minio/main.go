package minio

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitS3() (*minio.Client, error) {
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin123", ""),
		Secure: false, // set true if using HTTPS
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func EnsureBucket(minioClient *minio.Client, bucket string) error {
	ctx := context.Background()

	// Check if the bucket exists
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	if !exists {
		// Create the bucket
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		return nil
	}

	return nil
}

func UploadFile(minioClient *minio.Client, bucket, key, filepath string) (minio.UploadInfo, error) {
	// Ensure bucket exists
	EnsureBucket(minioClient, bucket)

	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("❌ failed to open file: %v", err)
	}
	defer file.Close()

	// Upload to MinIO
	uploadInfo, err := minioClient.PutObject(context.Background(), bucket, key, file, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("failed to upload file: %v", err)
	}

	return uploadInfo, nil
}
