package utils

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	Data any `json:"data,omitempty"`
}

type RestErr struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

// ReadJSON tries to read the body of a request and converts it into JSON
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSON takes a response status code and arbitrary data and writes a json response to the client
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for _, header := range headers {
			for key, value := range header {
				w.Header()[key] = value
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_, err = w.Write(out)
		if err != nil {
			return err
		}
	}

	return nil
}

// ErrorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload RestErr
	payload.Error = err.Error()
	payload.Status = statusCode

	return WriteJSON(w, statusCode, payload)
}
