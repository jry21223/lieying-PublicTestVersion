package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ptr(t time.Time) *time.Time {
	return &t
}
