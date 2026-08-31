package handler

import (
	"fmt"
	"net/http"

	"github.com/arthur404dev/gopportunities/schemas"
	"github.com/labstack/echo/v4"
)

func sendError(ctx echo.Context, code int, msg string) error {
	ctx.Response().Header().Set("Content-type", "application/json")
	return ctx.JSON(code, map[string]interface{}{
		"message":   msg,
		"errorCode": code,
	})
}

func sendSuccess(ctx echo.Context, op string, data interface{}) error {
	ctx.Response().Header().Set("Content-type", "application/json")
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("operation from handler: %s successfull", op),
		"data":    data,
	})
}

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type CreateOpeningResponse struct {
	Message string                  `json:"message"`
	Data    schemas.OpeningResponse `json:"data"`
}

type DeleteOpeningResponse struct {
	Message string                  `json:"message"`
	Data    schemas.OpeningResponse `json:"data"`
}
type ShowOpeningResponse struct {
	Message string                  `json:"message"`
	Data    schemas.OpeningResponse `json:"data"`
}
type ListOpeningsResponse struct {
	Message string                    `json:"message"`
	Data    []schemas.OpeningResponse `json:"data"`
}
type UpdateOpeningResponse struct {
	Message string                  `json:"message"`
	Data    schemas.OpeningResponse `json:"data"`
}
