package main

import (
	"RPS_SERVICE/internal/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 设置WebSocket和其他路由
	router.SetupRoutes()
	fmt.Println("WS_HTTP server started on :8080")
	// 启动服务器
	log.Println("HTTP server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
