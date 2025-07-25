package main

import (
	"log"      // For logging messages
	"os"       // For interacting with the operating system (e.g., signals)
	"os/signal" // For handling OS signals (e.g., Ctrl+C)
	"syscall"  // For specific system calls (e.g., SIGINT, SIGTERM)

	"go-backend/config" // Import your config package
	"go-backend/pkg/db" // Import your db package for MongoDB connection
	"go-backend/internal/server" // Import your server package
	"go-backend/pkg/utils" // ADDED: Import your utils package to initialize WebSocket Hub
)

func main() {
	// 1. Load application configuration from environment variables.
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("Failed to load configuration.")
	}

	// 2. Connect to MongoDB.
	db.ConnectDB(cfg)
	defer db.DisconnectDB()

	// 3. Initialize the WebSocket Hub.
	// This creates the Hub instance and starts its Run() method in a goroutine.
	// The Hub will now manage WebSocket connections and message broadcasting.
	hub := utils.InitWebSocketHub()
	// The hub.Run() is already started internally by InitWebSocketHub as a goroutine.

	// 4. Initialize the Gin server.
	appServer := server.NewServer(cfg)

	// 5. Setup all API routes.
	// Pass the initialized WebSocket Hub to the server setup, so it can be used
	// by the WebSocket handler.
	appServer.SetupRoutes(hub) // MODIFIED: Pass the hub to SetupRoutes

	// 6. Start the Gin HTTP server in a goroutine.
	go func() {
		appServer.Run()
	}()

	// 7. Set up graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Perform any cleanup operations here before exiting.
	// The `defer db.DisconnectDB()` will handle MongoDB disconnection.
	log.Println("Server gracefully stopped.")
}
