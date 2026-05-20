package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type outgoingMsg struct {
	From    string `json:"from"`
	Body    string `json:"body"`
	Channel string `json:"channel"`
}

type BotReply struct {
	Reply   string      `json:"reply"`
	Actions []BotAction `json:"actions"`
}

type BotAction struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func sendMessage(baseURL, from, body, channel string) (*BotReply, error) {
	payload, _ := json.Marshal(outgoingMsg{From: from, Body: body, Channel: channel})
	resp, err := http.Post(
		baseURL+"/api/towershare/message",
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var r BotReply
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, fmt.Errorf("bad response: %s", string(raw))
	}
	return &r, nil
}
