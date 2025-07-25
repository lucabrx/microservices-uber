package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	pb_driver "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	pb_trip "github.com/lukabrx/uber-clone/api/proto/trip/v1"
	"github.com/lukabrx/uber-clone/internal/jsn"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type HttpHandler struct {
	driverClient pb_driver.DriverServiceClient
	tripClient   pb_trip.TripServiceClient
}

func NewHttpHandler(driverClient pb_driver.DriverServiceClient, tripClient pb_trip.TripServiceClient) *HttpHandler {
	return &HttpHandler{driverClient, tripClient}
}

func (h *HttpHandler) RegisterDriver(w http.ResponseWriter, r *http.Request) {
	var req pb_driver.RegisterDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.driverClient.RegisterDriver(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsn.WriteJson(w, http.StatusCreated, res.Driver)
}

func (h *HttpHandler) FindAvailableDrivers(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)

	req := &pb_driver.FindAvailableDriversRequest{Lat: lat, Lon: lon}
	res, err := h.driverClient.FindAvailableDrivers(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsn.WriteJson(w, http.StatusOK, res.Drivers)
}

func (h *HttpHandler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	var req pb_trip.CreateTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.tripClient.CreateTrip(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsn.WriteJson(w, http.StatusCreated, res.Trip)
}

func (h *HttpHandler) CompleteTrip(w http.ResponseWriter, r *http.Request) {
	tripID := chi.URLParam(r, "id")
	if tripID == "" {
		http.Error(w, "trip_id is required in the URL path", http.StatusBadRequest)
		return
	}

	req := &pb_trip.CompleteTripRequest{TripId: tripID}
	res, err := h.tripClient.CompleteTrip(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsn.WriteJson(w, http.StatusOK, res.Trip)
}

func (h *HttpHandler) StreamAvailableDrivers(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()

	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)

	for {
		req := &pb_driver.FindAvailableDriversRequest{Lat: lat, Lon: lon}
		res, err := h.driverClient.FindAvailableDrivers(r.Context(), req)
		if err != nil {
			log.Println("find available drivers error:", err)
			if r.Context().Err() != nil {
				break
			}
			continue
		}

		if err := conn.WriteJSON(res.Drivers); err != nil {
			log.Println("write json error:", err)
			break
		}

		time.Sleep(3 * time.Second)
	}
}
