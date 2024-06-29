package main

import (
	"bytes"
	"chickenfile/client"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

func generateHMAC(input, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func upload(c echo.Context) error {
	keyword := c.FormValue("keyword")

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

	client := client.GetS3Client()
	uploader := manager.NewUploader(client)

	bucketName := os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "AWS_BUCKET_NAME not set in environment")
	}

	secretWord := os.Getenv("SECRET_WORD")
	if secretWord == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "SECRET_WORD not set in environment")
	}

	keyForAWS := generateHMAC(keyword, secretWord)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyForAWS),
		Body:   src,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upload file to S3")
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully to %s with secretword=%s.</p><p>Uploaded to: %s</p>", file.Filename, bucketName, keyword, result.Location))
}

func download(c echo.Context) error {
	bodyWord := c.FormValue("word")

	key := generateHMAC(bodyWord, os.Getenv("SECRET_WORD"))

	client := client.GetS3Client()

	bucketName := os.Getenv("AWS_BUCKET_NAME")
	output, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list object in S3 bucket")
	}
	defer output.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(output.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read file content")
	}

	return c.Blob(http.StatusOK, *output.ContentType, buf.Bytes())
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")

	e.POST("/upload", upload)
	e.POST("/download", download)

	e.Logger.Fatal(e.Start(":8080"))
}
