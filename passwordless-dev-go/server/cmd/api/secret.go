package main

import (
	"log/slog"
	"net/http"
	"webauthn.rasc.ch/internal/response"
)

type SecretOutput struct {
	Message string `json:"message"`
}

func (app *application) secret(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "userID")
	slog.Info("fetch secret", "User ID", userID)

	response.JSON(w, http.StatusOK, SecretOutput{
		Message: "This is a secret message",
	})
}
