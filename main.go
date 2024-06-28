package main

import (
	"bytes"
	"chickenfile/client"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client.InitializeS3Client()
}

func upload(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "File is required")
	}

	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxFileSize {
		return echo.NewHTTPError(http.StatusBadRequest, "File too big")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	/*
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		client := s3.NewFromConfig(cfg)
	*/

	client := client.GetS3Client()
	uploader := manager.NewUploader(client)

	bucketName := os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "AWS_BUCKET_NAME not set in environment")
	}

	key := file.Filename

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   src,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upload file to S3")
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully to %s with fields name=%s and email=%s.</p><p>Uploaded to: %s</p>", file.Filename, bucketName, name, email, result.Location))
}

func download(c echo.Context) error {
	// Access the form value
	bodyWord := c.FormValue("word")

	client := client.GetS3Client()

	bucketName := os.Getenv("AWS_BUCKET_NAME")
	// Example S3 operation: List objects in a bucket
	output, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bodyWord),
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list object in S3 bucket")
	}
	defer output.Body.Close()

	// Read the content of the file
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(output.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read file content")
	}

	// You can add additional logic to decode the hash, request the resource from the S3 bucket, etc.

	// Return the form value as a response for demonstration purposes
	return c.Blob(http.StatusOK, *output.ContentType, buf.Bytes())
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")

	e.POST("/upload", upload)
	e.POST("/download", download)

	e.Logger.Fatal(e.Start(":1323"))
}
