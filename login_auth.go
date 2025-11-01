package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

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

func IsValidLogin(hub *Hub, username string, password string) error {
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

	for client := range hub.clients {
		if(client.username == username) {
			return errors.New("User is already logged in")
		}
	}

	if (!ok) {

	}

	return nil
}