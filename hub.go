package main

import (
	"encoding/json"
	"log"
)


type Hub struct {
	clients map[*Client]bool 
	broadcast chan []byte //raw JSON bytes
	register chan *Client
	unregister chan *Client
}

type UserJoinMessage struct {
	Type	string	`json:"type"`
	UName	string	`json:"uname"`
}

type UserLeaveMessage struct {
	Type	string	`json:"type"`
	UName	string	`json:"uname"`
}


func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte, 256),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		//handle registers
		case client := <-h.register:
			h.clients[client] = true
			byteSlice, err := json.Marshal(&Message{
				Type: "user_joined",
				Username: client.username,
				Message: "",
			})

			if err != nil {
				log.Printf("Failed to marshal join message: %s", err.Error())
				return;
			}

			h.broadcast <- byteSlice
			log.Printf("Client connected - total: %d", len(h.clients))

		//handle unregisters
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				byteSlice, err := json.Marshal(&Message{
					Type: "user_left",
					Username: client.username,
					Message: "",
				})

				if err != nil {
					log.Printf("Failed to marshal join message: %s", err.Error())
					return;
				}

				h.broadcast <- byteSlice
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected - total: %d", len(h.clients))
			}

		//handle broadcasts to send message
		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					//client is dead / too slow -> drop it
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}