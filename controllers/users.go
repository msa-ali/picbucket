package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/msa-ali/picbucket/context"
	"github.com/msa-ali/picbucket/models"
	"github.com/msa-ali/picbucket/utils"
)

type Users struct {
	Templates struct {
		New            Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
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
	var data struct {
		Email    string
		Password string
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data.Email = r.PostForm.Get("email")
	data.Password = r.PostForm.Get("password")
	if !utils.ValidateEmail(data.Email) || len(data.Password) < 4 {
		http.Error(w, "create user: invalid email or password format", http.StatusBadRequest)
		return
	}
	user, err := u.UserService.Create(data.Email, data.Password)
	if err != nil {
		u.Templates.New.Execute(w, r, data, err)
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

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		fmt.Println(err)
		// TODO: handle other cases like user doesn;t exist,
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	val := url.Values{
		"token": {pwReset.Token},
	}
	resetUrl := utils.GetEnv().ServerUrl + "/reset-pw?" + val.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetUrl)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	// don't render the reset token here as we need the user to confirm they have
	// access to the email account to verify their identity
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")
	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		// TODO: handle other cases like user doesn;t exist,
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		fmt.Println(err)
		// TODO: handle other cases like user doesn;t exist,
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// sign the user in now that their password has been reset
	// any errors from this point onwards should redirect the user to the sign in page
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	utils.SetCookie(w, utils.CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}
