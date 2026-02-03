package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"log"
	"bale-moderator-bot/config"
)

func call(method string, payload any) error {
	url := fmt.Sprintf("https://tapi.bale.ai/bot%s/%s", config.BotToken, method)
	data, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("%s failed: %s", method, resp.Status)
	}
	return nil
}

func SendMessage(chatID int64, text string) {
	_ = call("sendMessage", map[string]any{
		"chat_id": chatID,
		"text":    text,
	})
}

func DeleteMessage(chatID, msgID int64) {
	_ = call("deleteMessage", map[string]any{
		"chat_id":    chatID,
		"message_id": msgID,
	})
}

func BanUser(chatID, userID int64) {
	_ = call("banChatMember", map[string]any{
		"chat_id": chatID,
		"user_id": userID,
	})
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
