package minio_pkg

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

func InitBucket(logger *zap.Logger, minioClient *minio.Client) {
	bucket_name := os.Getenv("MINIO_BUCKET_NAME")
	println(bucket_name)
	if bucket_name == "" {
		bucket_name = "media"
	}

	err := EnsureBucket(minioClient, bucket_name)
	if err != nil {
		logger.Fatal("❌ failed to ensure bucket:", zap.Error(err))
	}

	policy := `{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {"AWS": ["*"]},
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::` + bucket_name + `/*"]
            }
        ]
    }`

	err = minioClient.SetBucketPolicy(context.Background(), bucket_name, policy)
	if err != nil {
		logger.Fatal("❌ failed to set bucket policy:", zap.Error(err))
	}
}

func InitS3(logger *zap.Logger) *minio.Client {
	user := os.Getenv("MINIO_ACCESS_KEY")
	password := os.Getenv("MINIO_SECRET_KEY")
	host := os.Getenv("MINIO_ENDPOINT")
	useSSL := os.Getenv("MINIO_USE_SSL")
	region := os.Getenv("MINIO_REGION")
	minioClient, err := minio.New(host, &minio.Options{
		Creds:        credentials.NewStaticV4(user, password, ""),
		Secure:       (useSSL == "true"), // set true if using HTTPS
		BucketLookup: minio.BucketLookupAuto,
		Region:       region,
	})

	if err != nil {
		logger.Fatal("Fail connecting to Minio", zap.Error(err))
	}

	InitBucket(logger, minioClient)

	return minioClient
}

func EnsureBucket(minioClient *minio.Client, bucketName string) error {
	ctx := context.Background()
	region := os.Getenv("MINIO_REGION")

	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
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
