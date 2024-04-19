package handle

import (
	"RPS_SERVICE/internal/db"
	"encoding/json"
	"net/http"
)

func GetAllDevicesHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := db.GetAllDevices() // 假设这个函数已经存在于你的 db 包中，用于获取所有设备的ID
	if err != nil {
		http.Error(w, "Failed to get devices", http.StatusInternalServerError)
		return
	}

	// 设置响应类型为JSON
	w.Header().Set("Content-Type", "application/json")
	// 将设备列表编码为JSON并发送响应
	json.NewEncoder(w).Encode(devices)
}
