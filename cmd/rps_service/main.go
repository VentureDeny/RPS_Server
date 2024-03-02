package main

import (
	"RPS_SERVICE/pkg/websocket"
	"log"
	"net/http"
)

func main() {
	// 设置WebSocket路由
	http.HandleFunc("/ws", websocket.HandleConnections)

	// 启动新的协程，用于监听并处理WebSocket消息
	go websocket.HandleMessages()

	// 启动服务器
	log.Println("HTTP server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
