package router

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 在生产环境中应更严格地检查来源
	},
}

// dataClients 存储所有连接到/data的WebSocket客户端
var dataClients = make(map[*websocket.Conn]bool)

// SetupRoutes 配置WebSocket路由
func SetupRoutes() {
	http.HandleFunc("/data", handleDataWS)
	http.HandleFunc("/gps", handleGPSRPSWS) // 使用相同的处理函数来处理/gps和/rps
	http.HandleFunc("/rps", handleGPSRPSWS) // 因为它们的处理逻辑相同
}

// handleDataWS 处理连接到/data的WebSocket客户端
func handleDataWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// 将新的WebSocket连接添加到dataClients
	dataClients[conn] = true

	// 保持连接活跃，直到它断开
	for {
		if _, _, err := conn.NextReader(); err != nil {
			conn.Close()
			delete(dataClients, conn)
			break
		}
	}
}

// handleGPSRPSWS 处理/gps和/rps路由的WebSocket连接
func handleGPSRPSWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// 接收消息并转发到所有/data的客户端
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		forwardToDataClients(msg)
	}
}

// forwardToDataClients 将收到的消息转发到所有连接到/data的客户端
func forwardToDataClients(message []byte) {
	for client := range dataClients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("forward error: %v", err)
			client.Close()
			delete(dataClients, client)
		}
	}
}
