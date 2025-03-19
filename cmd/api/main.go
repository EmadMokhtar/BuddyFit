package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/EmadMokhtar/BuddyFit/internal"
	"github.com/rs/cors"
)

type BuddyFitRequest struct {
	Prompt string `json:"prompt"`
}

type BuddyFitResponse struct {
	Response string `json:"response"`
}

func askAIHandler(w http.ResponseWriter, r *http.Request) {
	var req BuddyFitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	model := "llama3.1:latest"

	response := internal.AskAI(req.Prompt, model)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	for msg := range response {
		res := BuddyFitResponse{Response: msg}
		if err := encoder.Encode(res); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		w.(http.Flusher).Flush()
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ask", askAIHandler)

	// Configure the API server
	portAPI := os.Getenv("BF_API_PORT")
	if portAPI == "" {
		portAPI = "8000"
	}

	hostAPI := os.Getenv("BF_API_HOST")
	if hostAPI == "" {
		hostAPI = "localhost"
	}

	portUI := os.Getenv("BF_UI_PORT")
	if portUI == "" {
		portUI = "3000"
	}

	hostUI := os.Getenv("BF_UI_HOST")
	if hostUI == "" {
		hostUI = "localhost"
	}

	addrAPI := fmt.Sprintf("%s:%s", hostAPI, portAPI)
	addrUI := fmt.Sprintf("http://%s:%s", hostUI, portUI)
	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{addrAPI, addrUI},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	log.Printf("Starting server on %s\n", addrAPI)
	if err := http.ListenAndServe(addrAPI, handler); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}
