package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/service"
	"github.com/velocity/server-monitoring/backend/pkg/response"
)

type MetricHandler struct {
	metricSvc service.MetricService
	serverSvc service.ServerService
}

func NewMetricHandler(metricSvc service.MetricService, serverSvc service.ServerService) *MetricHandler {
	return &MetricHandler{
		metricSvc: metricSvc,
		serverSvc: serverSvc,
	}
}

// AgentMiddleware verifies the API Key provided by the agent
func (h *MetricHandler) AgentMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error("missing authorization header"))
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error("invalid authorization format"))
		}

		apiKey := parts[1]
		server, err := h.serverSvc.GetServerByAPIKey(apiKey)
		if err != nil || server == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error("invalid api key"))
		}

		// Store server ID in context
		c.Locals("server_id", server.ID)
		return c.Next()
	}
}

func (h *MetricHandler) Ingest(c *fiber.Ctx) error {
	serverID := c.Locals("server_id").(uint)

	var metric domain.Metric
	if err := c.BodyParser(&metric); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request body"))
	}

	metric.ServerID = serverID
	if err := h.metricSvc.SaveMetric(&metric); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("could not save metric"))
	}

	return c.JSON(response.Success(nil))
}

func (h *MetricHandler) RegisterRoutes(router fiber.Router) {
	// API specifically for Agents to report metrics
	agentGroup := router.Group("/metrics", h.AgentMiddleware())
	agentGroup.Post("/", h.Ingest)
}
