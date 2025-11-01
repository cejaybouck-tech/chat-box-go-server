package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	// Allow cross-origin for development (remove in production)
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == os.Getenv("APP_ORIGIN")
	},
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte //raw JSON bytes
	username string
}

type AuthMessage struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Type    string `json:"type"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Message struct {
	Type 		string `json:"type"`
	Username	string `json:"username"`
	Message 	string `json:"message"`
}

type UpdateUsersMessage struct {
	Type	string	`json:"type"`
	Users	[]string `json:"users"`
}

func ServeWebSockets(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w,r, nil);

	if err != nil {
		log.Printf("Upgrade failed: %v", err)
		return
	}

	//read auth message
	_, msg, err := conn.ReadMessage()
	if err != nil {
		
		writeAuthError(conn, "Server failed to read message: "+ err.Error())
		conn.Close()
		return;
	}

	//format json
	var authMsg AuthMessage
	if err := json.Unmarshal(msg, &authMsg); err != nil {
		writeAuthError(conn, "Invalid Message format")
		return;
	}

	//validate message type
	if authMsg.Type != "authenticate" {
		writeAuthError(conn, "Expected 'authenticate' as message type")
		return;
	}

	//check if user exists and if password is correct, will create user if one doesnt exist
	if err := IsValidLogin(hub, authMsg.Username, authMsg.Password); err != nil {
			writeAuthError(conn, err.Error())
			return;
		}
		
	client := &Client {
		hub: hub,
		conn: conn,
		send: make(chan []byte, 256),
		username: authMsg.Username,
	}

	client.hub.register <- client

	//add send message to update connected client with all currently logged users
	users := make([]string, 0, len(hub.clients))
	users = append(users, authMsg.Username)

	for userClient := range hub.clients {
		users = append(users, userClient.username)
	}

	//Send update users message
	if err := conn.WriteJSON(UpdateUsersMessage {
		Type: "users",
		Users: users,
	}); err != nil {
			log.Println("Write error:", err)
	}

	go client.writePump()
	go client.readPump()
	
}

func writeAuthError(conn *websocket.Conn, msg string) {
	conn.WriteJSON(AuthResponse{
		Type:    "auth_response",
		Success: false,
		Message: msg,
	})
	conn.Close()
	log.Printf("Upgrade failed: %v", msg);
}

func (c *Client) writePump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	
	for msg := range c.send{
		var newMessage Message
		if err := json.Unmarshal(msg, &newMessage); err != nil {
			log.Println("json unmarshal error:", err)
			continue
		}

		if err := c.conn.WriteJSON(newMessage); err != nil {
			log.Println("Write error:", err)
		}
	}

}

func (c *Client) readPump() {

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512) //may need to be bigger

	for {

		_, raw, err := c.conn.ReadMessage()

		if (err != nil) {
			log.Println("Failed to read message:", err)
			break;
		}

		c.hub.broadcast <- raw

	}
}