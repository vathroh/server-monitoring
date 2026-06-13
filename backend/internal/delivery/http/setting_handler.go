package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/velocity/server-monitoring/backend/internal/service"
	"github.com/velocity/server-monitoring/backend/pkg/response"
)

type SettingHandler struct {
	settingSvc service.SettingService
}

func NewSettingHandler(settingSvc service.SettingService) *SettingHandler {
	return &SettingHandler{settingSvc: settingSvc}
}

func (h *SettingHandler) GetSettings(c *fiber.Ctx) error {
	settings, err := h.settingSvc.GetAllSettings()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Failed to fetch settings"))
	}
	return c.JSON(response.Success(settings))
}

func (h *SettingHandler) SaveSettings(c *fiber.Ctx) error {
	var req map[string]string
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid request payload"))
	}

	if err := h.settingSvc.UpsertSettings(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("Failed to save settings"))
	}

	return c.JSON(response.Success(nil))
}

func (h *SettingHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/", h.GetSettings)
	router.Post("/", h.SaveSettings)
}
