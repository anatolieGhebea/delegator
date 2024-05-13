package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type BranchCheck string

const (
	CurrentBranch  BranchCheck = "current"
	SpecificBranch BranchCheck = "specific"
)

type Response struct {
	Message string `json:"message"`
}

type ServerConfig struct {
	Port string
}

type ProjectEntry struct {
	Name         string
	AbsolutePath string
	SharedSecret string
	SyncBranch   BranchCheck // Default: current. [ current > , <branch_name>]
	BranchName   string      //
}

type TriggerObject struct {
	ProjectName  string `json:"project_name"`
	SharedSecret string `json:"shared_secret"`
}

// test only
var serverConf ServerConfig = ServerConfig{Port: ":9999"} // set defalut port
var testProjectEntry ProjectEntry

// test only

func loadConfig() {
	// test only
	// Load configuration from file
	fmt.Println("Loading configuration...")

	serverConf = ServerConfig{Port: ":9180"} // customize port

	testProjectEntry = ProjectEntry{
		Name:         "SGA-development",
		AbsolutePath: "/Users/anatolieghebea/Dev/www/other/go_apps/simple-git-agent",
		SharedSecret: "asjhhdalhh8uyr84hiahesd894qawa",
		SyncBranch:   CurrentBranch,
		BranchName:   "main",
	}
}

// WEB ENDPOINTS (handlers implementation)
func jsonResponseHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(w, req)
	}
}

func infoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// fmt.Fprintf(w, "This is the information endpoint")

	response := Response{Message: "This is the information endpoint"}
	json.NewEncoder(w).Encode(response)
}

func updateHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// fmt.Fprintf(w, "true")
	if testProjectEntry.SyncBranch != SpecificBranch && testProjectEntry.BranchName != "main" {
		http.Error(w, "The trigger for the current project is missconfigured! Try later or contact the server administrator.", http.StatusInternalServerError)
		return
	}

	var triggerObject TriggerObject
	err := json.NewDecoder(req.Body).Decode(&triggerObject)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if triggerObject.SharedSecret != testProjectEntry.SharedSecret || triggerObject.ProjectName != testProjectEntry.Name {
		http.Error(w, "ProjectName and SharedKey don't match", http.StatusUnauthorized)
		return
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("cd %s && git pull origin ", testProjectEntry.AbsolutePath))
	if err := cmd.Run(); err != nil {
		http.Error(w, "Error while updating the project", http.StatusInternalServerError)
		return
	}

	response := Response{Message: fmt.Sprintf("Trigger is set to update %s branch, for project %s with key %s", string(testProjectEntry.SyncBranch), triggerObject.ProjectName, triggerObject.SharedSecret)}

	json.NewEncoder(w).Encode(response)

}

func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	// http.Error(w, "404 Not Found", http.StatusNotFound)
	http.Error(w, "404 Not Found", http.StatusNotFound)
	return
}

// END > WEB ENDPOINTS

// ENTRY POINT
func main() {
	// set test data
	loadConfig()

	// define web endpoints
	http.HandleFunc("/info", jsonResponseHandler(infoHandler))
	http.HandleFunc("/trigger_update", jsonResponseHandler(updateHandler))
	http.HandleFunc("/", jsonResponseHandler(notFoundHandler))

	go func() {
		fmt.Printf("Start server on port %s\n", serverConf.Port)

		if err := http.ListenAndServe(serverConf.Port, nil); err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
			os.Exit(1)
		}

	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Shutting down...")
}
