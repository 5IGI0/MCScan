package main

type ChatComponent struct {
	Text string `json:"text"`

	Extra []ChatComponent `json:"extra"`
}
