package utils

import (
	"net/http"
	"runtime/debug"
)

func RecoverHttpPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recoveredValue := recover(); recoveredValue != nil {
				if recoveredValue == http.ErrAbortHandler {
					// net/http catches this panic and aborts only the response, without crashing the server.
					panic(recoveredValue)
				}

				Logger.Error(
					"panic in handler",
					"panic_cause", recoveredValue,
					"stack", string(debug.Stack()),
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
