package handle

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var countClients = make(map[*websocket.Conn]bool)

func HandleCountWS(w http.ResponseWriter, r *http.Request) {
	log.Println("CountHandle Setup")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			// 每秒执行这些函数发送数据
			SendOnlineDevicesCount()
		}
	}()

	// 将新的WebSocket连接添加到dataClients
	countClients[conn] = true

	// 保持连接活跃，直到它断开
	for {
		// NextReader 会阻塞直到收到一个消息或发生错误（比如连接关闭）
		if _, _, err := conn.NextReader(); err != nil {
			log.Printf("WebSocket closed with error: %v", err)
			conn.Close()
			delete(countClients, conn)
			break // 退出 for 循环
		}
		// 这里可以添加处理消息的逻辑
	}
}
