package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var offset = 0

const (
	TEMP_LIMIT          = 50.0
	TEMP_CHECK_INTERVAL = 5 * time.Second
	ALERT_INTERVAL      = 30 * time.Second
)

var (
	alertsEnabled = true
	lastAlertTime time.Time
	lastChatID    int64
)

func main() {
	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		panic("TG_BOT_TOKEN not set")
	}

	fmt.Println("🤖 Raspberry Pi Telegram Bot started")

	startTempMonitor(token)

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
			if u.Message != nil {
				handleMessage(token, u)
			}
			if u.CallbackQuery != nil {
				handleCallback(token, u.CallbackQuery)
			}

		}
	}
}

func handleMessage(token string, u Update) {
	if u.Message.Text == "" {
		return
	}

	chatID := u.Message.Chat.ID
	lastChatID = chatID

	switch u.Message.Text {

	case "/start":
		alertsEnabled = true
		sendMenu(token, chatID)

	case "/stop":
		alertsEnabled = false
		sendMessage(token, chatID, "⛔ Alerts *DISABLED*")

	case "/temp":
		sendMessage(token, chatID, fmt.Sprintf("🌡 %.1f°C", cpuTempValue()))

	case "/cpu":
		sendMessage(token, chatID, "🧠 CPU Usage: "+cpuUsage())

	case "/ram":
		sendMessage(token, chatID, "🧮 RAM: "+ramInfo())

	case "/uptime":
		sendMessage(token, chatID, "⏱ Uptime: "+uptime())

	case "/throttle":
		sendMessage(token, chatID, "⚡ Throttle: "+throttled())

	case "/status":
		sendStatus(token, chatID)

	case "/menu":
		sendMenu(token, chatID)

	default:
		sendMessage(token, chatID, "❓ Unknown command")
	}
}

func startTempMonitor(token string) {
	go func() {
		ticker := time.NewTicker(TEMP_CHECK_INTERVAL)
		defer ticker.Stop()

		for range ticker.C {

			if !alertsEnabled || lastChatID == 0 {
				continue
			}

			temp := cpuTempValue()
			if temp < TEMP_LIMIT {
				continue
			}

			if time.Since(lastAlertTime) < ALERT_INTERVAL {
				continue
			}

			msg := fmt.Sprintf(
				"🔥 *High Temperature Alert!*\n\n🌡 Temp: %.1f°C\n⚠️ Limit: %.1f°C",
				temp,
				TEMP_LIMIT,
			)

			sendMessage(token, lastChatID, msg)
			lastAlertTime = time.Now()
		}
	}()
}

func handleCallback(token string, cb *CallbackQuery) {
	chatID := cb.Message.Chat.ID
	lastChatID = chatID

	switch cb.Data {

	case "status":
		sendStatus(token, chatID)

	case "temp":
		sendMessage(token, chatID,
			fmt.Sprintf("🌡 Temperature: %.1f°C", cpuTempValue()),
		)

	case "start":
		alertsEnabled = true
		sendMessage(token, chatID, "✅ Temperature alerts ENABLED")

	case "stop":
		alertsEnabled = false
		sendMessage(token, chatID, "⛔ Temperature alerts DISABLED")

	}

	answerCallback(token, cb.ID)
}

func sendStatus(token string, chatID int64) {
	alertStatus := "⛔ Disabled"
	if alertsEnabled {
		alertStatus = "✅ Enabled"
	}

	msg := fmt.Sprintf(
		"📟 Raspberry Pi Status\n\n"+
			"🌡 Temp: %.1f°C\n"+
			"🧠 CPU: %s\n"+
			"🧮 RAM: %s\n"+
			"⏱ Uptime: %s\n"+
			"⚡ Throttle: %s\n"+
			"🚨 Alerts: %s\n",
		cpuTempValue(),
		cpuUsage(),
		ramInfo(),
		uptime(),
		throttled(),
		alertStatus,
	)

	sendMessageWithButtons(token, chatID, msg)
}
