package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/anatolieGhebea/delegator/handlers"
	"github.com/anatolieGhebea/delegator/middleware"
	"github.com/anatolieGhebea/delegator/models"
)

func loadConfig() {

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
}

// ENTRY POINT
func main() {
	// set test data
	loadConfig()

	// define web endpoints
	http.HandleFunc("/info", middleware.LoggingMiddleware(middleware.JsonResponseHandler(handlers.InfoHandler)))
	http.HandleFunc("/trigger_update", middleware.LoggingMiddleware(middleware.JsonResponseHandler(handlers.TriggerHandler)))
	http.HandleFunc("/", middleware.LoggingMiddleware(middleware.JsonResponseHandler(handlers.NotFoundHandler)))

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
