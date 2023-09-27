package main

import (
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
}
