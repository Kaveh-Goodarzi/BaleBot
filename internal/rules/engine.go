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
var admins = map[int64]bool{681706497: true} // ID ادمین‌ها

func Handle(chatID, userID, msgID int64, text string, replyTo int64,
	send SendFunc, del DeleteFunc, ban BanFunc) {

	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return
	}

	// بررسی mute فعال
	var until int64
	err := db.DB.QueryRow("SELECT until FROM muted WHERE user_id=?", userID).Scan(&until)
	if err == nil && time.Now().Unix() < until {
		del(chatID, msgID)
		return
	}

	// بررسی کلمات زشت
	for _, w := range badWords {
		if strings.Contains(text, w) {
			untilTime := time.Now().Add(10 * time.Minute).Unix()
			db.DB.Exec("INSERT OR REPLACE INTO muted VALUES(?,?)", userID, untilTime)
			del(chatID, msgID)
			send(chatID, "کاربر ۱۰ دقیقه mute شد")
			return
		}
	}

	// دستورات ادمین
	if admins[userID] && replyTo != 0 {
		switch {
		case strings.HasPrefix(text, "mute"):
			untilTime := time.Now().Add(10 * time.Minute).Unix()
			db.DB.Exec("INSERT OR REPLACE INTO muted VALUES(?,?)", replyTo, untilTime)
			send(chatID, "کاربر mute شد")

		case strings.HasPrefix(text, "ban"):
			ban(chatID, replyTo)
			send(chatID, "کاربر ban شد")
		}
	}
}
