package handle

import (
	"RPS_SERVICE/internal/db"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var fleetClients = make(map[*websocket.Conn]bool)

type FleetAction struct {
	Action string       `json:"action"`
	Data   db.FleetData `json:"data"` // Use the FleetData type from the db package
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
