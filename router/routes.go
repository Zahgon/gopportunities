package router

import (
	docs "github.com/arthur404dev/gopportunities/docs"
	"github.com/arthur404dev/gopportunities/handler"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func initializeRoutes(e *echo.Echo) {
	// Initialize Handler
	handler.InitializeHandler()
	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath
	v1 := e.Group(basePath)
	{
		v1.GET("/opening", handler.ShowOpeningHandler)
		v1.POST("/opening", handler.CreateOpeningHandler)
		v1.DELETE("/opening", handler.DeleteOpeningHandler)
		v1.PUT("/opening", handler.UpdateOpeningHandler)
		v1.GET("/openings", handler.ListOpeningsHandler)

	}
	// Initialize Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
