package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}

// note: gorilla/schema is 3rd party lib which can be used to handle complex forms
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email := r.PostForm.Get("email")
	pass := r.PostForm.Get("password")
	if email == "" || pass == "" {
		http.Error(w, "invalid form values", http.StatusBadRequest)
		return
	}
	fmt.Println(email, pass)
	fmt.Fprint(w, "User signed up!!")
}
