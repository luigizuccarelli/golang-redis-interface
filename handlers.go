package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

// Response schema
type Response struct {
	Name       string `json:"name"`
	StatusCode string `json:"statuscode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Counter    uint64 `json:"counter"`
	Payload    string `json:"payload"`
}

func SimpleHandler(w http.ResponseWriter, r *http.Request) {
	var response Response

	passFail := r.URL.Query().Get("flag")
	addHeaders(w, r)
	handleOptions(w, r)

	counter++

	if passFail == "" || passFail == "fail" {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: "The system detected an error (simulated)"}
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			response = handleError(w, "Could not read body data "+err.Error())
		} else {
			response = Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data uploaded succesfully", Counter: counter, Payload: string(body)}
		}
	}
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("SimpleHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \"1.0.2\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	var request []string
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	logger.Trace(fmt.Sprintf("Headers : %s", request))

	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}

// simple options handler
func handleOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
	}
	return
}

// simple error handler
func handleError(w http.ResponseWriter, msg string) Response {
	w.WriteHeader(http.StatusInternalServerError)
	r := Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: msg}
	return r
}
