package http

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/repository"
	"github.com/velocity/server-monitoring/backend/pkg/response"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error("missing or invalid token"))
		}

		tokenString := authHeader[7:]
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error("invalid token"))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error("invalid claims"))
		}

		// Set user ID to fiber context
		c.Locals("user_id", uint(claims["user_id"].(float64)))

		return c.Next()
	}
}

func AuditMiddleware(repo repository.AuditLogRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Proceed with the request
		err := c.Next()

		// Log mutative requests
		method := c.Method()
		if method == "POST" || method == "PUT" || method == "DELETE" {
			userID, _ := c.Locals("user_id").(uint) // 0 if unauthenticated
			
			logEntry := &domain.AuditLog{
				UserID:    userID,
				Method:    method,
				Path:      c.Path(),
				IPAddress: c.IP(),
			}

			// We launch this as a goroutine to not block the response
			go func() {
				_ = repo.Create(logEntry)
			}()
		}

		return err
	}
}
