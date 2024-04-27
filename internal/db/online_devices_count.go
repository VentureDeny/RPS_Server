package db

import (
	"log"
)

// GetOnlineDevicesCount 获取在线设备的数量
func GetOnlineDevicesCount() (int, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM onlinedevice`).Scan(&count)
	if err != nil {
		log.Printf("GetOnlineDevicesCount: %v", err)
		return 0, err
	}
	return count, nil
}
