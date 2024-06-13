package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anatolieGhebea/delegator/models"
)

func InfoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// fmt.Fprintf(w, "This is the information endpoint")

	response := models.Response{Message: "This is the information endpoint"}
	json.NewEncoder(w).Encode(response)
}
