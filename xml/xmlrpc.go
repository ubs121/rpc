// Copyright 2013 ubs121

package xml

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// ReadRequest parses the request object.
func ParseRequest(r *http.Request, v interface{}) error {
	decoder := xml.NewDecoder(r.Body)
	err := decoder.Decode(&v)

	// DEBUG ONLY
	// log.Printf("xmlrpc << %s: %s\n", r.URL.String(), string(s))

	return err
}

// WriteResponse encodes the response and writes it to the ResponseWriter.
func WriteResponse(w http.ResponseWriter, result interface{}, err error) {
	var xmlstr string
	if err != nil {
		var fault Fault
		switch err.(type) {
		case Fault:
			fault = err.(Fault)
		default:
			fault = FaultApplicationError
			fault.String += fmt.Sprintf(": %v", err)
		}
		xmlstr = fault2XML(fault)
	} else {
		xmlstr, _ = rpcResponse2XML(result)
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write([]byte(xmlstr))
}
