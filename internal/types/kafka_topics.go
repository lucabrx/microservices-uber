package types

var (
	DriverLocationTopic = "driver_locations"
	TripEventsTopic     = "trip_events"
)

type EventType string

const (
	TripCreatedEvent   EventType = "TRIP_CREATED"
	TripCompletedEvent EventType = "TRIP_COMPLETED"
)

type TripEvent struct {
	EventType EventType `json:"event_type"`
	TripID    string    `json:"trip_id"`
	DriverID  string    `json:"driver_id"`
}
