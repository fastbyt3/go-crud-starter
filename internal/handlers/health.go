package handlers

import (
	"net/http"

	"github.com/fastbyt3/diy-mssql-test/pkg/database"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	databaseService *database.Service
}

func NewHealthHandler(databaseService database.Service) *HealthHandler {
	return &HealthHandler{databaseService: &databaseService}
}

func (h *HealthHandler) ApiHealth(c echo.Context) error {
	resp := map[string]string{
		"message": "All Ok!",
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HealthHandler) DbHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, h.databaseService.Health())
}
