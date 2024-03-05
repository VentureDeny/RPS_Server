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
	DB, err = sql.Open("mysql", "用户名:密码@tcp(数据库地址:端口)/数据库名?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	// 测试数据库连接
	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("数据库连接成功！")
}

// SaveGPSData 保存GPS数据到数据库
func SaveGPSData(gpsID string, location string) {
	stmt, err := DB.Prepare("INSERT INTO gps_data(gps_id, location) VALUES(?, ?)")
	if err != nil {
		log.Println("Prepare statement error:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(gpsID, location)
	if err != nil {
		log.Println("Execute statement error:", err)
		return
	}

	fmt.Println("GPS数据插入成功！")
}

// SaveRPSData 保存RPS数据到数据库
func SaveRPSData(rpsID string, x int, y int) {
	stmt, err := DB.Prepare("INSERT INTO rps_data(rps_id, x, y) VALUES(?, ?, ?)")
	if err != nil {
		log.Println("Prepare statement error:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(rpsID, x, y)
	if err != nil {
		log.Println("Execute statement error:", err)
		return
	}

	fmt.Println("RPS数据插入成功！")
}
