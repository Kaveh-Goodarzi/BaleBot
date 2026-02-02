package parser

import (
	"strconv"
	"strings"
)

type MessageType int

const (
	Normal MessageType = iota
	AdminMute
	AdminBan
	AdminKick
	UserReport
)

type ParsedMessage struct {
	ChatID        int64
	FromUserID    int64
	Text          string

	IsReply       bool
	ReplyToUserID int64

	Type           MessageType
	DurationMinute int
}

type baleUpdate struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		MessageID int64 `json:"message_id"`
		Text      string `json:"text"`

		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`

		From struct {
			ID int64 `json:"id"`
		} `json:"from"`

		ReplyToMessage *struct {
			From struct {
				ID int64 `json:"id"`
			} `json:"from"`
		} `json:"reply_to_message"`
	} `json:"message"`
}


func ParseUpdate(u baleUpdate) ParsedMessage {
	msg := ParsedMessage{
		ChatID:     u.Message.Chat.ID,
		FromUserID: u.Message.From.ID,
		Text:       strings.TrimSpace(u.Message.Text),
		Type:       Normal,
	}

	// آیا ریپلای است؟
	if u.Message.ReplyToMessage != nil {
		msg.IsReply = true
		msg.ReplyToUserID = u.Message.ReplyToMessage.From.ID
	}

	parseCommand(&msg)

	return msg
}

func parseCommand(msg *ParsedMessage) {
	parts := strings.Fields(msg.Text)

	if len(parts) == 0 {
		return
	}

	cmd := parts[0]

	switch cmd {

	// کلمات اختصاصی آگا کاوه!
	case "سوکوت":
		msg.Type = AdminMute
		if len(parts) > 1 {
			min, err := strconv.Atoi(parts[1])
			if err == nil {
				msg.DurationMinute = min
			}
		}

	case "سیشمیش":
		msg.Type = AdminBan

	case "سیکیم":
		msg.Type = AdminKick

	case "ریپبزن":
		msg.Type = UserReport
	}
}
