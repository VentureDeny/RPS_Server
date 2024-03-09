package router

import (
	"log"
	"net/http"
)

// SetupRoutes 配置WebSocket路由
func SetupRoutes() {
	log.Println("Router Setup")
	http.HandleFunc("/data", handleDataWS)
	http.HandleFunc("/common", handleCommonWS) // 使用相同的处理函数来处理/gps和/rps
}
