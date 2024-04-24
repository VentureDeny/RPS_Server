package db

import (
	datastruct "RPS_SERVICE/internal/struct"
	"fmt"
	"log"
)

func SaveStatusData(deviceID string, battery string, MAC string, location string, speed string, accelerationX string, accelerationY string, accelerationZ string) {
	stmt, err := DB.Prepare(`INSERT INTO status_data(device_id, location, battery_level, mac_address, speed, accelerationX, accelerationY, accelerationZ)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE battery_level = VALUES(battery_level), mac_address = VALUES(mac_address), location = VALUES(location), speed = VALUES(speed), accelerationX = VALUES(accelerationX), accelerationY = VALUES(accelerationY), accelerationZ = VALUES(accelerationZ)`)
	if err != nil {
		log.Println("Prepare statement error:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(deviceID, location, battery, MAC, speed, accelerationX, accelerationY, accelerationZ)
	if err != nil {
		log.Println("Execute statement error:", err)
		return
	}

	fmt.Println("状态数据插入或更新成功！")
}

func QueryDeviceData() ([]datastruct.DeviceData, error) {
	// 查询数据
	rows, err := DB.Query("SELECT device_id, location, battery_level, mac_address, speed, accelerationX, accelerationY, accelerationZ FROM status_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 将查询结果转换为结构体切片
	var devices []datastruct.DeviceData
	for rows.Next() {
		var device datastruct.DeviceData
		err := rows.Scan(&device.DeviceID, &device.Location, &device.BatteryLevel, &device.MacAddress, &device.Speed, &device.AccelerationX, &device.AccelerationY, &device.AccelerationZ)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}
