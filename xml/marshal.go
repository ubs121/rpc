// Copyright 2013 ubs121

package xml

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type response struct {
	Name   xml.Name   `xml:"methodResponse"`
	Params []param    `xml:"params>param"`
	Fault  faultValue `xml:"fault,omitempty"`
}

type param struct {
	Value value `xml:"value"`
}

type value struct {
	Array    []value  `xml:"array>data>value"`
	Struct   []member `xml:"struct>member"`
	String   string   `xml:"string"`
	Int      string   `xml:"int"`
	Int4     string   `xml:"i4"`
	Double   string   `xml:"double"`
	Boolean  string   `xml:"boolean"`
	DateTime string   `xml:"dateTime.iso8601"`
	Base64   string   `xml:"base64"`
	Raw      string   `xml:",innerxml"` // the value can be defualt string
}

type member struct {
	Name  string `xml:"name"`
	Value value  `xml:"value"`
}

func rpcRequest2XML(method string, v interface{}) (string, error) {
	buffer := "<methodCall><methodName>"
	buffer += method
	buffer += "</methodName>"
	params, err := rpcParams2XML(v)
	buffer += params
	buffer += "</methodCall>"
	return buffer, err
}

func rpcResponse2XML(v interface{}) (string, error) {
	buffer := "<methodResponse>"
	params, err := rpcParams2XML(v)
	buffer += params
	buffer += "</methodResponse>"
	return buffer, err
}

func rpcParams2XML(v interface{}) (string, error) {
	var err error
	buffer := "<params>"
	for i := 0; i < reflect.ValueOf(v).Elem().NumField(); i++ {
		var xml string
		buffer += "<param>"
		xml, err = rpc2XML(reflect.ValueOf(v).Elem().Field(i).Interface())
		buffer += xml
		buffer += "</param>"
	}
	buffer += "</params>"
	return buffer, err
}

func rpc2XML(v interface{}) (string, error) {
	out := "<value>"
	switch reflect.ValueOf(v).Kind() {
	case reflect.Int:
		out += fmt.Sprintf("<int>%d</int>", v.(int))
	case reflect.Float64:
		out += fmt.Sprintf("<double>%f</double>", v.(float64))
	case reflect.String:
		out += string2XML(v.(string))
	case reflect.Bool:
		out += bool2XML(v.(bool))
	case reflect.Struct:
		if reflect.TypeOf(v).String() != "time.Time" {
			out += struct2XML(v)
		} else {
			out += time2XML(v.(time.Time))
		}
	case reflect.Slice, reflect.Array:
		// FIXME: is it the best way to recognize '[]byte'?
		if reflect.TypeOf(v).String() != "[]uint8" {
			out += array2XML(v)
		} else {
			out += base642XML(v.([]byte))
		}
	case reflect.Ptr:
		if reflect.ValueOf(v).IsNil() {
			out += "<nil/>"
		}
	}
	out += "</value>"
	return out, nil
}

func bool2XML(v bool) string {
	var b string
	if v {
		b = "1"
	} else {
		b = "0"
	}
	return fmt.Sprintf("<boolean>%s</boolean>", b)
}

func string2XML(v string) string {
	v = strings.Replace(v, "&", "&amp;", -1)
	v = strings.Replace(v, "\"", "&quot;", -1)
	v = strings.Replace(v, "<", "&lt;", -1)
	v = strings.Replace(v, ">", "&gt;", -1)
	return fmt.Sprintf("<string>%s</string>", v)
}

func struct2XML(v interface{}) (out string) {
	out += "<struct>"
	for i := 0; i < reflect.TypeOf(v).NumField(); i++ {
		field := reflect.ValueOf(v).Field(i)
		field_type := reflect.TypeOf(v).Field(i)
		var name string
		if field_type.Tag.Get("xml") != "" {
			name = field_type.Tag.Get("xml")
		} else {
			name = field_type.Name
		}
		field_value, _ := rpc2XML(field.Interface())
		field_name := fmt.Sprintf("<name>%s</name>", name)
		out += fmt.Sprintf("<member>%s%s</member>", field_name, field_value)
	}
	out += "</struct>"
	return
}

func array2XML(v interface{}) (out string) {
	out += "<array><data>"
	for i := 0; i < reflect.ValueOf(v).Len(); i++ {
		item_xml, _ := rpc2XML(reflect.ValueOf(v).Index(i).Interface())
		out += item_xml
	}
	out += "</data></array>"
	return
}

func time2XML(t time.Time) string {
	/*
		// TODO: find out whether we need to deal
		// here with TZ
		var tz string;
		zone, offset := t.Zone()
		if zone == "UTC" {
			tz = "Z"
		} else {
			tz = fmt.Sprintf("%03d00", offset / 3600 )
		}
	*/
	return fmt.Sprintf("<dateTime.iso8601>%04d%02d%02dT%02d:%02d:%02d</dateTime.iso8601>",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func base642XML(data []byte) string {
	str := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("<base64>%s</base64>", str)
}
