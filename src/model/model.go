package model

import "time"

type FriendRequest struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Status    string    `json:"status"` // "pending", "accepted", "rejected"
	CreatedAt time.Time `json:"created_at"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}
