package main

import (
	"fmt"
	"testing"
)

type wrongFormValue int

func (w wrongFormValue) ProcessPOST(map[string]string) error {
	return nil
}

func TestGetFieldsFailsWithoutStruct(t *testing.T) {
	var testValue wrongFormValue
	testValue = 1
	_, err := getFields(testValue)
	expectedMessage := "interface must be a struct is int"
	errorMessage := err.Error()
	if errorMessage != expectedMessage {
		t.Logf("error should be %q but is %q", expectedMessage, errorMessage)
		t.Fail()
	}
}

type rightForm struct {
}

func (r rightForm) ProcessPOST(map[string]string) error {
	return nil
}

func TestGetFieldsDereferencesPtr(t *testing.T) {
	testValue := rightForm{}
	_, err := getFields(testValue)
	if err != nil {
		t.Logf("error should be nil, is: %v", err)
		t.Fail()
	}
}

type formWithFields struct {
	FirstName string
	LastName  string
	Email     string `html:"email"`
	Website   string `html:"url"`
	Phone     string `html:"telephone"`
	Password  string `html:"secret"`
	Age       int
}

func (f formWithFields) ProcessPOST(map[string]string) error {
	return nil
}

func TestGetFieldsReturnsFields(t *testing.T) {
	testValue := formWithFields{}
	result, err := getFields(testValue)
	if err != nil {
		t.Logf("error should be nil, is: %v", err)
		t.Fail()
	}
	resultLenght := len(result)
	if resultLenght != 7 {
		t.Logf("result should contain 7 fields, contains %d", resultLenght)
		t.Fail()
	}
	expected := map[string]string{
		"Email":     "email",
		"Website":   "url",
		"Phone":     "tel",
		"Password":  "password",
		"Age":       "number",
		"FirstName": "text",
		"LastName":  "text",
	}
	for fieldName, fieldType := range expected {
		value, ok := result[fieldName]
		if !ok {
			t.Logf("key %q should be present in the map but is not", fieldName)
			t.Logf("obtained map is %#v", result)
			t.Logf("expected map is %#v", expected)
			t.Fail()
		}
		if value != fieldType {
			t.Logf("expected field type for field %q is %q but the obtained values is %q", fieldName, fieldType, value)
			t.Fail()
		}
	}
}

func TestDoHTMLReturnsExpectedHTML(t *testing.T) {
	inputMap := map[string]string{
		"Email":     "email",
		"Website":   "url",
		"Phone":     "tel",
		"Password":  "password",
		"Age":       "number",
		"FirstName": "text",
		"LastName":  "text",
	}
	expectedHTML := `<form name="goform" method="POST"><label for="Age">Age</label><input name="Age" type="number" />
<label for="Email">Email</label><input name="Email" type="email" />
<label for="FirstName">FirstName</label><input name="FirstName" type="text" />
<label for="LastName">LastName</label><input name="LastName" type="text" />
<label for="Password">Password</label><input name="Password" type="password" />
<label for="Phone">Phone</label><input name="Phone" type="tel" />
<label for="Website">Website</label><input name="Website" type="url" /></form>`
	outputHTML := doHTML(inputMap)
	if outputHTML != expectedHTML {
		t.Logf("got: \n%q\nexpected: \n%q", outputHTML, expectedHTML)
		t.Fail()
	}
}

func TestDoFillHTMLReturnsExpectedHTML(t *testing.T) {
	inputStruct := formWithFields{
		FirstName: "Horacio",
		LastName:  "Duran",
		Email:     "horacio.duran@gmail.com",
		Website:   "http://perri.to",
		Phone:     "+555555555",
		Password:  "a big secret",
		Age:       32,
	}
	expectedHTML := `<form name="goform" method="POST"><label for="FirstName">FirstName</label><input name="FirstName" type="text" value="Horacio" />
<label for="LastName">LastName</label><input name="LastName" type="text" value="Duran" />
<label for="Email">Email</label><input name="Email" type="email" value="horacio.duran@gmail.com" />
<label for="Website">Website</label><input name="Website" type="url" value="http://perri.to" />
<label for="Phone">Phone</label><input name="Phone" type="tel" value="+555555555" />
<label for="Password">Password</label><input name="Password" type="password" value="a big secret" />
<label for="Age">Age</label><input name="Age" type="number" value="32" /></form>`
	outputHTML, err := doFilledForm(inputStruct)
	if err != nil {
		t.Logf("no error expected, obtained %v", err)
		t.Fail()
	}
	fmt.Println(outputHTML)
	if outputHTML != expectedHTML {
		t.Logf("got: \n%q\nexpected: \n%q", outputHTML, expectedHTML)
		t.Fail()
	}
}
