package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ---------- Declaring the routes method ----------
// The routes method hangs off the application struct, so it can reach the
// config and logger, and so the handlers it registers (also methods on
// application) share the same dependencies. It returns an http.Handler, which
// is exactly what main() plugs into the server's Handler field. Returning the
// interface rather than the concrete router type keeps main() independent of
// which router library is used.

func (app *application) routes() http.Handler {
	// ---------- Creating the router ----------
	// httprouter.New() builds a new router instance. It is chosen over the
	// standard mux because it supports method-based routing and named URL
	// parameters like :id, and it automatically replies with 404 for unknown
	// paths and 405 for known paths with the wrong method.
	router := httprouter.New()

	// ---------- Registering the endpoints ----------
	// Each line maps an HTTP method and URL pattern to a handler function.
	// GET /v1/healthcheck reports that the service is up. POST /v1/movies
	// creates a new movie. GET /v1/movies/:id fetches one movie, where :id is a
	// placeholder the handler can read from the request (for example, /v1/movies/7).
	// HandlerFunc adapts a standard http.HandlerFunc-style method so httprouter
	// can use it, and the /v1 prefix makes it easy to introduce a v2 API later
	// without breaking existing clients.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	// ---------- Returning the router ----------
	// The configured router is returned to the caller. Since httprouter's
	// Router implements http.Handler, it can be given directly to the server,
	// which then dispatches every incoming request through the routes above.
	return router
}