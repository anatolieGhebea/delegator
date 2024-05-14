package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/anatolieGhebea/simple-git-agent/models"
)

func TriggerHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// check request validity
	triggerRequest := models.TriggerRequest{}
	err := json.NewDecoder(req.Body).Decode(&triggerRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// check if TriggerEntry exists in Config by Name
	triggerEntry := models.TriggerEntry{}
	found := false
	for _, item := range models.Configuration.Triggers {
		if item.Name == triggerRequest.Name {
			triggerEntry = item
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Trigger not found, check the name field and try again.", http.StatusNotFound)
		return
	}

	if triggerRequest.SharedSecret != triggerEntry.SharedSecret || triggerRequest.Name != triggerEntry.Name {
		http.Error(w, "TriggerName and SharedKey don't match", http.StatusUnauthorized)
		return
	}

	// fmt.Fprintf(w, "true")
	current_simulate_branch := "main" // capire come prendere da git
	if triggerEntry.SyncBranch == models.SpecificBranch && triggerEntry.BranchName != current_simulate_branch {
		http.Error(w, "The trigger for the current project is missconfigured! Try later or contact the server administrator.", http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && git pull origin ", triggerEntry.AbsolutePath))
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		http.Error(w, "Error while updating the project", http.StatusInternalServerError)
		return
	}

	response := models.Response{Message: fmt.Sprintf("Trigger is set to update %s branch, for project %s with key %s", string(triggerEntry.SyncBranch), triggerEntry.Name, triggerEntry.SharedSecret)}

	json.NewEncoder(w).Encode(response)

}
