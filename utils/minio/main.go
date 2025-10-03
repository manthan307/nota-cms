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
	user := os.Getenv("MINIO_ROOT_USER")
	password := os.Getenv("MINIO_ROOT_PASSWORD")
	// Initialize minio client object.
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:        credentials.NewStaticV4(user, password, ""),
		Secure:       false, // set true if using HTTPS
		BucketLookup: minio.BucketLookupAuto,
		Region:       "us-east-1",
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func EnsureBucket(minioClient *minio.Client, bucketName string) error {
	ctx := context.Background()

	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
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
