package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // 连接的客户端
var broadcast = make(chan []byte)            // 广播通道
var upgrader = websocket.Upgrader{}          // 将HTTP服务器升级到WebSocket协议
var mutex = &sync.Mutex{}                    // 保护clients

// HandleConnections 升级HTTP到WebSocket协议，并处理连接
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}
	defer ws.Close()

	// 注册新的客户端
	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			break
		}
		// 接收到消息，放入广播通道
		broadcast <- msg
	}
}

// HandleMessages 监听广播通道中的消息并广播给所有客户端
func HandleMessages() {
	for {
		// 从广播通道中获取消息
		msg := <-broadcast

		mutex.Lock()
		// 发送消息给所有连接的客户端
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Printf("Error writing message: %v\n", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
