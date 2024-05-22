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

	// check and create log file for the day
	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("logs/output_%s.log", currentDate)

	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
	}
	defer logFile.Close()

	// clean old log files
	// to move to a different function/file
	olderThan := fmt.Sprintf("+%d", models.Configuration.Server.LogRetentionDays)
	cleanLogFilesCmd := exec.Command("find", "logs/", "-type", "f", "-name", "output_*.log", "-mtime", olderThan, "-exec", "rm", "{}", ";")
	if err := cleanLogFilesCmd.Run(); err != nil {
		fmt.Fprintf(logFile, "error cleaning log files: %v", err)
	}

	eventSource := detectEventSource(req.Header)
	if eventSource == models.GitHubHook {
		handleGitHubEvent(w, req, logFile)
		return
	} else {
		handleGenericEvent(w, req, logFile)
		return
	}

}

func detectEventSource(header http.Header) models.EventSource {

	// check the header to set the correct request type
	values := header[models.HeaderGitHub]
	if len(values) > 0 {
		return models.GitHubHook
	}

	values = header[models.HeaderBitBucket]
	if len(values) > 0 {
		return models.BitBucketHook
	}

	return models.GenericHook
}

func handleGenericEvent(w http.ResponseWriter, req *http.Request, logFile *os.File) {
	// check request validity
	triggerRequest := models.GenericEventSource{}
	err := json.NewDecoder(req.Body).Decode(&triggerRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// check if TriggerEntry exists in Config by Name
	eventHook := models.EventHook{}
	found := false
	for _, item := range models.Configuration.Triggers {
		if item.Name == triggerRequest.Name {
			eventHook = item
			found = true
			break
		}
	}

	fmt.Fprintf(logFile, "Trigger action for %s.\n", triggerRequest.Name)
	if !found {
		notFoundRepository(w, logFile, triggerRequest.Name)
		return
	}

	if triggerRequest.SharedSecret != eventHook.SharedSecret || triggerRequest.Name != eventHook.Name {
		notAuthorized(w, logFile, triggerRequest.Name)
		return
	}

	currentBranch, err := getCurrentBranch(eventHook, logFile)
	if err != nil {
		http.Error(w, "An error occured while getting the current branch", http.StatusInternalServerError)
		return
	}

	if eventHook.SyncBranch == models.SpecificBranch && eventHook.BranchName != currentBranch {
		branchMismatch(w, logFile, eventHook.BranchName, currentBranch)
		return
	}

	// execute git pull command in the project folder to update the project
	cmd := exec.Command("git", "-C", eventHook.AbsolutePath, "pull", "origin", currentBranch)
	//	write the output to the log file
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Run(); err != nil {
		updateFailed(w, logFile, err)
		return
	}

	json.NewEncoder(w).Encode(models.Response{Message: "Operation completed"})

}
func handleGitHubEvent(w http.ResponseWriter, req *http.Request, logFile *os.File) {

	event := req.Header.Get(models.HeaderGitHub)

	if event == "ping" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Response{Message: "pong"})
		return
	} else if event != "push" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Response{Message: "accepted"})
		return
	}

	// check request validity
	triggerRequest := models.GitHubEventSource{}
	err := json.NewDecoder(req.Body).Decode(&triggerRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// check if TriggerEntry exists in Config by Name
	eventHook := models.EventHook{}
	found := false
	for _, item := range models.Configuration.Triggers {
		if item.RepositoryName == triggerRequest.Repository["full_name"] {
			eventHook = item
			found = true
			break
		}
	}

	if !found {
		notFoundRepository(w, logFile, triggerRequest.Repository["full_name"].(string))
		return
	}

	// check security

	//
	currentBranch, err := getCurrentBranch(eventHook, logFile)
	if err != nil {
		http.Error(w, "An error occured while getting the current branch", http.StatusInternalServerError)
		return
	}

	update_branch := ""
	if strings.Contains(triggerRequest.Ref, "heads") {
		parts := strings.Split(triggerRequest.Ref, "/")
		update_branch = parts[len(parts)-1]
	}

	if update_branch != currentBranch {
		// the event did not update the selected branch, no point in pulling the changes
		branchMismatch(w, logFile, update_branch, currentBranch)
		return
	}

	if eventHook.SyncBranch == models.SpecificBranch && eventHook.BranchName != currentBranch {
		// the pull request must be run only if the current branch is the one configured in the trigger
		branchMismatch(w, logFile, eventHook.BranchName, currentBranch)
		return
	}

	// execute git pull command in the project folder to update the project
	cmd := exec.Command("git", "-C", eventHook.AbsolutePath, "pull", "origin", currentBranch)
	//	write the output to the log file
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Run(); err != nil {
		updateFailed(w, logFile, err)
		return
	}

	json.NewEncoder(w).Encode(models.Response{Message: "Operation completed"})

}

func notFoundRepository(w http.ResponseWriter, logFile *os.File, repositoryName string) {
	fmt.Fprintf(logFile, "Repository not found in the configuration %s.\n", repositoryName)
	http.Error(w, "Repository not found in the configuration.", http.StatusNotFound)
}

func notAuthorized(w http.ResponseWriter, logFile *os.File, repositoryName string) {
	fmt.Fprintf(logFile, "Not authorized to trigger the repository %s.\n", repositoryName)
	http.Error(w, "Not authorized to trigger the repository.", http.StatusUnauthorized)
}

func branchMismatch(w http.ResponseWriter, logFile *os.File, configuredBranchName string, currentBranch string) {
	fmt.Fprintf(logFile, "The trigger is configured to run for branch %s current %s.\n", configuredBranchName, currentBranch)
	http.Error(w, "The trigger for the current project is missconfigured! Try later or contact the server administrator.", http.StatusBadRequest)
}

func updateFailed(w http.ResponseWriter, logFile *os.File, err error) {
	fmt.Printf("error executing git pull: %v", err)
	http.Error(w, "An error occured while executing the command. Check log file on the server for more detailes.", http.StatusInternalServerError)
}

func getCurrentBranch(eventHook models.EventHook, logFile *os.File) (string, error) {
	// Get current git branch
	getBranchNameCmd := exec.Command("git", "-C", eventHook.AbsolutePath, "rev-parse", "--abbrev-ref", "HEAD")
	branchName, err := getBranchNameCmd.Output()
	if err != nil {
		fmt.Printf("error getting current git branch: %v", err)
		return "", err
	}

	fmt.Fprintf(logFile, "Current branch: %s\n", branchName)
	return strings.TrimSpace(string(branchName)), nil
}
