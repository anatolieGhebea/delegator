package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func LoggingMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("Request Headers:\n")
		for key, values := range req.Header {
			fmt.Printf("  %s: %s\n", key, strings.Join(values, ", "))
		}

		handler(w, req)
	}
}
