package main

import (
    "fmt"
    "net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
)

const (
	Port               = 8080
	RequestMaxsize     = 10 * (1 << 20) // 10MB
)

func main() {
	// Define handlers for endpoints
	r := mux.NewRouter()
	r.HandleFunc("/validate", validateRequest).Methods("POST")
	// Create http server
	http.Handle("/", r)
	log.Printf("Listener started on port %d", Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), nil); err != nil {
		log.Fatal(err)
	}
}

func CloseRequestBody(r *http.Request) {
	_ = r.Body.Close()
}

func validateRequest(w http.ResponseWriter, r *http.Request) {
	validRequest := true
	errMsg := ""

	// Security guide, prevent users from send huge files by limiting reader
	defer CloseRequestBody(r)
	rdr := io.LimitReader(r.Body, RequestMaxsize)
	data, err := ioutil.ReadAll(rdr)
	if err != nil {
		log.Printf("read error: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // Send back post data

	// Send back json response
	resp := map[string]interface{}{
		"validRequest":  validRequest,
		"errors": errMsg,
		"data": string(data),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error sending: %s", err)
	}
}