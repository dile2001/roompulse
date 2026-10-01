package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorReading struct {
	ID          int64     `json:"id"`
	DeviceID    string    `json:"deviceId"`
	Temperature float64   `json:"temperature"`
	Timestamp   time.Time `json:"timestamp"`
	Humidity    float64   `json:"humidity"`
	ReceivedAt  time.Time `json:"receivedAt"`
}
type CreateReadingRequest struct {
	DeviceID    string     `json:"deviceId"`
	Temperature *float64   `json:"temperature"`
	Humidity    *float64   `json:"humidity"`
	Timestamp   *time.Time `json:"timestamp"`
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to create database pool:", err)
	}
	defer db.Close()
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	log.Println("Connected to PostgreSQL")
	http.HandleFunc("/readings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			limitStr := r.URL.Query().Get("limit")
			offsetStr := r.URL.Query().Get("offset")
			limit := 10
			offset := 0
			if limitStr != "" {
				value, err := strconv.Atoi(limitStr)
				if err == nil {
					if value > 100 {
						limit = 100
					} else if value >= 1 {
						limit = value
					}
				}
			}
			if offsetStr != "" {
				value, err := strconv.Atoi(offsetStr)
				if err == nil && value >= 0 {
					offset = value
				}
			}
			rows, err := db.Query(
				context.Background(),
				"SELECT id, device_id, temperature, humidity, measured_at, received_at FROM readings ORDER BY measured_at DESC LIMIT $1 OFFSET $2",
				limit, offset,
			)
			if err != nil {
				log.Println("Database query failed:", err)
				http.Error(w, "Failed to fetch readings", http.StatusInternalServerError)
				return
			}
			defer rows.Close()
			readings := []SensorReading{}
			for rows.Next() {
				var reading SensorReading

				err := rows.Scan(
					&reading.ID,
					&reading.DeviceID,
					&reading.Temperature,
					&reading.Humidity,
					&reading.Timestamp,
					&reading.ReceivedAt,
				)

				if err != nil {
					log.Println("Failed to scan reading:", err)
					continue
				}
				readings = append(readings, reading)
				fmt.Println(
					reading.ID,
					reading.DeviceID,
					reading.Temperature,
					reading.Humidity,
					reading.Timestamp,
					reading.ReceivedAt,
				)

			}
			if err := rows.Err(); err != nil {
				log.Println("Error occurred while iterating through rows:", err)
				http.Error(w, "Failed to fetch readings", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(readings)

			if err != nil {
				log.Println("Failed to encode readings:", err)
				return
			}
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var request CreateReadingRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&request)

		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if request.DeviceID == "" {
			http.Error(w, "deviceId is required", http.StatusBadRequest)
			return
		}
		if request.Temperature == nil {
			http.Error(w, "temperature is required", http.StatusBadRequest)
			return
		}

		if *request.Temperature < -50 || *request.Temperature > 100 {
			http.Error(w, "Temperature out of range", http.StatusBadRequest)
			return
		}
		if request.Humidity == nil {
			http.Error(w, "humidity is required", http.StatusBadRequest)
			return
		}
		humidity := *request.Humidity
		if humidity < 0 || humidity > 100 {
			http.Error(w, "humidity out of range", http.StatusBadRequest)
			return
		}
		if request.Timestamp == nil {
			http.Error(w, "timestamp is required", http.StatusBadRequest)
			return
		}

		timestamp := *request.Timestamp
		reading := SensorReading{
			DeviceID:    request.DeviceID,
			Temperature: *request.Temperature,
			Humidity:    *request.Humidity,
			Timestamp:   *request.Timestamp,
		}
		_, err = db.Exec(
			context.Background(),
			"INSERT INTO readings (device_id, temperature, humidity, measured_at) VALUES ($1, $2, $3, $4)",
			reading.DeviceID,
			reading.Temperature,
			reading.Humidity,
			reading.Timestamp,
		)
		if err != nil {
			log.Println("Database insert failed:", err)
			http.Error(w, "Failed to store reading", http.StatusInternalServerError)
			return
		}

		log.Println("Reading saved to PostgreSQL")
		log.Printf("Device: %s, Temperature: %.1f°C, Humidity: %.1f%%, Time: %s\n",
			request.DeviceID,
			*request.Temperature,
			*request.Humidity,
			timestamp.Format("15:04:05"),
		)
		w.Write([]byte("OK"))
	})
	fmt.Println("RoomPulse server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
