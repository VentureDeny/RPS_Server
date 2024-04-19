package handle

import (
	"RPS_SERVICE/internal/db"
	"encoding/json"
	"log"
)

type DeviceInfo struct {
	DeviceID     string `json:"device_id"`
	Location     string `json:"location,omitempty"` // omitempty 表示如果 Location 为空，则不包含这个字段
	BatteryLevel string `json:"battery_level,omitempty"`
	MACAddress   string `json:"mac_address,omitempty"`
}

func FetchAndSendDeviceData() {
	// 查询设备数据
	devices, err := db.QueryDeviceData()
	if err != nil {
		log.Fatal(err)
	}

	// 将结构体切片转换为 JSON
	jsonData, err := json.Marshal(devices)
	if err != nil {
		log.Fatal(err)
	}

	ForwardToDataClients(jsonData)
}

func SendOnlineDevicesCount() {
	count, err := db.GetOnlineDevicesCount()
	if err != nil {
		log.Printf("Error getting online devices count: %v", err)
		return
	}

	// 创建 JSON 对象
	data := struct {
		OnlineDevices int `json:"online_devices"`
	}{
		OnlineDevices: count,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling online devices count: %v", err)
		return
	}

	// 发送 JSON 数据
	ForwardToCountClients(jsonData)
}

func FetchAndSendFleetData() {
	//fleets, err := db.FetchFleets() // 确保此方法实现正确，返回[]FleetData
	//if err != nil {
	//	log.Printf("Error fetching fleet data: %v", err)
	//	return
	//}

	//jsonData, err := json.Marshal(fleets)
	//if err != nil {
	//		log.Printf("Error marshaling fleet data: %v", err)
	//		return
	//}

	//ForwardToFleetClients(jsonData) // 假设您已经实现了此方法来发送数据
}
