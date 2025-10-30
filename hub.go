package main

import (
	"log"
)


type Hub struct {
	clients map[*Client]bool 
	broadcast chan []byte //raw JSON bytes
	register chan *Client
	unregister chan *Client
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
			log.Printf("Client connected - total: %d", len(h.clients))

		//handle unregisters
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
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