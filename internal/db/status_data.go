package db

import (
	"log"
)

// GetStatusData 获取特定设备的最新 GPS 数据
func GetStatusData(deviceID string) (string, string, error) {
	var batteryLevel, macAddress string
	err := DB.QueryRow(`SELECT battery_level, mac_address FROM status_data WHERE device_id = ?`, deviceID).Scan(&batteryLevel, &macAddress)
	if err != nil {
		log.Printf("GetStatusData: %v", err)
		return "", "", err
	}

	return batteryLevel, macAddress, nil
}
