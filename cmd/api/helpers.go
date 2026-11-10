package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"greenlight.dheeraj.net/internal/models"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	opts := []json.Options{
		json.Deterministic(true),
		jsontext.Multiline(true),
		json.WithMarshalers(models.RuntimeMarshaler),
	}
	js, err := json.Marshal(data, opts...)
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
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.UnmarshalRead(r.Body, dst)
	if err != nil {
		var synErr *jsontext.SyntacticError
		var semErr *json.SemanticError

		switch {
		case errors.As(err, &synErr):
			return convertSyntacticErrorToBadRequest(synErr)

		case errors.As(err, &semErr):
			return &badRequestError{
				Message: "there was a semantic error",
			}

		default:
			return err
		}
	}

	return nil
}

func convertSyntacticErrorToBadRequest(synErr *jsontext.SyntacticError) *badRequestError {
	if errors.Is(synErr, io.ErrUnexpectedEOF) && synErr.ByteOffset == 0 {
		return &badRequestError{Message: "request body must not be empty"}
	}

	var buf strings.Builder

	pointer := string(synErr.JSONPointer)
	if pointer == "" {
		pointer = "request body"
	}
	fmt.Fprintf(&buf, "%s contains malformed JSON", pointer)

	offset := synErr.ByteOffset
	if offset > 0 {
		fmt.Fprintf(&buf, " at position %d", offset)
	}

	detail := synErr.Err.Error()
	if detail == "unexpected EOF" {
		detail = "unexpected end of JSON"
	}
	fmt.Fprintf(&buf, ": %s", detail)

	return &badRequestError{Message: buf.String()}
}
