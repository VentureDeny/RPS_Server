package db

import (
	"fmt"
	"log"
)

func SaveGPSData(deviceID string, location string) {
	stmt, err := DB.Prepare(`INSERT INTO gps_data(device_id, location)
	VALUES(?, ?)
	ON DUPLICATE KEY UPDATE location = VALUES(location)`)
	if err != nil {
		log.Println("Prepare statement error:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(deviceID, location)
	if err != nil {
		log.Println("Execute statement error:", err)
		return
	}

	fmt.Println("GPS数据插入或更新成功！")
}

func GetGPSData(deviceID string) (string, error) {
	var location string
	err := DB.QueryRow(`SELECT location FROM gps_data WHERE device_id = ?`, deviceID).Scan(&location)

	if err != nil {
		//log.Println("Query GPS data error:", err)
		return "", err
	}

	return location, nil
}
