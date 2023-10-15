package controllers

import (
	"fmt"
	"net/http"

	"github.com/msa-ali/picbucket/models"
	"github.com/msa-ali/picbucket/utils"
)

type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
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
	if !utils.ValidateEmail(email) || len(pass) < 4 {
		http.Error(w, "create user: invalid email or password format", http.StatusBadRequest)
		return
	}
	user, err := u.UserService.Create(email, pass)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error while creating user", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User created: %+v", user)
}
