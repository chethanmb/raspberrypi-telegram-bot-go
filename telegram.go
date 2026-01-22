package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
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
