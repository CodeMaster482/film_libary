package response

import (
	"net/http"

	"films_library/pkg/logger"

	"github.com/mailru/easyjson"
)

const (
	InvalidURLParameter = "invalid url parameter"
	InvalidBodyRequest  = "invalid input body"
	ForbiddenUser       = "user has no rights"
)

// const minErrorToLogCode = 500

//easyjson:json
type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

//easyjson:json
type ResponseError struct {
	Status int    `json:"status"`
	ErrMes string `json:"message"`
}

type NilBody struct{}

func NIL() NilBody {
	return NilBody{}
}

func ErrorResponse(w http.ResponseWriter, code int, message string, log logger.Interface) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	errorMsg := ResponseError{
		Status: code,
		ErrMes: message,
	}

	// if code < minErrorToLogCode {
	// 	log.Info("invalid request: %v:", err)
	// } else {
	// 	log.Error(err.Error())
	// }

	// Marshal response using easyjson
	_, _, err := easyjson.MarshalToHTTPResponseWriter(errorMsg, w)
	if err != nil {
		log.Error("Error failed to marshal error message: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		if _, writeErr := w.Write([]byte("Can't encode error message into json, message: " + message)); writeErr != nil {
			log.Error("Error writing response: %s", writeErr.Error())
		}
	}
}

func SuccessResponse[T any](w http.ResponseWriter, status int, response T) {
	date := Response{Status: status, Body: response}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// Marshal response using easyjson
	_, _, err := easyjson.MarshalToHTTPResponseWriter(date, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
