package routes

import (
	"chickenfile/controllers"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.POST("/upload", controllers.Upload)
	e.POST("/download", controllers.Download)
}
