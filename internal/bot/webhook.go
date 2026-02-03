package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"bale-moderator-bot/config"
	"bale-moderator-bot/internal/db"
	"bale-moderator-bot/internal/rules"
)

type Update struct {
	Message struct {
		MessageID int64  `json:"message_id"`
		Text      string `json:"text"`
		Chat      struct {
			ID int64 `json:"id"`
		}
		From struct {
			ID int64 `json:"id"`
		}
		ReplyToMessage *struct {
			From struct {
				ID int64 `json:"id"`
			}
		} `json:"reply_to_message"`
	} `json:"message"`
}

func Start() {
	db.Init()
	setWebhook()

	http.HandleFunc("/webhook", handle)
	fmt.Println("Bot running on :8080")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

func setWebhook() {
	call("setWebhook", map[string]any{
		"url": config.WebhookURL,
	})
}

func handle(w http.ResponseWriter, r *http.Request) {
	var u Update
	json.NewDecoder(r.Body).Decode(&u)

	replyTo := int64(0)
	if u.Message.ReplyToMessage != nil {
		replyTo = u.Message.ReplyToMessage.From.ID
	}

	rules.Handle(
		u.Message.Chat.ID,
		u.Message.From.ID,
		u.Message.MessageID,
		u.Message.Text,
		replyTo,
		SendMessage,
		DeleteMessage,
		BanUser,
	)
}
