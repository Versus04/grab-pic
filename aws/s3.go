package aws

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func UploadToS3(file *os.File) (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error Loading the .env File ", err)

	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("error loading the configuration ", err)
	}
	client := s3.NewFromConfig(cfg)
	key := "uploads/" + filepath.Base(file.Name())
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("photo-bucket-67"),
		Key:    aws.String(key),

		Body: file,
	})
	if err != nil {
		return "", err
	}
	presignClient := s3.NewPresignClient(client)
	req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("photo-bucket-67"),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err

	}
	return req.URL, nil
}
