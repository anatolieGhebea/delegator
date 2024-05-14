package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

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

	// check and create log file for the day
	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("logs/output_%s.log", currentDate)

	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
	}
	defer logFile.Close()

	fmt.Fprintf(logFile, "Trigger action for %s.\n", triggerRequest.Name)
	if !found {
		// print error to log file
		fmt.Fprintf(logFile, "Trigger config not found for %s.\n", triggerRequest.Name)
		http.Error(w, "Trigger not found, check the name field and try again.", http.StatusNotFound)
		return
	}

	if triggerRequest.SharedSecret != triggerEntry.SharedSecret || triggerRequest.Name != triggerEntry.Name {
		fmt.Fprintf(logFile, "TriggerName and SharedKey don't match %s > %s.\n", triggerRequest.Name, triggerRequest.SharedSecret)
		http.Error(w, "TriggerName and SharedKey don't match", http.StatusUnauthorized)
		return
	}

	// Get current git branch
	getBranchNameCmd := exec.Command("git", "-C", triggerEntry.AbsolutePath, "rev-parse", "--abbrev-ref", "HEAD")
	branchName, err := getBranchNameCmd.Output()
	if err != nil {
		fmt.Printf("error getting current git branch: %v", err)
		return
	}

	currentBranch := strings.TrimSpace(string(branchName))
	fmt.Fprintf(logFile, "Current branch: %s\n", currentBranch)
	// current_simulate_branch := "main" // capire come prendere da git

	// check if the trigger is set to a specific branch and if the current branch is the same
	if triggerEntry.SyncBranch == models.SpecificBranch && triggerEntry.BranchName != currentBranch {
		fmt.Fprintf(logFile, "The trigger is configured to run for branch %s.\n", triggerEntry.BranchName)
		http.Error(w, "The trigger for the current project is missconfigured! Try later or contact the server administrator.", http.StatusInternalServerError)
		return
	}

	// execute git pull command in the project folder to update the project
	cmd := exec.Command("git", "-C", triggerEntry.AbsolutePath, "pull", "origin", currentBranch)
	//	write the output to the log file
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Run(); err != nil {
		fmt.Printf("error executing git pull: %v", err)
		http.Error(w, "An error occured while executing the command. Check log file on the server for more detailes.", http.StatusInternalServerError)
		return
	}

	response := models.Response{Message: "Operation completed"}

	// clean old log files
	olderThan := fmt.Sprintf("+%d", models.Configuration.Server.LogRetentionDays)
	cleanLogFilesCmd := exec.Command("find", "logs/", "-type", "f", "-name", "output_*.log", "-mtime", olderThan, "-exec", "rm", "{}", ";")
	if err := cleanLogFilesCmd.Run(); err != nil {
		fmt.Fprintf(logFile, "error cleaning log files: %v", err)
	}

	json.NewEncoder(w).Encode(response)

}
