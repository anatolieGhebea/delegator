package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/anatolieGhebea/simple-git-agent/handlers"
	"github.com/anatolieGhebea/simple-git-agent/middleware"
	"github.com/anatolieGhebea/simple-git-agent/models"
)

// test only
// var Config models.Config = models.Config{
// 	Server:   models.Server{Port: ":9919"}, // set default port to 9919
// 	Triggers: []models.TriggerEntry{},
// }

// test only

func loadConfig() {
	// test only
	// Load configuration from file
	fmt.Println("Loading configuration...")

	// check if json config file exists
	if _, err := os.Stat("config/config.json"); os.IsNotExist(err) {
		fmt.Println("Config file not found")
		os.Exit(2)
	}

	file, _ := os.ReadFile("config/config.json")
	err := json.Unmarshal(file, &models.Configuration)
	if err != nil {
		fmt.Println("Error while loading configuration")
		os.Exit(3)
	}

	// Print the server info from config
	fmt.Println("Server:")
	fmt.Printf("Server port: %s\n", models.Configuration.Server.Port)

	// Print the loaded items
	fmt.Printf("Items loaded (%d):\n", len(models.Configuration.Triggers))
	for _, item := range models.Configuration.Triggers {
		fmt.Printf("Name: %s, \tPath: %s\n", item.Name, item.AbsolutePath)
	}

	// Config.Server = models.Server{Port: ":9180"} // customize port
}

// WEB ENDPOINTS (handlers implementation)

// END > WEB ENDPOINTS

// ENTRY POINT
func main() {
	// set test data
	loadConfig()

	// define web endpoints
	http.HandleFunc("/info", middleware.JsonResponseHandler(handlers.InfoHandler))
	http.HandleFunc("/trigger_update", middleware.JsonResponseHandler(handlers.TriggerHandler))
	http.HandleFunc("/", middleware.JsonResponseHandler(handlers.NotFoundHandler))

	go func() {
		fmt.Printf("Start server on port %s\n", models.Configuration.Server.Port)

		if err := http.ListenAndServe(models.Configuration.Server.Port, nil); err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
			os.Exit(1)
		}

	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Shutting down...")
}
