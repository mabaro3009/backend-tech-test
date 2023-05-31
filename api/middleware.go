package api

import (
	"errors"
	"net/http"
)

func JSONResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func RecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil { // IDE marks this as wrong but its not (IDE bug)
				RespondError(w, Error{
					Err:        errors.New("ERR_PANIC"),
					HTTPStatus: http.StatusInternalServerError,
					Reason:     Internal,
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
