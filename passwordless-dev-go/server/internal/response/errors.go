package response

import (
	"fmt"
	"log/slog"
	"net/http"
)

func errorMessage(w http.ResponseWriter, status int, message string, headers http.Header) {
	JSONWithHeaders(w, status, map[string]string{"error": message}, headers)
}

func InternalServerError(w http.ResponseWriter, err error) {
	slog.Error(err.Error(), "error", err)

	message := "The server encountered a problem and could not process your request"
	errorMessage(w, http.StatusInternalServerError, message, nil)
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	message := "The requested resource could not be found"
	errorMessage(w, http.StatusNotFound, message, nil)
}

func Forbidden(w http.ResponseWriter) {
	message := "You are not authorized to access this resource"
	errorMessage(w, http.StatusForbidden, message, nil)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	errorMessage(w, http.StatusMethodNotAllowed, message, nil)
}

func BadRequest(w http.ResponseWriter, err error) {
	errorMessage(w, http.StatusBadRequest, err.Error(), nil)
}
