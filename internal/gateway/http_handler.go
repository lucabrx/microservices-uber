package gateway

import (
	"encoding/json"
	"net/http"
	"strconv"

	pb_driver "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	pb_trip "github.com/lukabrx/uber-clone/api/proto/trip/v1"
)

type HttpHandler struct {
	driverClient pb_driver.DriverServiceClient
	tripClient   pb_trip.TripServiceClient
}

func NewHttpHandler(driverClient pb_driver.DriverServiceClient, tripClient pb_trip.TripServiceClient) *HttpHandler {
	return &HttpHandler{driverClient: driverClient, tripClient: tripClient}
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
	json.NewEncoder(w).Encode(res)
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
	json.NewEncoder(w).Encode(res)
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
	json.NewEncoder(w).Encode(res)
}
