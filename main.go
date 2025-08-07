package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"personalization-content-converter/utils"
	"time"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TranslationResponse struct {
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	
	response := HealthResponse{Status: "ok"}
	json.NewEncoder(w).Encode(response)
	
	slog.Info("Health check request",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
		"status", 200,
		"duration_ms", time.Since(start).Milliseconds(),
	)
}

func uoToCommonHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")

	var uoRequest utils.UOCurrentRequestFormat
	if err := json.NewDecoder(r.Body).Decode(&uoRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		
		slog.Error("UO to Common translation failed - invalid JSON",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"status", 400,
			"error", err.Error(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return
	}

	translator := &utils.UOToCommonTranslator{}
	commonRequest, err := translator.Translate(&uoRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		
		slog.Error("UO to Common translation failed - translator error",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"status", 500,
			"error", err.Error(),
			"user_id", uoRequest.IsEvent.User.ID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return
	}

	response := TranslationResponse{
		Request:  uoRequest,
		Response: commonRequest,
	}
	json.NewEncoder(w).Encode(response)
	
	slog.Info("UO to Common translation successful",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
		"status", 200,
		"user_id", uoRequest.IsEvent.User.ID,
		"action", uoRequest.IsEvent.Action,
		"duration_ms", time.Since(start).Milliseconds(),
	)
}

func commonToUoHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")

	var commonRequest utils.CommonRequestFormat
	if err := json.NewDecoder(r.Body).Decode(&commonRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		
		slog.Error("Common to UO translation failed - invalid JSON",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"status", 400,
			"error", err.Error(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return
	}

	translator := &utils.CommonToUOTranslator{}
	uoRequest, err := translator.Translate(&commonRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		
		slog.Error("Common to UO translation failed - translator error",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"status", 500,
			"error", err.Error(),
			"user_id", commonRequest.User.ID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return
	}

	response := TranslationResponse{
		Request:  commonRequest,
		Response: uoRequest,
	}
	json.NewEncoder(w).Encode(response)
	
	slog.Info("Common to UO translation successful",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
		"status", 200,
		"user_id", commonRequest.User.ID,
		"event_type", commonRequest.Event.Type,
		"duration_ms", time.Since(start).Milliseconds(),
	)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	
	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /translate/uo-to-common", uoToCommonHandler)
	mux.HandleFunc("POST /translate/common-to-uo", commonToUoHandler)
	
	slog.Info("Server starting", "port", 8080)
	
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("Server failed to start", "error", err.Error())
		os.Exit(1)
	}
}