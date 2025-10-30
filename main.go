package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

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

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	// Allow cross-origin for development (remove in production)
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Printf("Origin: %s", r.Header.Get("Origin"))
		return origin == os.Getenv("APP_ORIGIN")
	},
}

const userPath = "users.json"

/* Checks for a valid users file, if file doesnt exist create new users file */
func checkValidUserFile() error {
	_, err := os.Stat(userPath);

	if(os.IsNotExist(err)) {
		empty := map[string]string{"Test":"butsandsuch9874123"}
		data, err := json.MarshalIndent(empty, "", "  ")
		if err != nil {
			log.Printf("Json Marshal Failed")
			return err
		}
		err = os.WriteFile(userPath, data, 0644)
		if err != nil {
			log.Printf("File creation failed")
			return err
		}
		log.Printf("File Created Successfully")
	}

	return nil
}

/* read users file and return all users saved so far */
func readUsers() (map[string]string, error) {
	err := checkValidUserFile()
	if (err != nil) {
		return nil, err
	}

	file, err := os.ReadFile(userPath)
	if(err != nil) {
		log.Printf("Error reading users.json file")
		return nil, err
	}

	var users map[string]string
	err = json.Unmarshal(file, &users)
	if(err != nil) {
		log.Printf("Error unmarshalling users.json")
		return nil, err
	}

	return  users, err
}

func addUserToFile(users map[string]string, username string, password string) error{

	users[username] = password;
	data, err := json.MarshalIndent(users, "", "  ");

	if(err != nil) {
		log.Printf("failed to marshel user: %v", username)
		return err
	}

	return os.WriteFile(userPath, data, 0644)
}

func isValidLogin(username string, password string) error {
	users, err := readUsers()

	if(err != nil) {
		return err
	}

	userPassword, ok := users[username]
	
	if(!ok) {
		if err := addUserToFile(users, username, password); err != nil {
			return err;
		}
		return nil
	}

	if (password != userPassword) {
		return errors.New("password is incorrect")
	}

	return nil
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w,r, nil);

	if err != nil {
		log.Printf("Upgrade failed: %v/n", err)
		return
	}

	defer conn.Close();

	//Authentication message
	msgType, msg, err := conn.ReadMessage()

	if err != nil {
		log.Printf("message type: %v\n", msgType)
		log.Printf("Read error: %v\n", err)
		return;
	}

	var authMsg AuthMessage
	//unmarshal the json format
	if err := json.Unmarshal(msg, &authMsg); err != nil {
		log.Println("Unmarshal error:", err)
		conn.WriteJSON(AuthResponse{Type: "auth_response", Success: false, Message: "Invalid message format"})
		return;
	}

	//validate login
	if authMsg.Type == "authenticate" {

		if err := isValidLogin(authMsg.Username, authMsg.Password); err != nil {
			log.Printf("Validation error: %v", err)
			conn.WriteJSON(AuthResponse{Type: "auth_response", Success: false, Message: err.Error()})
			return;
		}
		
		//todo: add user to logged in array
		//continue with authenticated connection
		for {
		messageType, message, err := conn.ReadMessage();

		
		fmt.Printf("Read Message: %v\n", message)
		if err != nil {
			log.Printf("Read Failed: %v", err);
			break;
		}

		
		fmt.Println("Write Message")
		if writeErr := conn.WriteMessage(messageType, message); writeErr != nil {
			log.Printf("Write Failed: %v", writeErr);
			break;
		}

		fmt.Printf("message-type: %v, message: %v", messageType, message);
		}
	}

	
}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err);
		return;
	}
	http.HandleFunc("/chat", messageHandler);

	// Start the HTTP server
	log.Printf("Starting server at %s", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("ListenAndServer: %v", err)
	}

}