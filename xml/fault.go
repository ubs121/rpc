// Copyright 2013 ubs121

package xml

import (
	"fmt"
)

// Default Faults
// NOTE: XMLRPC spec doesn't specify any Fault codes.
// These codes seems to be widely accepted, and taken from the http://xmlrpc-epi.sourceforge.net/specs/rfc.fault_codes.php
var (
	FaultInvalidParams        = Fault{Code: -32602, String: "Invalid Method Parameters"}
	FaultWrongArgumentsNumber = Fault{Code: -32602, String: "Wrong Arguments Number"}
	FaultInternalError        = Fault{Code: -32603, String: "Internal Server Error"}
	FaultApplicationError     = Fault{Code: -32500, String: "Application Error"}
	FaultSystemError          = Fault{Code: -32400, String: "System Error"}
	FaultDecode               = Fault{Code: -32700, String: "Parsing error: not well formed"}
)

// Fault represents XML-RPC Fault.
type Fault struct {
	Code   int    `xml:"faultCode"`
	String string `xml:"faultString"`
}

// Error satisifies error interface for Fault.
func (f Fault) Error() string {
	return fmt.Sprintf("%d: %s", f.Code, f.String)
}

// Fault2XML is a quick 'marshalling' replacement for the Fault case.
func fault2XML(fault Fault) string {
	buffer := "<methodResponse><fault>"
	xml, _ := rpc2XML(fault)
	buffer += xml
	buffer += "</fault></methodResponse>"
	return buffer
}

type faultValue struct {
	Value value `xml:"value"`
}

// IsEmpty returns true if faultValue contain fault.
//
// faultValue should be a struct with 2 members.
func (f faultValue) IsEmpty() bool {
	return len(f.Value.Struct) == 0
}
