package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type SensorReading struct {
	DeviceID    string    `json:"deviceId"`
	Temperature float64   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
}

func main() {
	reading := SensorReading{
		DeviceID:    "room-01",
		Temperature: 27.4,
		Timestamp:   time.Now(),
	}
	for {

		fmt.Printf("Device: %s\n", reading.DeviceID)
		fmt.Printf("Temperature: %.1f°C\n", reading.Temperature)
		reading.Temperature += (rand.Float64() - 0.5) * 0.6
		data, err := json.Marshal(reading)

		if err != nil {
			fmt.Println("JSON error:", err)
			continue
		}
		resp, err := http.Post(
			"http://localhost:8080/readings",
			"application/json",
			bytes.NewBuffer(data),
		)

		if err != nil {
			fmt.Println("JSON error:", err)
			continue
		}

		resp.Body.Close()
		time.Sleep(2 * time.Second)
	}
}
