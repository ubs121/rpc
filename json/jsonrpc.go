// Copyright 2013 ubs121
/*
	Custom  json-rpc over http implementation.
	Idea is to avoid using reflect on methods (like net/rpc) for faster serialization.
	Standard encoding/json package is used for JSON marshalling/unmarshalling.
*/

package json

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	httpMethods string = "POST, GET, OPTIONS"
)

// ParseRequest parses a json request.
func ParseRequest(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&v)

	// DEBUG
	log.Printf("jsonrpc << %s: %v\n", r.URL.String(), v)

	return err
}

// WriteResponse converts data into json format and writes into response.
func WriteResponse(r *http.Request, w http.ResponseWriter, result interface{}, e error) {
	doc := map[string]interface{}{}
	if e != nil {
		doc["Error"] = e.Error()
	} else {
		doc["Result"] = result
	}

	// DEBUG
	log.Printf("jsonrpc >> %v\n", doc)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	setCORSHeaders(r, w)

	encoder := json.NewEncoder(w)
	encoder.Encode(doc)

}

// set CORS headers
func setCORSHeaders(r *http.Request, w http.ResponseWriter) {
	// TODO: copy Access-Control-Request-Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", httpMethods)
	w.Header().Add("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
}
