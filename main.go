package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var offset = 0

func main() {
	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		panic("TG_BOT_TOKEN not set")
	}

	fmt.Println("🤖 Raspberry Pi Telegram Bot started")

	for {
		url := fmt.Sprintf(
			"https://api.telegram.org/bot%s/getUpdates?timeout=30&offset=%d",
			token,
			offset,
		)

		resp, err := http.Get(url)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		var data struct {
			Result []Update `json:"result"`
		}

		json.NewDecoder(resp.Body).Decode(&data)
		resp.Body.Close()

		for _, u := range data.Result {
			offset = u.UpdateID + 1
			handleMessage(token, u)
		}
	}
}

func handleMessage(token string, u Update) {
	if u.Message.Text == "" {
		return
	}

	chatID := u.Message.Chat.ID
	cmd := u.Message.Text

	switch cmd {
	case "/status":
		msg := fmt.Sprintf(
			"📟 Raspberry Pi Status\n\n🌡 CPU: %s\n⏱ Uptime: %s",
			cpuTemp(),
			uptime(),
		)
		sendMessage(token, chatID, msg)

	case "/temp":
		sendMessage(token, chatID, "🌡 "+cpuTemp())

	default:
		sendMessage(token, chatID, "❓ Unknown command")
	}
}
