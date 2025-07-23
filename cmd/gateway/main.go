package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lukabrx/uber-clone/internal/gateway"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/drivers", gateway.RegisterDriver).Methods("POST")
	r.HandleFunc("/drivers/{id}", gateway.UnregisterDriver).Methods("DELETE")
	r.HandleFunc("/drivers/{id}/availability", gateway.CheckDriverAvailability).Methods("GET")

	r.HandleFunc("/trips", gateway.BookTrip).Methods("POST")

	log.Println("Gateway server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
