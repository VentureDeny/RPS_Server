package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB 是全局数据库连接实例
var DB *sql.DB

// 初始化数据库连接
func init() {
	var err error
	DB, err = sql.Open("mysql", "root:123456@tcp(47.99.133.66:3306)/rps?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	// 测试数据库连接
	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("数据库连接成功！")
}

func SaveGPSData(deviceID string, location string) {
	stmt, err := DB.Prepare(`
		INSERT INTO gps_data(device_id, location)
		VALUES(?, ?)
		ON DUPLICATE KEY UPDATE location = VALUES(location)
	`)
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

// SaveRPSData 保存RPS数据到数据库
func SaveRPSData(deviceID string, x int, y int) {
	stmt, err := DB.Prepare(`
		INSERT INTO rps_data(device_id, x, y)
		VALUES(?, ?, ?)
		ON DUPLICATE KEY UPDATE x = VALUES(x), y = VALUES(y)
	`)
	if err != nil {
		log.Println("Prepare statement error:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(deviceID, x, y)
	if err != nil {
		log.Println("Execute statement error:", err)
		return
	}

	fmt.Println("RPS数据插入或更新成功！")
}

func SaveStatusData(deviceID string, battery string, MAC string) {
	stmt, err := DB.Prepare(`
		INSERT INTO status_data(device_id, battery_level, mac_address)
		VALUES(?, ?, ?)
		ON DUPLICATE KEY UPDATE battery_level = VALUES(battery_level), mac_address = VALUES(mac_address)
	`)
	if err != nil {
		log.Println("Prepare statement error:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(deviceID, battery, MAC)
	if err != nil {
		log.Println("Execute statement error:", err)
		return
	}

	fmt.Println("状态数据插入或更新成功！")
}

// AddDeviceToOnlineAndAll 添加设备到 onlinedevice 和 alldevice 数据库
func AddDeviceToOnlineAndAll(deviceID string) {
	// 添加到 alldevice
	stmtAll, err := DB.Prepare(`
        INSERT INTO alldevice(device_id)
        VALUES(?)
        ON DUPLICATE KEY UPDATE device_id = VALUES(device_id)
    `)
	if err != nil {
		log.Println("Prepare statement error (alldevice):", err)
		return
	}
	defer stmtAll.Close()

	_, err = stmtAll.Exec(deviceID)
	if err != nil {
		log.Println("Execute statement error (alldevice):", err)
		return
	}

	// 添加到 onlinedevice
	stmtOnline, err := DB.Prepare(`
        INSERT INTO onlinedevice(device_id)
        VALUES(?)
        ON DUPLICATE KEY UPDATE device_id = VALUES(device_id)
    `)
	if err != nil {
		log.Println("Prepare statement error (onlinedevice):", err)
		return
	}
	defer stmtOnline.Close()

	_, err = stmtOnline.Exec(deviceID)
	if err != nil {
		log.Println("Execute statement error (onlinedevice):", err)
		return
	}

	fmt.Println("设备已添加到 onlinedevice 和 alldevice！")
}

// RemoveDeviceFromOnline 从 onlinedevice 数据库中移除设备
func RemoveDeviceFromOnline(deviceID string) {
	stmt, err := DB.Prepare(`
        DELETE FROM onlinedevice WHERE device_id = ?
    `)
	if err != nil {
		log.Println("Prepare statement error:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(deviceID)
	if err != nil {
		log.Println("Execute statement error:", err)
		return
	}

	fmt.Println("设备已从 onlinedevice 移除！")
}

// GetAllDevices 获取 alldevice 表中的所有设备ID
func GetAllDevices() ([]string, error) {
	rows, err := DB.Query(`
		SELECT device_id FROM alldevice
	`)
	if err != nil {
		log.Println("Query alldevice error:", err)
		return nil, err
	}
	defer rows.Close()

	var devices []string
	for rows.Next() {
		var deviceID string
		if err := rows.Scan(&deviceID); err != nil {
			log.Println("Scan device_id error:", err)
			continue // 或返回错误
		}
		devices = append(devices, deviceID)
	}

	return devices, nil
}

// GetOnlineDevices 获取 onlinedevice 表中的所有在线设备ID
func GetOnlineDevices() ([]string, error) {
	rows, err := DB.Query(`
		SELECT device_id FROM onlinedevice
	`)
	if err != nil {
		log.Println("Query onlinedevice error:", err)
		return nil, err
	}
	defer rows.Close()

	var devices []string
	for rows.Next() {
		var deviceID string
		if err := rows.Scan(&deviceID); err != nil {
			log.Println("Scan device_id error:", err)
			continue // 或返回错误
		}
		devices = append(devices, deviceID)
	}

	return devices, nil
}

// GetGPSData 获取特定设备的最新 GPS 数据
// GetGPSData 获取特定设备的 GPS 数据
func GetGPSData(deviceID string) (string, error) {
	var location string
	err := DB.QueryRow(`
        SELECT location FROM gps_data WHERE device_id = ?
    `, deviceID).Scan(&location)

	if err != nil {
		log.Println("Query GPS data error:", err)
		return "", err
	}

	return location, nil
}

// GetStatusData 获取特定设备的最新状态数据
// GetStatusData 获取特定设备的状态数据
func GetStatusData(deviceID string) (string, string, error) {
	var batteryLevel, macAddress string
	err := DB.QueryRow(`
        SELECT battery_level, mac_address FROM status_data WHERE device_id = ?
    `, deviceID).Scan(&batteryLevel, &macAddress)

	if err != nil {
		log.Println("Query status data error:", err)
		return "", "", err
	}

	return batteryLevel, macAddress, nil
}

// GetOnlineDevicesCount 获取在线设备的数量
func GetOnlineDevicesCount() (int, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM onlinedevice`).Scan(&count)
	if err != nil {
		log.Printf("Error querying online devices count: %v", err)
		return 0, err
	}
	return count, nil
}
