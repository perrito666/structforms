package main

import (
	"fmt"
	"strconv"
)

type sample struct {
	FirstName string
	LastName  string
	Email     string `html:"email"`
	Website   string `html:"url"`
	Phone     string `html:"telephone"`
	Password  string `html:"secret"`
	Age       int
}

// ProcessPost implements POSTFillable.
func (s *sample) ProcessPOST(post map[string]string) error {
	if value, ok := post["FirstName"]; ok {
		s.FirstName = value
	}
	if value, ok := post["LastName"]; ok {
		s.LastName = value
	}
	if value, ok := post["Email"]; ok {
		s.Email = value
	}
	if value, ok := post["Website"]; ok {
		s.Website = value
	}
	if value, ok := post["Phone"]; ok {
		s.Phone = value
	}
	if value, ok := post["Password"]; ok {
		s.Password = value
	}
	if value, ok := post["Age"]; ok {
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("unexpected value for Age %q", value)
		}
		s.Age = int(i)
	}
	return nil
}

func main() {
	fields, err := getFields(&sample{})
	if err != nil {
		fmt.Printf("<h1>cannot obtain fields for type: %v</h1>", err)
		return
	}
	html := doHTML(fields)
	fmt.Println(html)

	html, err = doFilledForm(&sample{
		FirstName: "Horacio",
		LastName:  "Duran",
		Email:     "horacio.duran@gmail.com",
		Website:   "http://perri.to",
		Phone:     "+555555555",
		Password:  "a big secret",
		Age:       32,
	})
	if err != nil {
		fmt.Printf("<h1>cannot create a filled form: %v</h1>", err)
	}
	fmt.Println(html)

}
