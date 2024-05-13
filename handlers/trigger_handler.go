package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/anatolieGhebea/simple-git-agent/models"
)

func UpdateHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	projectEntry := models.ProjectEntry{
		Name:         "SGA-development",
		AbsolutePath: "/Users/anatolieghebea/Dev/www/other/go_apps/simple-git-agent",
		SharedSecret: "asjhhdalhh8uyr84hiahesd894qawa",
		SyncBranch:   models.CurrentBranch,
		BranchName:   "main",
	}

	// fmt.Fprintf(w, "true")
	if projectEntry.SyncBranch != models.SpecificBranch && projectEntry.BranchName != "main" {
		http.Error(w, "The trigger for the current project is missconfigured! Try later or contact the server administrator.", http.StatusInternalServerError)
		return
	}

	var triggerObject models.TriggerObject
	err := json.NewDecoder(req.Body).Decode(&triggerObject)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if triggerObject.SharedSecret != projectEntry.SharedSecret || triggerObject.ProjectName != projectEntry.Name {
		http.Error(w, "ProjectName and SharedKey don't match", http.StatusUnauthorized)
		return
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && git pull origin ", projectEntry.AbsolutePath))
	if err := cmd.Run(); err != nil {
		http.Error(w, "Error while updating the project", http.StatusInternalServerError)
		return
	}

	response := models.Response{Message: fmt.Sprintf("Trigger is set to update %s branch, for project %s with key %s", string(projectEntry.SyncBranch), triggerObject.ProjectName, triggerObject.SharedSecret)}

	json.NewEncoder(w).Encode(response)

}
