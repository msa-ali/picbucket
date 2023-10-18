package utils

import (
	"fmt"
	"net/http"
)

const (
	CookieSession = "session"
)

func GetCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
}

func SetCookie(w http.ResponseWriter, name, value string) {
	cookie := GetCookie(name, value)
	http.SetCookie(w, cookie)
}

func ReadCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("%s : %w", name, err)
	}
	return c.Value, nil
}
