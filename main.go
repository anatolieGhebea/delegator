package main

import (
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
var serverConf models.ServerConfig = models.ServerConfig{Port: ":9999"} // set defalut port
// test only

func loadConfig() {
	// test only
	// Load configuration from file
	fmt.Println("Loading configuration...")
	serverConf = models.ServerConfig{Port: ":9180"} // customize port
}

// WEB ENDPOINTS (handlers implementation)

// END > WEB ENDPOINTS

// ENTRY POINT
func main() {
	// set test data
	loadConfig()

	// define web endpoints
	http.HandleFunc("/info", middleware.JsonResponseHandler(handlers.InfoHandler))
	http.HandleFunc("/trigger_update", middleware.JsonResponseHandler(handlers.InfoHandler))
	http.HandleFunc("/", middleware.JsonResponseHandler(handlers.NotFoundHandler))

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
