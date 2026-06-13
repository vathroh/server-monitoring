package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/velocity/server-monitoring/backend/internal/service"
	"github.com/velocity/server-monitoring/backend/pkg/response"
)

type AuthHandler struct {
	svc service.UserService
}

func NewAuthHandler(svc service.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request"))
	}

	accessToken, refreshToken, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error("invalid credentials"))
	}

	return c.JSON(response.Success(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}))
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request"))
	}

	aToken, rToken, err := h.svc.Refresh(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error("invalid refresh token"))
	}

	return c.JSON(response.Success(fiber.Map{
		"access_token":  aToken,
		"refresh_token": rToken,
	}))
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request"))
	}

	user, err := h.svc.Register(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	return c.JSON(response.Success(user))
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error("unauthorized"))
	}

	user, err := h.svc.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error("user not found"))
	}

	return c.JSON(response.Success(user))
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	type ForgotPasswordRequest struct {
		Email string `json:"email"`
	}

	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request"))
	}

	// Always return success to prevent email enumeration
	_ = h.svc.ForgotPassword(req.Email)
	return c.JSON(response.Success("If the email exists, a reset link will be sent."))
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	type ResetPasswordRequest struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("invalid request"))
	}

	if err := h.svc.ResetPassword(req.Token, req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	return c.JSON(response.Success("Password reset successfully"))
}

func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	router.Post("/login", h.Login)
	router.Post("/register", h.Register)
	router.Post("/refresh", h.Refresh)
	router.Post("/forgot-password", h.ForgotPassword)
	router.Post("/reset-password", h.ResetPassword)
	
	// Protected routes
	protected := router.Group("/", AuthMiddleware())
	protected.Get("/profile", h.GetProfile)
}
