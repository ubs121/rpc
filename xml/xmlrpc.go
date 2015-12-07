// Copyright 2013 ubs121

package xml

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

// TODO: xml-rpc request type

// ReadRequest parses the request object.
func ParseRequest(r *http.Request, v interface{}) error {
	decoder := xml.NewDecoder(r.Body)
	err := decoder.Decode(&v)

	// DEBUG
	log.Printf("xmlrpc << %s: %v\n", r.URL.String(), v)

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

	log.Printf("xmlrpc >> %s\n", xmlstr)

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write([]byte(xmlstr))
}
