package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/velocity/server-monitoring/backend/internal/agent"
	"github.com/velocity/server-monitoring/backend/internal/alert"
	"github.com/velocity/server-monitoring/backend/internal/config"
	delivery "github.com/velocity/server-monitoring/backend/internal/delivery/http"
	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/repository"
	"github.com/velocity/server-monitoring/backend/internal/scheduler"
	"github.com/velocity/server-monitoring/backend/internal/service"
)

func main() {
	_ = godotenv.Load()

	// Connect to Database
	config.ConnectDB()

	// Run Migrations
	err := config.DB.AutoMigrate(
		&domain.User{}, &domain.Server{}, &domain.Metric{},
		&domain.Alert{}, &domain.Setting{}, &domain.NotificationLog{},
		&domain.PasswordResetToken{}, &domain.AuditLog{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	// Rate Limiter Middleware
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too many requests, please try again later."})
		},
	}))

	// Initialize dependencies
	userRepo := repository.NewUserRepository(config.DB)
	resetRepo := repository.NewPasswordResetRepository(config.DB)
	userSvc := service.NewUserService(userRepo, resetRepo)
	authHandler := delivery.NewAuthHandler(userSvc)

	settingRepo := repository.NewSettingRepository(config.DB)
	settingSvc := service.NewSettingService(settingRepo)
	settingHandler := delivery.NewSettingHandler(settingSvc)

	notifLogRepo := repository.NewNotificationLogRepository(config.DB)
	notifSvc := service.NewNotificationService(settingSvc, notifLogRepo)
	notifSvc.StartWorker()
	defer notifSvc.StopWorker()

	alertRepo := repository.NewAlertRepository(config.DB)
	alertEngine := alert.NewEngine(alertRepo, notifSvc)
	alertSvc := service.NewAlertService(alertRepo)
	alertHandler := delivery.NewAlertHandler(alertSvc)

	serverRepo := repository.NewServerRepository(config.DB)
	serverSvc := service.NewServerService(serverRepo, alertRepo)
	serverHandler := delivery.NewServerHandler(serverSvc)

	metricRepo := repository.NewMetricRepository(config.DB)
	metricSvc := service.NewMetricService(metricRepo, alertEngine)
	metricHandler := delivery.NewMetricHandler(metricSvc, serverSvc)

	auditLogRepo := repository.NewAuditLogRepository(config.DB)

	statusChecker := scheduler.NewStatusChecker(serverRepo, metricRepo, alertEngine)
	statusChecker.Start()
	defer statusChecker.Stop()

	selfMonitor := agent.NewSelfMonitor(metricSvc, serverSvc)
	selfMonitor.Start()
	defer selfMonitor.Stop()

	dashboardHandler := delivery.NewDashboardHandler(serverSvc, metricSvc)

	// API Routing
	api := app.Group("/api/v1")
	// Apply Audit Middleware globally to api routes
	api.Use(delivery.AuditMiddleware(auditLogRepo))

	authHandler.RegisterRoutes(api.Group("/auth"))
	metricHandler.RegisterRoutes(api) // Agent metrics route

	protected := api.Group("/", delivery.AuthMiddleware())
	serverHandler.RegisterRoutes(protected.Group("/servers"))
	dashboardHandler.RegisterRoutes(protected.Group("/dashboard"))
	alertHandler.RegisterRoutes(protected.Group("/alerts"))
	settingHandler.RegisterRoutes(protected.Group("/settings"))

	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, err := config.DB.DB()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "ERROR", "database": "DISCONNECTED"})
		}
		if err := sqlDB.Ping(); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "ERROR", "database": "DISCONNECTED"})
		}
		return c.JSON(fiber.Map{"status": "OK", "database": "OK"})
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8001"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
