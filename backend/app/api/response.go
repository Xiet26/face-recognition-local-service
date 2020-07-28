package api

import (
	"encoding/json"
	"git.cyradar.com/utilities/data/datatable"
	"git.cyradar.com/utilities/data/providers/mongo"
	"git.cyradar.com/utilities/data/providers/redis"
	"net/http"
	"strings"
)

const MESSAGE_INTERNAL_SERVER_ERROR = "Internal Server Error"

var jsonContentType = []string{"application/json; charset=utf-8"}

type ResponseBody struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseError(w http.ResponseWriter, r *http.Request, err error) { //check error type
	switch err.(type) {
	case ValidationError:
		WriteJSON(w, http.StatusUnprocessableEntity, datatable.ResponseBody{
			Code:    http.StatusUnprocessableEntity,
			Message: strings.Title(err.Error()),
		})
		return

	case mongo.DatabaseExecutionError, redis.DatabaseExecutionError:
		message := err.Error()

		logLevel, ok := r.Context().Value("logLevel").(int)
		if ok && logLevel <= 2 {
			message = MESSAGE_INTERNAL_SERVER_ERROR
		}

		WriteJSON(w, http.StatusInternalServerError, datatable.ResponseBody{
			Code:    http.StatusInternalServerError,
			Message: strings.Title(message),
		})
		return

	case AuthenticationError:
		WriteJSON(w, http.StatusUnauthorized, datatable.ResponseBody{
			Code:    http.StatusUnauthorized,
			Message: strings.Title(err.Error()),
		})
		return
	case AuthorizationError:
		WriteJSON(w, http.StatusForbidden, datatable.ResponseBody{
			Code:    http.StatusForbidden,
			Message: strings.Title(err.Error()),
		})
		return
	}

	WriteJSON(w, http.StatusUnprocessableEntity, datatable.ResponseBody{
		Code:    http.StatusUnprocessableEntity,
		Message: strings.Title(err.Error()),
	})
}

func WriteJSON(w http.ResponseWriter, code int, obj interface{}) error {
	w.WriteHeader(code)
	writeContentType(w, jsonContentType)

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		for _, v := range value {
			w.Header().Set("Content-Type", v)
		}
	}
}
