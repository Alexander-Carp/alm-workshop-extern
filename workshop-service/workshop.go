package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

type Workshop struct {
	Name         string   `json:"name"`
	Date         string   `json:"date"`
	Presentator  string   `json:"presentator"`
	Participants []string `json:"participants"`
	SweaterScore int      `json:"sweaterScore"`
}

var defaultSweaterScore = readDefaultSweaterScore()

var workshop = Workshop{
	Name:         "ALM Workshop",
	Date:         "1/12/2025",
	Presentator:  "AE Consultants",
	Participants: []string{"John Doe", "Mary Little Lamb", "Chuck Norris", "Carp"},
	SweaterScore: defaultSweaterScore,
}

func readDefaultSweaterScore() int {
	// Read default from env var `SWEATER_SCORE_DEFAULT`. Fall back to 5 if unset/invalid.
	s := os.Getenv("SWEATER_SCORE_DEFAULT")
	if s == "" {
		return 5
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 1 || v > 10 {
		return 5
	}
	return v
}

func getWorkshopHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the struct to JSON and write it to the response
	json.NewEncoder(w).Encode(workshop)
}

func postWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON data into a new Workshop struct
	var newWorkshop Workshop
	err := json.NewDecoder(r.Body).Decode(&newWorkshop)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON data"))
		return
	}

	// If SweaterScore not provided (zero), set default. Otherwise validate range 1-10.
	if newWorkshop.SweaterScore == 0 {
		newWorkshop.SweaterScore = defaultSweaterScore
	} else if newWorkshop.SweaterScore < 1 || newWorkshop.SweaterScore > 10 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("SweaterScore must be between 1 and 10"))
		return
	}

	// Update the workshop details
	workshop = newWorkshop

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workshop)
}

func WorkshopHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getWorkshopHandler(w, r)
	case "POST":
		postWorkshopHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}
