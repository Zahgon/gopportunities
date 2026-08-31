package router

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRouter builds the fully wired HTTP handler with all routes registered.
// It returns the stdlib http.Handler interface so callers (including tests) do
// not depend on the concrete web framework type.
func SetupRouter() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())
	initializeRoutes(e)
	return e
}

func Initialize() {
	// Initialize Router
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	// Initialize Routes
	initializeRoutes(e)

	// Get the port from the environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run the server
	e.Logger.Fatal(e.Start("0.0.0.0:" + port))
}
