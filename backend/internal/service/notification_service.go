package service

import (
	"fmt"
	"log"
	"time"

	"github.com/velocity/server-monitoring/backend/internal/domain"
	"github.com/velocity/server-monitoring/backend/internal/notification"
	"github.com/velocity/server-monitoring/backend/internal/repository"
)

type NotificationService interface {
	QueueAlert(alert *domain.Alert)
	StartWorker()
	StopWorker()
}

type notificationService struct {
	queue      chan *domain.Alert
	settingSvc SettingService
	logRepo    repository.NotificationLogRepository
	quit       chan struct{}
}

func NewNotificationService(settingSvc SettingService, logRepo repository.NotificationLogRepository) NotificationService {
	return &notificationService{
		queue:      make(chan *domain.Alert, 100),
		settingSvc: settingSvc,
		logRepo:    logRepo,
		quit:       make(chan struct{}),
	}
}

func (s *notificationService) QueueAlert(alert *domain.Alert) {
	select {
	case s.queue <- alert:
	default:
		log.Println("Notification queue full, dropping alert notification")
	}
}

func (s *notificationService) StartWorker() {
	go func() {
		for {
			select {
			case alert := <-s.queue:
				s.processAlert(alert)
			case <-s.quit:
				return
			}
		}
	}()
}

func (s *notificationService) StopWorker() {
	close(s.quit)
}

func (s *notificationService) processAlert(alert *domain.Alert) {
	// TELEGRAM
	tgEnabled, _ := s.settingSvc.GetSetting("telegram_enabled")
	if tgEnabled == "true" {
		tgToken, _ := s.settingSvc.GetSetting("telegram_bot_token")
		tgChat, _ := s.settingSvc.GetSetting("telegram_chat_id")
		if tgToken != "" && tgChat != "" {
			provider := notification.NewTelegramProvider(tgToken, tgChat)
			s.sendWithRetry(provider, "TELEGRAM", alert)
		}
	}

	// WHATSAPP
	waEnabled, _ := s.settingSvc.GetSetting("whatsapp_enabled")
	if waEnabled == "true" {
		waProvider, _ := s.settingSvc.GetSetting("whatsapp_provider")
		waEndpoint, _ := s.settingSvc.GetSetting("whatsapp_endpoint")
		waChat, _ := s.settingSvc.GetSetting("whatsapp_chat_id")
		
		if waProvider == "waha" && waEndpoint != "" && waChat != "" {
			provider := notification.NewWahaProvider(waEndpoint, waChat)
			s.sendWithRetry(provider, "WHATSAPP", alert)
		}
	}
}

func (s *notificationService) sendWithRetry(provider notification.Provider, providerName string, alert *domain.Alert) {
	var subject string
	if alert.State == "RESOLVED" {
		subject = "✅ Server Alert Resolved"
	} else if alert.Severity == "CRITICAL" {
		subject = "🚨 CRITICAL Server Alert"
	} else {
		subject = "⚠️ WARNING Server Alert"
	}

	message := alert.Message
	if alert.Server != nil {
		message = "Server: " + alert.Server.Name + "\n" + message
	} else {
		message = fmt.Sprintf("Server ID: %d\n%s", alert.ServerID, message)
	}

	maxRetries := 3
	backoff := 2 * time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		err := provider.SendNotification(subject, message)
		if err == nil {
			s.logRepo.Create(&domain.NotificationLog{
				Provider: providerName,
				AlertID:  alert.ID,
				Status:   "SUCCESS",
			})
			return
		}
		lastErr = err
		log.Printf("Notification %s attempt %d failed: %v", providerName, i+1, err)
		time.Sleep(backoff)
		backoff *= 2
	}

	// Failed after all retries
	s.logRepo.Create(&domain.NotificationLog{
		Provider: providerName,
		AlertID:  alert.ID,
		Status:   "FAILED",
		Error:    lastErr.Error(),
	})
	log.Printf("Notification %s completely failed for Alert ID %d", providerName, alert.ID)
}
