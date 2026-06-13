package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/service"
	"github.com/velocity/server-monitoring/backend/pkg/response"
)

type ServerHandler struct {
	svc service.ServerService
}

func NewServerHandler(svc service.ServerService) *ServerHandler {
	return &ServerHandler{svc: svc}
}

func (h *ServerHandler) Create(c *fiber.Ctx) error {
	var server domain.Server
	if err := c.BodyParser(&server); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request"))
	}

	if server.Name == "" || server.Hostname == "" || server.IPAddress == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("name, hostname, and ip_address are required"))
	}

	if err := h.svc.CreateServer(&server); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("could not create server"))
	}

	return c.JSON(response.Success(server))
}

func (h *ServerHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result, err := h.svc.GetServers(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("could not retrieve servers"))
	}

	return c.JSON(response.Success(result))
}

func (h *ServerHandler) Detail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid id"))
	}

	server, err := h.svc.GetServerByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error("server not found"))
	}

	return c.JSON(response.Success(server))
}

func (h *ServerHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid id"))
	}

	server, err := h.svc.GetServerByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error("server not found"))
	}

	if err := c.BodyParser(&server); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request"))
	}

	if err := h.svc.UpdateServer(server); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("could not update server"))
	}

	return c.JSON(response.Success(server))
}

func (h *ServerHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid id"))
	}

	if err := h.svc.DeleteServer(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error("could not delete server"))
	}

	return c.JSON(response.Success(nil))
}

func (h *ServerHandler) RegisterRoutes(router fiber.Router) {
	// All routes are protected by AuthMiddleware (set up in main)
	router.Post("/", h.Create)
	router.Get("/", h.List)
	router.Get("/:id", h.Detail)
	router.Put("/:id", h.Update)
	router.Delete("/:id", h.Delete)
}
