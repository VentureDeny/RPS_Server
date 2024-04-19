package handle

import (
	"RPS_SERVICE/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var commonClients = make(map[*websocket.Conn]bool)

func HandleCommonWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New WS Connetced")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	defer func() {
		// 在设备断开连接时执行清理操作
		handleDeviceDisconnection(conn)
		delete(commonClients, conn)
		conn.Close()
	}()
	commonClients[conn] = true
	for {
		// 设置读取超时，如果3秒内没有收到消息，则触发超时错误
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// 超时或其他错误，根据需要处理
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			} else {
				log.Println("read timeout or error:", err)
			}
			handleDeviceDisconnection(conn)
			break // 断开连接
		}
		// 接收消息并根据类型存储到相应的数据库
		ForwardToDataClients(msg)
		//log.Println(msg)
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
		} else if signupData, ok := data["Signup"]; ok {
			handleSignupData(signupData, conn)
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
	//fmt.Println(gpsID, location)
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
		//log.Println("Invalid status data format")
		return
	}
	statusID, okID := statusData["id"].(string)
	battery, okBattery := statusData["battery"].(string)
	MAC, okMAC := statusData["MAC"].(string)
	if !okID || !okBattery || !okMAC {
		//log.Println("Invalid status data received")
		return
	}

	// 调用SaveStatusData存储提取的数据
	db.SaveStatusData(statusID, battery, MAC)
	//fmt.Println(statusID, battery, MAC)
}
func handleSignupData(rawData interface{}, conn *websocket.Conn) {
	signupData, ok := rawData.(map[string]interface{})
	if !ok {
		log.Println("Invalid Signup data format")
		return
	}
	deviceID, ok := signupData["id"].(string)
	if !ok {
		log.Println("Invalid Signup data, missing 'id'")
		return
	}
	connToDeviceID[conn] = deviceID
	// 在这里调用数据库函数添加设备到 onlinedevice 和 alldevice
	db.AddDeviceToOnlineAndAll(deviceID)
	log.Printf("Device %s signed up and added to databases", deviceID)
}

func ForwardToCountClients(message []byte) {
	mu.Lock()         // 在写操作前锁定
	defer mu.Unlock() // 确保函数退出时解锁
	for client := range countClients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("forward error: %v", err)
			client.Close()
			delete(countClients, client)
		}
	}
}
func ForwardToFleetClients(message []byte) {
	mu.Lock()         // 在写操作前锁定
	defer mu.Unlock() // 确保函数退出时解锁
	for client := range fleetClients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("forward error: %v", err)
			client.Close()
			delete(fleetClients, client)
		}
	}
}
