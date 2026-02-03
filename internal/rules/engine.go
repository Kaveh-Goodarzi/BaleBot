package rules

import (
	"strings"
	"time"

	"bale-moderator-bot/internal/db"
)

type SendFunc func(chatID int64, text string)
type DeleteFunc func(chatID, msgID int64)
type BanFunc func(chatID, userID int64)

var badWords = []string{"kir", "jende"}
var admins = map[int64]bool{681706497: true}

func Handle(chatID, userID, msgID int64, text string, replyTo int64,
	send SendFunc, del DeleteFunc, ban BanFunc) {

	// check mute
	var until int64
	err := db.DB.QueryRow("SELECT until FROM muted WHERE user_id=?", userID).Scan(&until)
	if err == nil && time.Now().Unix() < until {
		del(chatID, msgID)
		return
	}

	lower := strings.ToLower(text)

	// profanity auto mute
	for _, w := range badWords {
		if strings.Contains(lower, w) {
			untilTime := time.Now().Add(10 * time.Minute).Unix()
			db.DB.Exec("INSERT OR REPLACE INTO muted VALUES(?,?)", userID, untilTime)
			del(chatID, msgID)
			send(chatID, "کاربر ۱۰ دقیقه mute شد")
			return
		}
	}

	// admin commands
	if admins[userID] && replyTo != 0 {
		switch {
		case strings.HasPrefix(lower, "mute"):
			untilTime := time.Now().Add(10 * time.Minute).Unix()
			db.DB.Exec("INSERT OR REPLACE INTO muted VALUES(?,?)", replyTo, untilTime)
			send(chatID, "کاربر mute شد")

		case strings.HasPrefix(lower, "ban"):
			ban(chatID, replyTo)
			send(chatID, "کاربر ban شد")
		}
	}
}
