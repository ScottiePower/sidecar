package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ValidationRequest struct {
	Rego string `json:"rego"`
}

const (
	RequestMaxsize = 10 * (1 << 20) // 10MB
	Port           = "8071"
)

func main() {
	health := healthcheck.NewHandler()
	// Our app is not happy if we've got more than 100 goroutines running.
	health.AddLivenessCheck("Liveness", healthcheck.GoroutineCountCheck(100))
	health.AddReadinessCheck("Readiness", healthcheck.GoroutineCountCheck(100))

	// Define handlers for endpoints
	r := mux.NewRouter()
	r.HandleFunc("/live", health.LiveEndpoint)
	r.HandleFunc("/ready", health.ReadyEndpoint)
	r.HandleFunc("/validate", validateRequestJSON).Methods("POST")

	// Create http server
	http.Handle("/", r)
	port, err := strconv.Atoi(Port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listener started on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}

}

func validateRequestJSON(w http.ResponseWriter, r *http.Request) {
	log.Printf("HttpMethod: %s\t RequestURI: %s\n", r.Method, r.RequestURI)
	isValid := true
	errMsg := ""
	var validationRequest ValidationRequest
	err := json.NewDecoder(r.Body).Decode(&validationRequest)
	if err != nil || len(validationRequest.Rego) == 0 {
		isValid = false
		errMsg = "Invalid validation request."
		sendValidationResponse(w, isValid, errMsg)
		return
	}
	log.Printf("rego = %s\n", validationRequest.Rego)
	// Do nothing, just send valid response
	// Send back json response
	sendValidationResponse(w, isValid, errMsg)
}

func validateRequestTEXT(w http.ResponseWriter, r *http.Request) {
	log.Printf("HttpMethod: %s\t RequestURI: %s\n", r.Method, r.RequestURI)
	isValid := true
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
	log.Printf("rego = %s\n", string(data))
	// Do nothing, just send valid response
	// Send back json response
	sendValidationResponse(w, isValid, errMsg)
}

func sendValidationResponse(w http.ResponseWriter, isValid bool, errMsg string) {
	resp := map[string]interface{}{
		"valid":  isValid,
		"errors": errMsg,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error sending: %s", err)
	}
}

func CloseRequestBody(r *http.Request) {
	_ = r.Body.Close()
}
