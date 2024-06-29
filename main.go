package main

import (
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
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	fmt.Println("Checking ENV variable:", os.Getenv("ENV"))

	if os.Getenv("ENV") == "development" {
		fmt.Println("Loading .env file...")

		fmt.Println(".env file loaded successfully.")
	} else {
		fmt.Println("Not in development mode.")
	}

	fmt.Println("Initialization completed.")
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
	log.Printf("AWS_BUCKET_NAME: %s", bucketName)
	if bucketName == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "AWS_BUCKET_NAME not set in environment")
	}

	secretWord := os.Getenv("SECRET_WORD")
	log.Printf("SECRET_WORD: %s", secretWord)
	if secretWord == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "SECRET_WORD not set in environment")
	}

	keyForAWS := generateHMAC(keyword, secretWord)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file.Filename),
		Metadata: map[string]string{
			"hash": keyForAWS,
		},
		Body: src,
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

	resp, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list objects in S3 bucket")
	}

	var matchedObjectKey string
	for _, item := range resp.Contents {
		headResp, err := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    item.Key,
		})
		if err != nil {
			continue
		}
		if headResp.Metadata["hash"] == key {
			matchedObjectKey = aws.ToString(item.Key)
			break
		}
	}

	if matchedObjectKey == "" {
		return echo.NewHTTPError(http.StatusNotFound, "File not found")
	}

	getObjectOutput, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(matchedObjectKey),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get object from S3 bucket")
	}
	defer getObjectOutput.Body.Close()

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", matchedObjectKey))
	c.Response().Header().Set(echo.HeaderContentType, aws.ToString(getObjectOutput.ContentType))

	return c.Stream(http.StatusOK, aws.ToString(getObjectOutput.ContentType), getObjectOutput.Body)
}

func main() {
	e := echo.New()

	e.Static("/", "public")

	e.POST("/upload", upload)
	e.POST("/download", download)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
