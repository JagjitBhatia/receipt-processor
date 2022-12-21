package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JagjitBhatia/receipt-processor/processor"
	"github.com/gorilla/mux"
)

// ProcessResponse represents the JSON response body for a successful process request
type ProcessResponse struct {
	ID string `json:"id"`
}

// PointsResponse represents the JSON response body for a successful points request
type PointsResponse struct {
	Points int `json:"points"`
}

// ErrorResponse represents the JSON response for an error message
type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	router := mux.NewRouter()
	rp := processor.NewReceiptProcessor()
	router.HandleFunc("/receipts/process", func(w http.ResponseWriter, r *http.Request) {
		var receipt processor.Receipt
		err := json.NewDecoder(r.Body).Decode(&receipt)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ErrorResponse{Error: "request body is not a valid receipt"})
			return
		}
		id, err := rp.ProcessReceipt(receipt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("receipt processing failed with error: %v", err.Error())})
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ProcessResponse{ID: id})
	}).Methods(http.MethodPost)

	router.HandleFunc("/receipts/{id}/points", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ErrorResponse{Error: "malformed request; no ID found"})
			return
		}

		points, err := rp.GetReceipt(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("points for receipt id %s could not be located", id)})
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PointsResponse{Points: points})
	}).Methods(http.MethodGet)

	http.ListenAndServe(":8080", router)
}
