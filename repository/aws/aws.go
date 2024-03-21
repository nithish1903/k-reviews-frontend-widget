package aws

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

func CreateAwsSession() *session.Session {
	region := os.Getenv("REGION")
	aws_id := os.Getenv("AWS_ID")
	aws_secret := os.Getenv("AWS_SECRET")

	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(aws_id, aws_secret, ""),
		})

	if err != nil {
		zap.L().Fatal("Failed to initialize new aws session", zap.Any("err", err))
	}
	zap.L().Info("Creating AWS session")
	return sess
}

func CreateFolderInsideS3Bucket(bucketName, folderName string) *s3.PutObjectInput {
	return &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(folderName + "/"), // Ensure the folder name ends with "/"
	}
}
