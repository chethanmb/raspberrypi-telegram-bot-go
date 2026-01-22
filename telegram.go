package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Update struct {
	UpdateID int `json:"update_id"`
	Message  *struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

func sendMessage(token string, chatID int64, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	body := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	b, _ := json.Marshal(body)
	http.Post(url, "application/json", bytes.NewBuffer(b))
}

type CallbackQuery struct {
	ID      string `json:"id"`
	Data    string `json:"data"`
	Message struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

type InlineKeyboard struct {
	InlineKeyboard [][]map[string]string `json:"inline_keyboard"`
}

func sendMessageWithButtons(token string, chatID int64, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	keyboard := InlineKeyboard{
		InlineKeyboard: [][]map[string]string{
			{
				{"text": "📟 Status", "callback_data": "status"},
				{"text": "🌡 Temp", "callback_data": "temp"},
			},
			{
				{"text": "▶️ Start", "callback_data": "start"},
				{"text": "⛔ Stop", "callback_data": "stop"},
			},
		},
	}

	body := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": keyboard,
	}

	b, _ := json.Marshal(body)
	http.Post(url, "application/json", bytes.NewBuffer(b))
}

func answerCallback(token, callbackID string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", token)
	body := map[string]string{"callback_query_id": callbackID}
	b, _ := json.Marshal(body)
	http.Post(url, "application/json", bytes.NewBuffer(b))
}

func sendMenu(token string, chatID int64) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	keyboard := InlineKeyboard{
		InlineKeyboard: [][]map[string]string{
			{
				{"text": "📟 Status", "callback_data": "status"},
				{"text": "🌡 Temp", "callback_data": "temp"},
			},
			{
				{"text": "▶️ Start", "callback_data": "start"},
				{"text": "⛔ Stop", "callback_data": "stop"},
			},
		},
	}

	body := map[string]interface{}{
		"chat_id":      chatID,
		"text":         "Choose an action 👇",
		"reply_markup": keyboard,
	}

	b, _ := json.Marshal(body)
	http.Post(url, "application/json", bytes.NewBuffer(b))
}
