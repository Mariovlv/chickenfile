package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)
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

func main() {
	e := echo.New()

	e.Use(middleware.BodyLimit("10M"))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":1323"))
}
