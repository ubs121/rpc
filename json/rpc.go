// Copyright 2013 ubs121
package json

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	httpMethods string = "POST, GET, OPTIONS"
)

// GzipResponseWriter is used for gzipped response.
type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// ParseRequest parses a json request.
func ParseRequest(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&v)

	// DEBUG ONLY
	s, _ := json.Marshal(v)
	log.Printf("rpc << %s: %s\n", r.URL.String(), string(s))

	return err
}

// WriteResponse writes data into HTTP response.
func WriteResponse(r *http.Request, w http.ResponseWriter, result interface{}, e error) {
	// copy Access-Control-Request-Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", httpMethods)
	w.Header().Add("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	doc := map[string]interface{}{}
	if e != nil {
		doc["Error"] = e.Error()
	} else {
		doc["Result"] = result
	}

	// debug
	s, _ := json.Marshal(doc)
	log.Printf("rpc >> %s\n", string(s))

	encoder := json.NewEncoder(w)
	encoder.Encode(doc)

}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// makes gzip handler
func MakeGzipHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// gzip
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()

			handler.ServeHTTP(GzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
		} else {
			handler.ServeHTTP(w, r)
		}

	})
}
