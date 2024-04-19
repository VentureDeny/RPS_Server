package handle

import (
	"RPS_SERVICE/internal/db"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 在生产环境中应更严格地检查来源
	},
}

func HandleRegisterWS(w http.ResponseWriter, r *http.Request) {
	log.Println("Register WS Endpoint Hit")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}

	var user map[string]string
	json.Unmarshal(msg, &user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user["password"]), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user["username"], hashedPassword)
	if err != nil {
		log.Println("Error inserting new user:", err)
		return
	}

	response := map[string]string{
		"status": "success",
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling response:", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		log.Println("Error sending response:", err)
	}
}

// HandleLoginWS 处理登录的WebSocket连接
func HandleLoginWS(w http.ResponseWriter, r *http.Request) {
	log.Println("Login WS Endpoint Hit")
	conn, err := upgrader.Upgrade(w, r, nil) // 使用之前定义的upgrader进行协议升级
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close() // 确保在函数退出时关闭连接

	// 等待接收一条消息
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}

	// 打印接收到的消息，仅用于调试目的
	log.Printf("Received message: %s", msg)

	// 准备响应消息
	response := map[string]string{
		"status": "success",
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling response:", err)
		return
	}

	// 发送响应消息回客户端
	if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		log.Println("Error sending response:", err)
	}
}
