package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)


type Message struct{
	Sender string `json:"sender"`
	Content string `json:"content"`
}


func main(){
	//Prompt for username
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")

	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)


	// connect to WebSocket server

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("[ERROR] Connection error:",err)
	}

	defer conn.Close()


	// Read message from server (in background)

	go func(){
		for {
			_, msg, err := conn.ReadMessage()

			if err != nil {
				log.Println("[ERROR] Read error:",err)
				return
			}
			
			fmt.Printf("\r<< %s\n>> ", string(msg))
		}
	}()


	//Read input from user and send to server

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(">> ")

	for scanner.Scan() {
		text := scanner.Text()
		if strings.TrimSpace(text) == "/exit" {
			fmt.Println("Exiting chat...")
			break
		}

		message := Message{
			Sender: username,
			Content: text,
		}

		msgBytes, err := json.Marshal(message)
		if err != nil {
			log.Println("[ERROR] Marshal error", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Println("[ERROR] Write error:",err)
			break
		}
		fmt.Print(">> ")

	}

}
