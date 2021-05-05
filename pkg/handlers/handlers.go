package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-redisapi/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-redisapi/pkg/schema"
	"github.com/microlib/simple"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

func GraphDataHandler(w http.ResponseWriter, r *http.Request, logger *simple.Logger, con connectors.Clients) {
	var response *schema.Response
	var data []byte

	id := r.URL.Query().Get("id")
	graph := r.URL.Query().Get("graph")
	addHeaders(w, r)
	if id == "" {
		id = "default"
	}
	if graph == "" {
		id = "line"
	}

	if r.Method == http.MethodGet {
		data, _ = con.Get(id + "-i" + graph)
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "KO", Message: fmt.Sprintf("Could not read body data %v\n", err), Payload: ""}
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			con.Set(id+"-"+graph, body, 0)
			data = body
		}
	}
	response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data processed succesfully", Payload: string(data)}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("LineDataHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func ApiCallHandler(w http.ResponseWriter, r *http.Request, logger *simple.Logger, con connectors.Clients) {
	var response *schema.Response

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "KO", Message: fmt.Sprintf("Could not read body data %v\n", err), Payload: ""}
		w.WriteHeader(http.StatusInternalServerError)
	}
	req, err := http.NewRequest("POST", os.Getenv("URL"), bytes.NewBuffer(body))
	resp, err := con.Do(req)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		logger.Error(fmt.Sprintf("Http request %v\n", err))
		response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "KO", Message: fmt.Sprintf("Http request error %v\n", err), Payload: ""}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "KO", Message: fmt.Sprintf("Could not read body data %v\n", err), Payload: ""}
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "ApiCall request processed succesfully", Payload: string(body)}
			w.WriteHeader(http.StatusOK)
		}
	}
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("ApiCallHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \"1.0.2\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
