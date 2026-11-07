package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// ---------- Declaring the helper method ----------
// readIDParam is a small helper attached to the application struct so any
// handler can reuse it instead of repeating the same parsing code. It takes
// the incoming request and returns the id from the URL as an int, or an error
// if the id is missing or invalid. Returning (int, error) lets each handler
// decide how to respond to a bad id, usually with a 404 Not Found.

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

	// ---------- Returning the valid id ----------
	// If we get here the id passed every check, so we return it with a nil
	// error to signal success.
	return id, nil
}