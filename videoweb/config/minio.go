package config

import (
	"context"
	"github.com/minio/minio-go/v7"
	"log"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinIoClient() {
	var err error
	MinioClient, err = minio.New(EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKeyID, SecretAccessKey, ""),
		Secure: SSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	err = MinioClient.MakeBucket(context.Background(), BucketName, minio.MakeBucketOptions{})
	if err != nil {
		log.Println(err)
	}
}
