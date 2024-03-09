package router

import (
	"RPS_SERVICE/pkg/db"
	"encoding/json"
	"fmt"
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
var commonClients = make(map[*websocket.Conn]bool)

// SetupRoutes 配置WebSocket路由
func SetupRoutes() {
	log.Println("Router Setup")
	http.HandleFunc("/data", handleDataWS)
	http.HandleFunc("/common", handleCommonWS) // 使用相同的处理函数来处理/gps和/rps
}

// handleDataWS 处理连接到/data的WebSocket客户端
func handleDataWS(w http.ResponseWriter, r *http.Request) {
	log.Println("DataHandle Setup")
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

func handleCommonWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New WS Connetced")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()
	commonClients[conn] = true

	// 接收消息并根据类型存储到相应的数据库
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		forwardToDataClients(msg)

		// 使用interface{}来处理不同格式的数据
		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Println("json unmarshal error:", err)
			continue
		}

		// 根据数据类型处理GPS或RPS数据
		if gpsData, ok := data["GPS"]; ok {
			handleGPSData(gpsData)
		} else if rpsData, ok := data["RPS"]; ok {
			handleRPSData(rpsData)
		} else if statusData, ok := data["Status"]; ok {
			handleStatusData(statusData)
		} else {
			log.Println("Invalid data type received")
		}
	}
	// 保持连接活跃，直到它断开
	for {
		if _, _, err := conn.NextReader(); err != nil {
			conn.Close()
			delete(commonClients, conn)
			break
		}
	}
}

func handleGPSData(rawData interface{}) {
	log.Println("GPSHandle Setup")
	// 转换并处理GPS数据
	gpsData, ok := rawData.(map[string]interface{})
	if !ok {
		log.Println("Invalid GPS data format")
		return
	}
	gpsID, okID := gpsData["id"].(string)
	location, okLocation := gpsData["location"].(string)
	if !okID || !okLocation {
		log.Println("Invalid GPS data received")
		return
	}

	// 直接调用SaveGPSData存储提取的数据，无需分割location
	db.SaveGPSData(gpsID, location)
	fmt.Println(gpsID, location)
}

func handleRPSData(rawData interface{}) {
	// 转换并处理RPS数据
	rpsData, ok := rawData.(map[string]interface{})
	if !ok {
		log.Println("Invalid RPS data format")
		return
	}
	for id, coords := range rpsData {
		coordsMap := coords.(map[string]interface{})
		x := int(coordsMap["x"].(float64))
		y := int(coordsMap["y"].(float64))
		db.SaveRPSData(id, x, y) // 正确传递参数
		fmt.Println(id, x, y)
	}
}
func handleStatusData(rawData interface{}) {
	statusData, ok := rawData.(map[string]interface{})
	if !ok {
		log.Println("Invalid status data format")
		return
	}
	statusID, okID := statusData["id"].(string)
	battery, okBattery := statusData["battery"].(string)
	MAC, okMAC := statusData["MAC"].(string)
	if !okID || !okBattery || !okMAC {
		log.Println("Invalid status data received")
		return
	}

	// 调用SaveStatusData存储提取的数据
	db.SaveStatusData(statusID, battery, MAC)
	fmt.Println(statusID, battery, MAC)
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

/* func GetCommonClientsCount() int {
	return len(commonClients)
}
*/
