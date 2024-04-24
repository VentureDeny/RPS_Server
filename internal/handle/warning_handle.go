package handle

import (
	"RPS_SERVICE/internal/db"
	datastruct "RPS_SERVICE/internal/struct"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ClientMessage2 struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type ClientMessage struct {
	Action string                 `json:"action"`
	Data   datastruct.WarningData `json:"data"`
}

func HandleWarningWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// 在连接建立时发送现有的警告数据
	sendExistingWarnings(conn)

	for {

		// 这里可以添加处理消息的逻辑
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket closed with error: %v", err)
			break
		}

		var msg ClientMessage2
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		switch msg.Action {
		case "deleteError":
			//log.Println(msg.Data)
			if timestamp, ok := msg.Data.(string); ok {
				deleteWarning(timestamp)
				log.Println(timestamp)
			} else {
				log.Println("Delete Error")
			}
		case "getErrors":
			sendExistingWarnings(conn)
		default:
			log.Printf("Unsupported action received: %s", msg.Action)
		}
	}
}

func SaveWarningToDB(data datastruct.WarningData) error {
	// 假设有一个数据库操作函数实现了存储逻辑
	return db.SaveWarning(data)
}
func deleteWarning(timestamp string) {
	err := db.DeleteWarningByTimestamp(timestamp)
	if err != nil {
		log.Printf("Failed to delete warning: %v", err)
	}
}
func sendExistingWarnings(conn *websocket.Conn) {
	warnings, err := db.FetchAllWarnings()
	if err != nil {
		log.Printf("Failed to fetch warnings: %v", err)
		return
	}
	err = conn.WriteJSON(warnings)
	if err != nil {
		log.Printf("Failed to send warnings: %v", err)
	}
}

// 更新数据库代码
func UpdateWarning(data datastruct.WarningData) error {
	return db.SaveWarning(data)
}
