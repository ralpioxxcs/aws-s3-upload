package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

var (
	s3Region          string
	s3AccessKey       string
	s3SecretAccessKey string
	s3Bucket          string
)

func loadAwsConfig() {
	s3Region = os.Getenv("S3_REGION")
	s3AccessKey = os.Getenv("S3_ACCESS_KEY")
	s3SecretAccessKey = os.Getenv("S3_SECRET_ACCESS_KEY")
	s3Bucket = os.Getenv("S3_BUCKET_NAME")

	log.Printf("*region : %v\n*access key : %v\n*secret access key : %v\n*bucket : %v\n", s3Region, s3AccessKey, s3SecretAccessKey, s3Bucket)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s [filename]\n", os.Args[0])
	}
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to open file (%v)\n", err)
	}
	defer file.Close()

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file (%v)\n", err)
	}
	loadAwsConfig()

	// connect S3
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Region),
		Credentials: credentials.NewStaticCredentials(
			s3AccessKey, s3SecretAccessKey, "",
		),
	})

	if err != nil {
		log.Fatalln(err)
	}

	// upload
	uploader := s3manager.NewUploader(session)
	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	log.Println(output)

	log.Printf("success to upload %v", filename)
}
