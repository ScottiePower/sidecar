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

type ValidationRequest struct {
	Rego string `json:"rego"`
}

const (
	RequestMaxsize     = 10 * (1 << 20) // 10MB
	Port               = 8071
)

func main() {
	// Define handlers for endpoints
	r := mux.NewRouter()
	r.HandleFunc("/policy/v0/policies/rego/validate", validateRequestJSON).Methods("POST")

	// Create http server
	http.Handle("/", r)
	log.Printf("Listener started on port %d", Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), nil); err != nil {
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

