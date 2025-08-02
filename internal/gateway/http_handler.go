package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	pb_auth "github.com/lukabrx/uber-clone/api/proto/auth/v1"
	pb_driver "github.com/lukabrx/uber-clone/api/proto/driver/v1"
	pb_trip "github.com/lukabrx/uber-clone/api/proto/trip/v1"
	"github.com/lukabrx/uber-clone/internal/jsn"
	"golang.org/x/oauth2"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type HttpHandler struct {
	driverClient pb_driver.DriverServiceClient
	tripClient   pb_trip.TripServiceClient
	authClient   pb_auth.AuthServiceClient
	googleOauth  *oauth2.Config

	hub *Hub
}

func NewHttpHandler(
	driverClient pb_driver.DriverServiceClient,
	tripClient pb_trip.TripServiceClient,
	authClient pb_auth.AuthServiceClient,
	hub *Hub,
	googleOauth *oauth2.Config,
) *HttpHandler {

	return &HttpHandler{
		driverClient: driverClient,
		tripClient:   tripClient,
		authClient:   authClient,
		hub:          hub,
		googleOauth:  googleOauth,
	}
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

	h.hub.AddClient(conn)
	defer h.hub.RemoveClient(conn)

	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	res, err := h.driverClient.FindAvailableDrivers(
		r.Context(),
		&pb_driver.FindAvailableDriversRequest{Lat: lat, Lon: lon},
	)
	if err != nil {
		log.Println("initial find available drivers error:", err)
	} else {
		if err := conn.WriteJSON(res.Drivers); err != nil {
			log.Println("initial write json error:", err)
			return
		}
	}

	// Keep the connection open and listen for new messages
	for {
		msgType, msg, err := conn.NextReader()
		log.Println("receive message type:", msgType)
		log.Println("receive message:", msg)
		if err != nil {
			break
		}
	}
}

func (h *HttpHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	// The state parameter is a security measure to prevent CSRF attacks.
	url := h.googleOauth.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *HttpHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		jsn.ErrorJson(w, errors.New("google did not return a code"), http.StatusBadRequest)
		return
	}

	res, err := h.authClient.AuthenticateWithGoogle(r.Context(), &pb_auth.AuthenticateWithGoogleRequest{Code: code})
	if err != nil {
		log.Printf("Failed to authenticate with google: %v", err)
		jsn.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	redirectURL := frontendURL + "/auth/callback?token=" + res.Token
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

type contextKey string

const UserIDKey contextKey = "userID"

func (h *HttpHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsn.ErrorJson(w, errors.New("authorization header is required"), http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			jsn.ErrorJson(w, errors.New("authorization header format must be Bearer <token>"), http.StatusUnauthorized)
			return
		}

		token := parts[1]
		res, err := h.authClient.VerifyToken(r.Context(), &pb_auth.VerifyTokenRequest{Token: token})
		if err != nil {
			jsn.ErrorJson(w, errors.New("invalid token: "+err.Error()), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, res.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
