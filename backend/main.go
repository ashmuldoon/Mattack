package backend

import (
	"backend/server/upgrader"
	"log"
	"net/http"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	defer conn.Close()

	log.Println("client connected")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error (client likely disconnected):", err)
			break
		}

		log.Printf("received: %s", message)

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}

	log.Println("client disconnected")
}

func main() {
	http.HandleFunc("/ws", echoHandler)

	log.Println("server starting on :8080, connect to ws://localhost:8080/ws")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("server error:", err)
	}
}
