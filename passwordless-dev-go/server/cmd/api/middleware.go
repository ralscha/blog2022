package main

import (
	"net/http"
	"webauthn.rasc.ch/internal/response"
)

func (app *application) authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := app.sessionManager.GetString(r.Context(), "userID")
		if userID != "" {
			next.ServeHTTP(w, r)
		} else {
			response.Forbidden(w)
		}
	})

}
