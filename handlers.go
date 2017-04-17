package main 
import (
    "net/http"
	"fmt"
	"encoding/json"
)

type Response struct { 
	Name string `json:"name"`
	Email string `json:"email"`
	Status string `json:"status"`
}

func SimpleHandler(w http.ResponseWriter, r *http.Request) {
	var response Response
    response = Response{Status: "200", Name: "Napoleon", Email: "napoleon@dynamite.co"}

    b, _ := json.Marshal(response)
    fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ok version 1.0")
}