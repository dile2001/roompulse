package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SensorReading struct {
	DeviceID    string    `json:"deviceId"`
	Temperature float64   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}

func main() {
	http.HandleFunc("/readings", func(w http.ResponseWriter, r *http.Request) {
		var reading SensorReading
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&reading)

		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if reading.Temperature < -50 || reading.Temperature > 100 {
			http.Error(w, "Temperature out of range", http.StatusBadRequest)
			return
		}
		fmt.Printf("Device: %s, Temperature: %.1f°C, Time: %s\n",
			reading.DeviceID,
			reading.Temperature,
			reading.Timestamp.Format("15:04:05"),
		)
		w.Write([]byte("OK"))
	})
	fmt.Println("RoomPulse server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
