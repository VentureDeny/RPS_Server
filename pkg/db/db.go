package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBConn 是一个全局数据库连接实例
var DBConn *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error
	// 请根据你的数据库配置调整DSN（数据源名称）
	dsn := "username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Database connection successfully established")
}

// CloseDB 关闭数据库连接
func CloseDB() {
	db, err := DBConn.DB()
	if err != nil {
		log.Fatalf("Error when closing the database connection: %v", err)
	}
	db.Close()
}
