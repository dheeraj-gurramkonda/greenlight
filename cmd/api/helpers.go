package main
package main

import (
	"encoding/json/v2"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data, json.Deterministic(true))
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readIDParam(r *http.Request) (int, error) {
	// ---------- Reading the URL parameters ----------
	// httprouter stores the values of named URL parameters (like :id in
	// /v1/movies/:id) inside the request's context. ParamsFromContext pulls
	// them out so we can look them up by name.
	params := httprouter.ParamsFromContext(r.Context())

	// ---------- Converting and validating the id ----------
	// Parameter values always arrive as strings, so ByName("id") returns text
	// such as "7". strconv.Atoi converts it to an integer and returns an error
	// if it isn't a number (e.g. "abc"). We also reject any id below 1, since
	// valid movie IDs start at 1. If either check fails, we return 0 along
	// with a generic error, so the caller never has to deal with a half-valid
	// value.
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
