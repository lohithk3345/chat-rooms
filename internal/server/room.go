package server

import (
	// "chat"
	"github.com/google/uuid"
	"sync"
)

type Room struct {
	id      string
	clients map[*Client]bool
	name    string
	mu      sync.Mutex
}

func RandomIDGenerator() string {
	return uuid.NewString()
}

func CreateRoom() (*Room, string) {
	id := RandomIDGenerator()
	return &Room{
		id:      id,
		clients: make(map[*Client]bool),
		name:    "New Room",
	}, id
}
