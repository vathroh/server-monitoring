package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/velocity/server-monitoring/backend/internal/service"
	"github.com/velocity/server-monitoring/backend/pkg/response"
)

type DashboardHandler struct {
	serverSvc service.ServerService
	metricSvc service.MetricService
}

func NewDashboardHandler(serverSvc service.ServerService, metricSvc service.MetricService) *DashboardHandler {
	return &DashboardHandler{
		serverSvc: serverSvc,
		metricSvc: metricSvc,
	}
}

func (h *DashboardHandler) Summary(c *fiber.Ctx) error {
	summary, err := h.serverSvc.GetDashboardSummary()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("could not fetch dashboard summary"))
	}
	return c.JSON(response.Success(summary))
}

func (h *DashboardHandler) Trend(c *fiber.Ctx) error {
	serverIDStr := c.Params("id")
	serverID, err := strconv.Atoi(serverIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid server id"))
	}

	timeRange := c.Query("range", "1h") // 1h, 24h, 7d

	trend, err := h.metricSvc.GetMetricsTrend(uint(serverID), timeRange)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("could not fetch metrics trend"))
	}

	return c.JSON(response.Success(trend))
}

func (h *DashboardHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/summary", h.Summary)
	router.Get("/servers/:id/trend", h.Trend)
}
