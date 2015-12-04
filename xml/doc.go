/*
In current go's rpc implementation doesn't handle interface{} types. So, marshalling is implemented manually.

This package converts response data directly to the string XML representation using standard encoding/xml package.

Types for marshalling

The following types are supported:

    XML-RPC             Golang
    -------             ------
    int, i4             int
    double              float64
    boolean             bool
    stringi             string
    dateTime.iso8601    time.Time
    base64              []byte
    struct              struct
    array               []interface{}
    nil                 nil
*/
package xml
