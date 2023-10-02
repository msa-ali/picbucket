package main

import (
	"errors"
	"fmt"
	"html/template"
	"os"
)

type User struct {
	Name string
	Age  int
	Meta UserMeta
	Bio  string //template.HTML
}

type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("cmd/exp/hello.html")
	if err != nil {
		panic(err)
	}
	user := User{
		Name: "Altamash Ali",
		Age:  27,
		Meta: UserMeta{
			Visits: 10,
		},
		Bio: `<script>alert("haha, you have been hacked !!")</script>`,
	}
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

	// errorf
	fmt.Println("\n", CreateOrg())
}

func Connect() error {
	return errors.New("connection failed")
}

func CreateUser() error {
	err := Connect()
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	// continue processing
	return nil
}

func CreateOrg() error {
	err := CreateUser()
	if err != nil {
		return fmt.Errorf("create Org: %w", err)
	}
	// continue processing
	return nil
}
