package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TelegramProvider struct {
	BotToken string
	ChatID   string
}

func NewTelegramProvider(botToken, chatID string) *TelegramProvider {
	return &TelegramProvider{
		BotToken: botToken,
		ChatID:   chatID,
	}
}

func (p *TelegramProvider) SendNotification(subject, message string) error {
	if p.BotToken == "" || p.ChatID == "" {
		return fmt.Errorf("telegram configuration missing")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", p.BotToken)
	
	text := fmt.Sprintf("<b>%s</b>\n%s", subject, message)
	
	payload := map[string]interface{}{
		"chat_id":    p.ChatID,
		"text":       text,
		"parse_mode": "HTML",
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status: %d", resp.StatusCode)
	}

	return nil
}
