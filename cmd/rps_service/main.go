package main

import (
	"RPS_SERVICE/pkg/router"
	"RPS_SERVICE/pkg/websocket"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 设置WebSocket和其他路由
	router.SetupRoutes()

	// 启动新的协程，用于监听并处理WebSocket消息
	go websocket.HandleMessages()
	fmt.Println("WS_HTTP server started on :8080")
	// 启动服务器
	log.Println("HTTP server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
