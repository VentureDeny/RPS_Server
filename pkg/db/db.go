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
	DB, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/rps?parseTime=true")
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
