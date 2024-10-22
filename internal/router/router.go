package router

import (
	"RPS_SERVICE/internal/handle"
	"log"
	"net/http"
)

// SetupRoutes 配置WebSocket路由
func SetupRoutes() {
	log.Println("Router Setup")
	http.HandleFunc("/data", handle.HandleDataWS)
	http.HandleFunc("/count", handle.HandleCountWS)
	http.HandleFunc("/common", handle.HandleCommonWS) // 使用相同的处理函数来处理/gps和/rps
	http.HandleFunc("/login", handle.HandleLoginWS)
	http.HandleFunc("/fleet", handle.HandleFleetWS)
	http.HandleFunc("/warning", handle.HandleWarningWS)
	http.HandleFunc("/api/devices", handle.GetAllDevicesHandler)
	log.Println("Fleet WebSocket router setup complete.")
}
