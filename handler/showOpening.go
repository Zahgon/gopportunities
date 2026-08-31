package handler

import (
	"net/http"

	"github.com/arthur404dev/gopportunities/schemas"
	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Show opening
// @Description Show a job opening
// @Tags Openings
// @Accept json
// @Produce json
// @Param id query string true "Opening identification"
// @Success 200 {object} ShowOpeningResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /opening [get]
func ShowOpeningHandler(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	if id == "" {
		return sendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "queryParameter").Error())
	}
	opening := schemas.Opening{}
	if err := db.First(&opening, id).Error; err != nil {
		return sendError(ctx, http.StatusNotFound, "opening not found")
	}

	return sendSuccess(ctx, "show-opening", opening)
}
