package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool) // 连接的客户端
	broadcast = make(chan []byte)              // 广播通道
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	mutex sync.Mutex
)

// HandleConnections 处理WebSocket连接
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// 升级初始GET请求到一个WebSocket连接
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// 注册新的客户端
	clients[ws] = true

	// 监听新的消息
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// 发送消息到广播通道
		broadcast <- msg
	}
}

// HandleMessages 监听广播通道中的消息并广播给所有客户端
func HandleMessages() {
	for {
		msg := <-broadcast
		log.Printf("Got message: %s\n", msg)
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
