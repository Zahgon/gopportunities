package handler

import (
	"net/http"

	"github.com/arthur404dev/gopportunities/schemas"
	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Create opening
// @Description Create a new job opening
// @Tags Openings
// @Accept json
// @Produce json
// @Param request body CreateOpeningRequest true "Request body"
// @Success 200 {object} CreateOpeningResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /opening [post]
func CreateOpeningHandler(ctx echo.Context) error {
	request := CreateOpeningRequest{}

	ctx.Bind(&request)

	if err := request.Validate(); err != nil {
		logger.Errorf("validation error: %v", err.Error())
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}

	opening := schemas.Opening{
		Role:     request.Role,
		Company:  request.Company,
		Location: request.Location,
		Remote:   *request.Remote,
		Link:     request.Link,
		Salary:   request.Salary,
	}

	if err := db.Create(&opening).Error; err != nil {
		logger.Errorf("error creating opening: %v", err.Error())
		return sendError(ctx, http.StatusInternalServerError, "error creating opening on database")
	}

	return sendSuccess(ctx, "create-opening", opening)
}
