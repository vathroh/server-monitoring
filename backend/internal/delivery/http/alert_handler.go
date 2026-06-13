package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/velocity/server-monitoring/backend/internal/service"
	"github.com/velocity/server-monitoring/backend/pkg/response"
)

type AlertHandler struct {
	alertSvc service.AlertService
}

func NewAlertHandler(alertSvc service.AlertService) *AlertHandler {
	return &AlertHandler{alertSvc: alertSvc}
}

func (h *AlertHandler) GetAlerts(c *fiber.Ctx) error {
	state := c.Query("state")
	serverIDStr := c.Query("server_id")
	var serverID uint
	if serverIDStr != "" {
		id, err := strconv.Atoi(serverIDStr)
		if err == nil {
			serverID = uint(id)
		}
	}

	alerts, err := h.alertSvc.GetAlerts(state, serverID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("could not fetch alerts"))
	}

	return c.JSON(response.Success(alerts))
}

func (h *AlertHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/", h.GetAlerts)
}
