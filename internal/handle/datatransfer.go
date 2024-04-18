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
	// 假设我们已经有了一组设备ID，可以是通过 GetAllDevices 函数获取
	deviceIDs, err := db.GetAllDevices()
	if err != nil {
		log.Println("Error fetching device IDs:", err)
		return
	}

	// 为了演示，我们创建一个 slice 来存储所有的设备信息
	var devicesInfo []DeviceInfo

	for _, deviceID := range deviceIDs {
		var deviceInfo DeviceInfo
		deviceInfo.DeviceID = deviceID

		// 获取设备的 GPS 数据
		location, err := db.GetGPSData(deviceID)
		if err == nil {
			deviceInfo.Location = location
		}

		// 获取设备的状态数据
		batteryLevel, macAddress, err := db.GetStatusData(deviceID)
		if err == nil {
			deviceInfo.BatteryLevel = batteryLevel
			deviceInfo.MACAddress = macAddress
		}

		devicesInfo = append(devicesInfo, deviceInfo)
	}

	// 将设备信息序列化为 JSON
	jsonData, err := json.Marshal(devicesInfo)
	if err != nil {
		log.Println("Error marshalling devices info:", err)
		return
	}

	// 发送 JSON 数据给所有连接到 /data 的客户端
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
