package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)


//Upgrade HTTP to websocket

 var upgrader = websocket.Upgrader{
  //To satisfy the Same-Origin Policy (SOP)
	//checks if the ws is established from the same origin (same protocol,host,port) 
	//not for production tho
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Sender *websocket.Conn
	Content string
}

//thread-safe list of clients

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var mutex = &sync.Mutex{} 

func main(){
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080",nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	//Upgrade GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] Failed to upgrade HTTP to WebSocket:", err)
		return
  }
  log.Println("[INFO] WebSocket connection established")

	defer ws.Close()


	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	//listen for messages from client

	for {
		// _ is message type, msg is actual message but in []byte
		_, msg, err := ws.ReadMessage()

		if err != nil {
			mutex.Lock()
			delete(clients,ws)
			mutex.Unlock()
			break
		}
		log.Println("[INFO] Message recieved from the client")
		//sending the recieved message on broadcast channel but first converting it into string
		broadcast <- Message{
			Sender: ws,
			Content: string(msg),
		}
	}
}


func handleMessages() {
	for{
		// wait for message from any client

		msg := <-broadcast //Blocks until message is received from any client

		//send to all connnected clients
		mutex.Lock()
		for client := range clients {
			if client != msg.Sender {
				err := client.WriteMessage(websocket.TextMessage, []byte(msg.Content))
				if err != nil {
					log.Printf("Write error: %v",err)
					client.Close()
					delete(clients, client)
				}
				log.Println("[INFO] Message sent to all clients")
			}
		}
		mutex.Unlock()
	}
}


