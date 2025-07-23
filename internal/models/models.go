package models

import "time"

type Driver struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	IsAvailable bool    `json:"is_available"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

type Trip struct {
	ID          string    `json:"id"`
	RiderID     string    `json:"rider_id"`
	DriverID    string    `json:"driver_id,omitempty"`
	StartLat    float64   `json:"start_lat"`
	StartLon    float64   `json:"start_lon"`
	EndLat      float64   `json:"end_lat"`
	EndLon      float64   `json:"end_lon"`
	Status      string    `json:"status"` // e.g., "requested", "in_progress", "completed"
	Price       float64   `json:"price,omitempty"`
	RequestTime time.Time `json:"request_time"`
}
