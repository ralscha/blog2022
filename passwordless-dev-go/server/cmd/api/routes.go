package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"net/http"
	"time"
	"webauthn.rasc.ch/internal/config"
	"webauthn.rasc.ch/internal/response"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(response.NotFound)
	mux.MethodNotAllowed(response.MethodNotAllowed)

	mux.Use(middleware.RealIP)
	if app.config.Environment == config.Development {
		mux.Use(middleware.Logger)
	}

	mux.Use(middleware.Recoverer)
	mux.Use(httprate.LimitAll(1_000, 1*time.Minute))
	mux.Use(middleware.Timeout(15 * time.Second))
	mux.Use(middleware.NoCache)

	mux.Route("/api/v1", func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)
		r.Group(func(r chi.Router) {
			r.Post("/logout", app.logoutHandler)
			r.Post("/create-token", app.createToken)
			r.Post("/signin", app.signin)
		})

		/*
			r.Group(func(r chi.Router) {
				r.Post("/authentication/start", app.authenticationStart)
				r.Post("/authentication/finish", app.authenticationFinish)
				r.Post("/registration/start", app.registrationStart)
				r.Post("/registration/finish", app.registrationFinish)
			})
		*/

		r.Group(func(r chi.Router) {
			r.Use(app.authenticatedOnly)
			r.Post("/authenticate", app.authenticateHandler)
			r.Get("/secret", app.secret)
		})
	})

	return mux
}
