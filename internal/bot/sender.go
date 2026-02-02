package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"bale-moderator-bot/config"
)

// SendMessage پیام به گروه یا کاربر می‌فرسته
func SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("https://tapi.bale.ai/bot%s/sendMessage", config.BotToken)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	data, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sendMessage failed: %s", resp.Status)
	}

	return nil
}

// DeleteMessage پیام مشخصی را از گروه حذف می‌کند
func DeleteMessage(chatID int64, messageID int64) error {
	url := fmt.Sprintf("https://tapi.bale.ai/bot%s/deleteMessage", config.BotToken)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
	}

	data, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deleteMessage failed: %s", resp.Status)
	}

	return nil
}
