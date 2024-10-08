package db

import (
	datastruct "RPS_SERVICE/internal/struct"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB 是全局数据库连接实例
var DB *sql.DB

// 初始化数据库连接
func init() {
	var err error
	DB, err = sql.Open("mysql", "root:1850560Dwc@tcp(localhost:3306)/rps?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	// 测试数据库连接
	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("数据库连接成功！")
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

func CreateFleet(name string, vehicles []string) error {
	vehiclesJSON, err := json.Marshal(vehicles)
	if err != nil {
		return err
	}

	stmt, err := DB.Prepare(`INSERT INTO fleet (name, vehicles) VALUES (?, ?)`)
	if err != nil {
		log.Println("Prepare statement error:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, vehiclesJSON)
	if err != nil {
		log.Println("Execute statement error:", err)
		return err
	}

	fmt.Println("车队创建成功！")
	return nil
}

// UpdateFleet 更新数据库中的现有车队
func UpdateFleet(fleetID, name string, vehicles []string) error {
	vehiclesJSON, err := json.Marshal(vehicles)
	if err != nil {
		return err
	}

	stmt, err := DB.Prepare(`UPDATE fleet SET name = ?, vehicles = ? WHERE fleet_id = ?`)
	if err != nil {
		log.Println("Prepare statement error:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, vehiclesJSON, fleetID)
	if err != nil {
		log.Println("Execute statement error:", err)
		return err
	}

	fmt.Println("车队更新成功！")
	return nil
}

// DeleteFleet 从数据库中删除一个车队
func DeleteFleet(fleetID string) error {
	stmt, err := DB.Prepare(`DELETE FROM fleet WHERE fleet_id = ?`)
	if err != nil {
		log.Println("Prepare statement error:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fleetID)
	if err != nil {
		log.Println("Execute statement error:", err)
		return err
	}

	fmt.Println("车队删除成功！")
	return nil
}
func GetFleets() ([]datastruct.FleetData, error) {
	var fleets []datastruct.FleetData

	// 编写 SQL 查询语句
	query := `SELECT fleet_id, name, vehicles FROM fleet`
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error querying fleet data: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f datastruct.FleetData
		var vehiclesJSON string // 用字符串接收JSON数据

		err := rows.Scan(&f.ID, &f.Name, &vehiclesJSON)
		if err != nil {
			log.Printf("Error scanning fleet data: %v", err)
			continue // 也可以选择返回错误
		}

		// 将 JSON 字符串转换为 string 切片
		err = json.Unmarshal([]byte(vehiclesJSON), &f.Vehicles)
		if err != nil {
			log.Printf("Error unmarshalling vehicles JSON: %v", err)
			continue // 同上，根据需要处理错误
		}

		fleets = append(fleets, f)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
		return nil, err
	}

	return fleets, nil
}
func SaveWarning(data datastruct.WarningData) error {
	stmt, err := DB.Prepare("INSERT INTO warnings (id, type, message, level, timestamp) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.ID, data.Type, data.Message, data.Level, data.Timestamp)
	return err
}
func FetchAllWarnings() ([]datastruct.WarningData, error) {
	var warnings []datastruct.WarningData
	rows, err := DB.Query("SELECT id, type, message, level, timestamp FROM warnings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var wd datastruct.WarningData
		if err := rows.Scan(&wd.ID, &wd.Type, &wd.Message, &wd.Level, &wd.Timestamp); err != nil {
			log.Println("Failed to scan warning:", err)
			continue
		}
		warnings = append(warnings, wd)
	}

	return warnings, nil
}

func DeleteWarningByTimestamp(timestamp string) error {
	_, err := DB.Exec("DELETE FROM warnings WHERE timestamp = ?", timestamp)
	return err
}
