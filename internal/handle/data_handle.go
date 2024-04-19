package handle

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"RPS_SERVICE/internal/db"

	"github.com/gorilla/websocket"
)

var mu sync.Mutex
var dataClients = make(map[*websocket.Conn]bool)
var connToDeviceID = make(map[*websocket.Conn]string)

func HandleDataWS(w http.ResponseWriter, r *http.Request) {
	log.Println("DataHandle Setup")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			// 每秒执行这些函数发送数据
			FetchAndSendDeviceData()
			//SendOnlineDevicesCount()
		}

	}()

	// 将新的WebSocket连接添加到dataClients
	dataClients[conn] = true

	// 保持连接活跃，直到它断开
	for {
		// NextReader 会阻塞直到收到一个消息或发生错误（比如连接关闭）
		if _, _, err := conn.NextReader(); err != nil {
			log.Printf("WebSocket closed with error: %v", err)
			conn.Close()
			delete(dataClients, conn)
			break // 退出 for 循环
		}
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket closed with error: %v", err)
			break
		}

		var msg ClientMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		switch msg.Action {
		case "logWarning":
			handleWarning(msg.Data)
		default:
			log.Printf("Unsupported action received: %s", msg.Action)
		}
		// 这里可以添加处理消息的逻辑
	}
}

func handleWarning(data db.WarningData) {
	log.Printf("Received warning: %+v", data)
	// 在这里调用数据库操作函数，将警告信息存储到数据库中
	err := SaveWarningToDB(data)
	if err != nil {
		log.Printf("Error saving warning to DB: %v", err)
	}
}

func ForwardToDataClients(message []byte) {
	mu.Lock()         // 在写操作前锁定
	defer mu.Unlock() // 确保函数退出时解锁
	for client := range dataClients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("forward error: %v", err)
			client.Close()
			delete(dataClients, client)
		}
	}
}

func handleDeviceDisconnection(conn *websocket.Conn) {
	// 使用 conn 查找设备ID
	deviceID, exists := connToDeviceID[conn]
	if !exists {
		log.Println("Device ID for the disconnected device not found")
		return
	}

	// 从 onlinedevice 数据库中移除设备
	db.RemoveDeviceFromOnline(deviceID)
	log.Printf("Device %s removed from online devices database", deviceID)

	// 从映射中移除此连接
	delete(connToDeviceID, conn)
}
