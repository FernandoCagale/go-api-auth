package handlers

import (
	"net/http"

	"github.com/FernandoCagale/go-api-auth/src/checker"
	"github.com/labstack/echo"
)

type HealthzHandler struct {
	checkers map[string]checker.Checker
}

func NewHealthzHandler(checkers map[string]checker.Checker) *HealthzHandler {
	return &HealthzHandler{checkers}
}

func (h *HealthzHandler) HealthzIndex(c echo.Context) error {
	payload := make(map[string]bool)

	for k, v := range h.checkers {
		payload[k] = v.IsAlive()
	}

	return c.JSON(http.StatusOK, payload)
}
