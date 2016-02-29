// Copyright 2016 Horacio Duran.
// Licensed under the MIT, see LICENCE file for details.

package main

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/template"
)

const (
	TagSecret = "secret"
	TagEmail  = "email"
	TagURL    = "url"
	TagTel    = "telephone"
)

const (
	HTMLFieldText     = "text"
	HTMLFieldPassword = "password"
	HTMLFieldEmail    = "email"
	HTMLFieldURL      = "url"
	HTMLFieldTel      = "tel" // This seems to be useful in phone browsers
	HTMLFieldNumber   = "number"
)

const (
	inputTemplate          = `<label for="{{.Name}}">{{.Name}}</label><input name="{{.Name}}" type="{{.Type}}" />`
	inputWithValueTemplate = `<label for="{{.Name}}">{{.Name}}</label><input name="{{.Name}}" type="{{.Type}}" value="{{.Value}}" />`
	formTemplate           = `<form name="goform" method="POST">%s</form>`
)

var (
	htmlLineTemplate          = template.Must(template.New("htmlInputLine").Parse(inputTemplate))
	htmlLineWithValueTemplate = template.Must(template.New("htmlInputLine").Parse(inputWithValueTemplate))
)

var GOHTML = map[string]string{
	"string":  HTMLFieldText,
	"int":     HTMLFieldNumber,
	"float":   HTMLFieldNumber,
	"int64":   HTMLFieldNumber,
	"float64": HTMLFieldNumber,
}

// POSTFillable represents a type that can fill its attributes from the
// results of a POSTed form.
type POSTFillable interface {
	// ProcessPOST takes a map containing key/values from an HTTP POST (or GET)
	// and fills the struct with them.
	ProcessPOST(postFormValues map[string]string) error
}

// getFields will return a map of the field names and their HTML
// form input type.
// Since this function is only for educational purposes it handles
// only the most basic types and will return error when a more complex
// one is passed.
func getFields(goType POSTFillable) (map[string]string, error) {
	// Only valid for structs for the purposes of this example.
	t := reflect.TypeOf(goType)
	// We require a struct for this sample so if this is a ptr, lets get
	// the concrete type pointed to.
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("interface must be a struct is %v", t.Kind())
	}
	fieldCount := t.NumField()
	// we know how many fields there are, we assume there are all valid.
	ret := make(map[string]string, fieldCount)
	for i := 0; i < fieldCount; i++ {
		f := t.Field(i)
		// Normally I could very well use the type and a function with
		// a switch going over the possible types here
		// but it would make this example a bit more complicated.
		fieldHTMLType, valid := GOHTML[f.Type.Name()]
		if !valid {
			return nil, fmt.Errorf("cannot find an HTML equivalent for %q of type ", f.Name)
		}

		// Adding tags to an attribute you can add extra metadata.
		tag := f.Tag.Get("html")
		if fieldHTMLType == HTMLFieldText {
			switch tag {
			case TagSecret:
				fieldHTMLType = HTMLFieldPassword
			case TagEmail:
				fieldHTMLType = HTMLFieldEmail
			case TagURL:
				fieldHTMLType = HTMLFieldURL
			case TagTel:
				fieldHTMLType = HTMLFieldTel
			}
		}

		ret[f.Name] = fieldHTMLType
	}
	return ret, nil
}

// htmlInput represents the data that can be used in htmlLoneTemplate.
type htmlInput struct {
	Name  string
	Type  string
	Value string
}

// doHTML returns a string containing the HTML form for the passed
// fields map.
func doHTML(fields map[string]string) string {
	htmlLines := make([]string, len(fields))
	fieldNames := make([]string, len(fields))
	fieldNo := 0
	for fieldName := range fields {
		fieldNames[fieldNo] = fieldName
		fieldNo++
	}
	sort.Sort(sort.StringSlice(fieldNames))
	for lineNo, fieldName := range fieldNames {
		fieldType := fields[fieldName]
		field := htmlInput{
			Name: fieldName,
			Type: fieldType,
		}
		var htmlLine bytes.Buffer
		htmlLineTemplate.Execute(&htmlLine, field)
		htmlLines[lineNo] = htmlLine.String()
	}
	return fmt.Sprintf(formTemplate, strings.Join(htmlLines, "\n"))
}

// valueStringer tries to return a string representing the value
// of the passed reflect.Value and a boolean indicating if it
// was possible.
func valueStringer(value reflect.Value) (string, bool) {
	var stringValue string
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := value.Int()
		stringValue = fmt.Sprintf("%d", v)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := value.Uint()
		stringValue = fmt.Sprintf("%d", v)
	case reflect.Float32, reflect.Float64:
		v := value.Float()
		stringValue = fmt.Sprintf("%f", v)
	case reflect.String:
		stringValue = value.String()
	default:
		return "", false
	}
	return stringValue, true

}

// doFilledForm returns a string containing an HTML form with values
// extracted from the passed POSTFillable or error if this is not possible.
func doFilledForm(contents POSTFillable) (string, error) {
	fields, err := getFields(contents)
	if err != nil {
		return "", fmt.Errorf("cannot obtain fields: %v", err)
	}
	t := reflect.ValueOf(contents)
	// We require a struct for this sample so if this is a ptr, lets get
	// the concrete type pointed to.
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	htmlLines := make([]string, len(fields))
	fieldNames := make([]string, len(fields))
	fieldNo := 0
	for fieldName := range fields {
		fieldNames[fieldNo] = fieldName
		fieldNo++
	}
	for lineNo, fieldName := range fieldNames {
		fieldType := fields[fieldName]
		value, valid := valueStringer(t.FieldByName(fieldName))
		if !valid {
			return "", fmt.Errorf("cannot determine the value for %q", fieldName)
		}
		field := htmlInput{
			Name:  fieldName,
			Type:  fieldType,
			Value: value,
		}
		var htmlLine bytes.Buffer
		htmlLineWithValueTemplate.Execute(&htmlLine, field)
		htmlLines[lineNo] = htmlLine.String()
	}
	return fmt.Sprintf(formTemplate, strings.Join(htmlLines, "\n")), nil
}
