package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bale-moderator-bot/config"
	"bale-moderator-bot/internal/parser"
)

var offset int64 = 0

func StartPolling() {
	for {
		updates, err := getUpdates()
		if err != nil {
			fmt.Println("error getting updates:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, u := range updates {
			offset = u.UpdateID + 1

			parsed := parser.ParseUpdate(u)

			fmt.Printf("Parsed Message: %+v\n", parsed)
		}
	}
}

func getUpdates() ([]parser.BaleUpdate, error) {
	url := fmt.Sprintf(
		"https://tapi.bale.ai/bot%s/getUpdates?offset=%d",
		config.BotToken,
		offset,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result []parser.BaleUpdate `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Result, nil
}
