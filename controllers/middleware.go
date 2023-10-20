package controllers

import (
	"net/http"

	"github.com/msa-ali/picbucket/context"
	"github.com/msa-ali/picbucket/models"
	"github.com/msa-ali/picbucket/utils"
)

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		token, err := utils.ReadCookie(r, utils.CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithUser(r.Context(), user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}
