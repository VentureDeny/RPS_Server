package handle

import (
	"RPS_SERVICE/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/websocket"
)

type FleetAction struct {
	Action string       `json:"action"`
	Data   db.FleetData `json:"data"` // Use the FleetData type from the db package
}

// 定义接收到的警告数据结构

// 定义从前端接收的消息格式
type ClientMessage struct {
	Action string         `json:"action"`
	Data   db.WarningData `json:"data"`
}
type ClientMessage2 struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

var mu sync.Mutex // 创建一个互斥锁变量
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 在生产环境中应更严格地检查来源
	},
}

// dataClients 存储所有连接到/data的WebSocket客户端
var dataClients = make(map[*websocket.Conn]bool)
var commonClients = make(map[*websocket.Conn]bool)
var connToDeviceID = make(map[*websocket.Conn]string)
var fleetClients = make(map[*websocket.Conn]bool)
var countClients = make(map[*websocket.Conn]bool)

// handleDataWS 处理连接到/data的WebSocket客户端
func HandleDataWS(w http.ResponseWriter, r *http.Request) {
	log.Println("DataHandle Setup")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			// 每秒执行这些函数发送数据
			//FetchAndSendDeviceData()
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

func SaveWarningToDB(data db.WarningData) error {
	// 假设有一个数据库操作函数实现了存储逻辑
	return db.SaveWarning(data)
}

func GetAllDevicesHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := db.GetAllDevices() // 假设这个函数已经存在于你的 db 包中，用于获取所有设备的ID
	if err != nil {
		http.Error(w, "Failed to get devices", http.StatusInternalServerError)
		return
	}

	// 设置响应类型为JSON
	w.Header().Set("Content-Type", "application/json")
	// 将设备列表编码为JSON并发送响应
	json.NewEncoder(w).Encode(devices)
}
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

func HandleFleetWS(w http.ResponseWriter, r *http.Request) {
	log.Println("Fleet WebSocket endpoint hit")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// 发送现有车队数据到前端
	sendExistingFleetData(conn)

	// 持续监听消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket closed with error: %v", err)
			break
		}

		var fa FleetAction
		err = json.Unmarshal(message, &fa)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		switch fa.Action {
		case "createFleet":
			createFleet(conn, fa.Data)
		case "updateFleet":
			updateFleet(conn, fa.Data)
		case "deleteFleet":
			deleteFleet(conn, fa.Data.ID)
		case "getFleets":
			sendExistingFleetData(conn)

		default:
			log.Printf("Unsupported fleet action: %s", fa.Action)
		}
	}
}

// 发送现有车队数据
func sendExistingFleetData(conn *websocket.Conn) {
	fleets, err := db.GetFleets() // 假设有这样一个函数来获取所有车队数据
	if err != nil {
		log.Printf("Error fetching fleets: %v", err)
		return
	}
	// 将获取的车队数据发送给前端
	conn.WriteJSON(fleets)

	device, err := db.GetAllDevices()
	if err != nil {
		log.Printf("Error fetching fleets: %v", err)
		return
	}
	err = conn.WriteJSON(device)
	if err != nil {
		log.Printf("Error sending fleet data: %v", err)
		return
	}
}

func createFleet(conn *websocket.Conn, data db.FleetData) {
	log.Printf("Creating fleet: %+v", data)
	err := db.CreateFleet(data.Name, data.Vehicles)
	if err != nil {
		log.Printf("Error creating fleet: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error creating fleet"))
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Fleet created successfully"))
}

func updateFleet(conn *websocket.Conn, data db.FleetData) {
	log.Printf("Updating fleet: %+v", data)

	err := db.UpdateFleet(data.ID, data.Name, data.Vehicles)
	if err != nil {
		log.Printf("Error updating fleet: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error updating fleet"))
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Fleet updated successfully"))
}

func deleteFleet(conn *websocket.Conn, id string) {
	log.Printf("Deleting fleet with ID: %s", id)
	err := db.DeleteFleet(id)
	if err != nil {
		log.Printf("Error deleting fleet: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error deleting fleet"))
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Fleet deleted successfully"))
}
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

// forwardToDataClients 将收到的消息转发到所有连接到/data的客户端
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

// handleDeviceDisconnection 在设备断开连接时调用

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

// HandleRegisterWS 处理注册的WebSocket连接
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

func deleteWarning(timestamp string) {
	err := db.DeleteWarningByTimestamp(timestamp)
	if err != nil {
		log.Printf("Failed to delete warning: %v", err)
	}
}
