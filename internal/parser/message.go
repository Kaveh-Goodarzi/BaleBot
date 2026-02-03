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
	Type          MessageType
	DurationMinute int
}

type BaleUpdate struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		MessageID int64 `json:"message_id"`
		Text      string
		Chat      struct{ ID int64 `json:"id"` }
		From      struct{ ID int64 `json:"id"` }
		ReplyToMessage *struct {
			From struct{ ID int64 `json:"id"` }
		} `json:"reply_to_message"`
	} `json:"message"`
}

func ParseUpdate(u BaleUpdate) ParsedMessage {
	msg := ParsedMessage{
		ChatID: u.Message.Chat.ID,
		FromUserID: u.Message.From.ID,
		Text: strings.TrimSpace(u.Message.Text),
		Type: Normal,
	}

	if u.Message.ReplyToMessage != nil {
		msg.IsReply = true
		msg.ReplyToUserID = u.Message.ReplyToMessage.From.ID
	}

	parseCommand(&msg)
	return msg
}

func parseCommand(msg *ParsedMessage) {
	parts := strings.Fields(msg.Text)
	if len(parts) == 0 { return }

	switch strings.ToUpper(parts[0]) {
	case "MUTE":
		msg.Type = AdminMute
		if len(parts) > 1 {
			if min, err := strconv.Atoi(parts[1]); err == nil {
				msg.DurationMinute = min
			}
		}
	case "BAN":
		msg.Type = AdminBan
	case "KICK":
		msg.Type = AdminKick
	case "REPORT":
		msg.Type = UserReport
	}
}