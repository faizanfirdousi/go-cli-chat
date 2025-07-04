package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)



func main(){
	// connect to WebSocket server

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)

	if err != nil {
		log.Fatal("Connection error:",err)
	}

	defer conn.Close()


	// Read message from server (in background)

	go func(){
		for {
			_, msg, err := conn.ReadMessage()

			if err != nil {
				log.Println("Read error:",err)
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
		if text == "/exit" {
			fmt.Println("Exiting...")
			break
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("Write error:",err)
			break
		}
		fmt.Print(">> ")

	}

}
