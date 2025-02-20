package bootstrap

import (
	"context"
	"log"
	_ "net/http"
	_ "net/url"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ProvideMinIOClient(ctx context.Context, cfg *config.Config) *minio.Client {
	client, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, ""),
		Secure: cfg.S3UseSSL,
	})
	if err != nil {
		panic(err)
	}

	// test connection
	_, err = client.BucketExists(ctx, cfg.S3BucketName)
	if err != nil {
		log.Fatal("🚫 Cannot connect to MinIO | ", err)
	} else {
		log.Println("🫙 Connected to MinIO")
	}
	// buckets, err := client.ListBuckets(ctx)
	// // test connection
	// if err != nil {
	// 	log.Fatal("🚫 Cannot connect to MinIO | ", err)
	// } else {
	// 	log.Println("🫙 Connected to MinIO | ", buckets)
	// }
	return client
}
