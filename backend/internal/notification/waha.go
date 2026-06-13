package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WahaProvider struct {
	Endpoint string
	ChatID   string
}

func NewWahaProvider(endpoint, chatID string) *WahaProvider {
	return &WahaProvider{
		Endpoint: endpoint,
		ChatID:   chatID,
	}
}

func (p *WahaProvider) SendNotification(subject, message string) error {
	if p.Endpoint == "" || p.ChatID == "" {
		return fmt.Errorf("waha configuration missing")
	}

	text := fmt.Sprintf("*%s*\n%s", subject, message)
	
	payload := map[string]interface{}{
		"chatId": p.ChatID,
		"text":   text,
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(p.Endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("waha API returned status: %d", resp.StatusCode)
	}

	return nil
}
