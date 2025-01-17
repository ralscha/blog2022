package main

import (
	"fmt"
	"github.com/AJAYK-01/passwordless-go/passwordless"
	"github.com/google/uuid"
	"net/http"
	"webauthn.rasc.ch/internal/request"
	"webauthn.rasc.ch/internal/response"
)

func (app *application) authenticateHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.sessionManager.Destroy(r.Context()); err != nil {
		response.InternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type CreateTokenInput struct {
	Username string `json:"username"`
}

type CreateTokenOutput struct {
	Token string `json:"token"`
}

func (app *application) createToken(w http.ResponseWriter, r *http.Request) {
	var input CreateTokenInput
	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		response.BadRequest(w, err)
		return
	}

	if input.Username == "" {
		response.BadRequest(w, fmt.Errorf("username is required"))
		return
	}

	params := passwordless.RegisterRequest{
		UserId:   uuid.New().String(),
		Username: input.Username,
	}
	resp, err := app.passwordlessClient.CreateRegisterToken(params)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, CreateTokenOutput{Token: resp.Token})
}

type SigninInput struct {
	Token string `json:"token"`
}

func (app *application) signin(w http.ResponseWriter, r *http.Request) {
	var input SigninInput
	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		response.BadRequest(w, err)
		return
	}

	resp, err := app.passwordlessClient.VerifySignin(input.Token)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	if !resp.Success {
		if err := app.sessionManager.Destroy(r.Context()); err != nil {
			response.InternalServerError(w, err)
			return
		}
		response.Forbidden(w)
		return
	}
	app.sessionManager.Put(r.Context(), "userID", resp.UserId)

	w.WriteHeader(http.StatusNoContent)
}
