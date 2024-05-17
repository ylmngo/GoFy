package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp = append(resp, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)

	return nil
}

// lint:ignore U1000 Ignore unused functions temporarily for debugging
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// limit the size of the request body to 1MB
	mb := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(mb))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // only allow fields that are part of the dst struct

	if err := dec.Decode(dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field: %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains badly formed JSON (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", mb)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	// decode to an anonymous struct to ensure that there are no more JSON values
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must only contain a sinlge JSON value")
	}

	return nil
}

func (app *application) nestMiddlewares(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			mw := mws[i]
			next = mw(next)
		}

		return next
	}
}
