package datastruct

type DeviceData struct {
	DeviceID      string `json:"device_id"`
	Location      string `json:"location,omitempty"`
	BatteryLevel  int    `json:"battery_level,omitempty"`
	MacAddress    string `json:"mac_address,omitempty"`
	Speed         string `json:"speed,omitempty"`
	AccelerationX string `json:"accelerationX,omitempty"`
	AccelerationY string `json:"accelerationY,omitempty"`
	AccelerationZ string `json:"accelerationZ,omitempty"`
}
type FleetData struct {
	ID       string   `json:"id,omitempty"`
	Name     string   `json:"name"`
	Vehicles []string `json:"vehicles"`
}
type WarningData struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
}
