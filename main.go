package main

import (
	"chickenfile/client"
	"chickenfile/routes"

	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func init() {
	/*
 		Change this is for production or development
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file, maybe in production enviroment:", err)
		}
	*/

	fmt.Println("Initialization completed.")
	client.InitializeS3Client()
}

func main() {
	e := echo.New()

	e.Static("/", "public")

	routes.InitRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
