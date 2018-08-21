package main

import (
	"time"
)

// message represent single message
type message struct {
	Name      string    `json:"name,omitempty"`
	Message   string    `json:"message,omitempty"`
	When      time.Time `json:"when,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
}
