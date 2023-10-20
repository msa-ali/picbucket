package controllers

import (
	"fmt"
	"net/http"

	"github.com/msa-ali/picbucket/context"
	"github.com/msa-ali/picbucket/models"
	"github.com/msa-ali/picbucket/utils"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
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
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println("Create: ", err)
		// @TODO: Show user a warning about not being able to signin
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	utils.SetCookie(w, utils.CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) Authenticate(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Authenticate(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error while signing user", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println("Authenticate: ", err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	utils.SetCookie(w, utils.CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}

func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	token, err := utils.ReadCookie(r, utils.CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	utils.DeleteCookie(w, utils.CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}
