package server

import (
	"encoding/json"
	"mini-crypto-exchange/internal/engine"
	"mini-crypto-exchange/internal/util"
	"net/http"
	"strings"

	"log"
)

type CreatePairRequest struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

type CreatePairResponse struct {
	Pair  string `json:"pair,omitempty"`
	Error string `json:"error,omitempty"`
}

// CreatePairHandler handles POST /api/pairs
func CreatePairHandler(engine *engine.MatchingEngine, config *util.RouterConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		var req CreatePairRequest
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			log.Printf("Failed to decode request: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CreatePairResponse{Error: "Invalid request body"})
			return
		}

		if strings.TrimSpace(req.Base) == "" || strings.TrimSpace(req.Quote) == "" {
			log.Println("Missing base or quote")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CreatePairResponse{Error: "Base and quote are required"})
			return
		}

		pair := req.Base + "/" + req.Quote
		engine.CreatePair(pair)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CreatePairResponse{Pair: pair})
	}
}
