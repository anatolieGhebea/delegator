package middleware

import (
	"net/http"
)

func JsonResponseHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(w, req)
	}
}
