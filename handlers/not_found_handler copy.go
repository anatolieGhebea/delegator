package handlers

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	// http.Error(w, "404 Not Found", http.StatusNotFound)
	http.Error(w, "404 Not Found", http.StatusNotFound)
}
