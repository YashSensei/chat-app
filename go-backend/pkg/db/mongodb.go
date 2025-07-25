package db

import (
	"context" // For managing request-scoped values, cancellation signals, and deadlines
	"fmt"     // For formatted I/O (like printing to console)
	"log"     // For logging messages, especially errors
	"time"    // For specifying timeouts

	"go-backend/config" // Import your config package. IMPORTANT: Replace "chat-app-backend" with your actual Go module name from go.mod

	"go.mongodb.org/mongo-driver/mongo"          // The main MongoDB driver package
	"go.mongodb.org/mongo-driver/mongo/options"  // For setting client options
	"go.mongodb.org/mongo-driver/mongo/readpref" // For pinging the database
)

// Global variables to hold the MongoDB client and database instance.
// These will be initialized once and then used throughout the application.
var(
	Client *mongo.Client
	DB *mongo.Database
)

// ConnectDB establishes a connection to MongoDB.
// It takes a pointer to your application's Config struct, which contains the MongoDB URI.
func ConnectDB(cfg *config.Config){
	// 1. Create a new context with a timeout for the connection attempt.
	//    It's good practice to set a reasonable timeout for network operations.
	//    Example: 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is cancelled when ConnectDB exits

	// 2. Create a new MongoDB client instance.
	//    Use `options.Client().ApplyURI()` to specify the connection string from your config.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBURI))
	if err != nil{
		// If connection fails, log a fatal error and exit the application.
		log.Fatalf("MongoDB connection error: %v", err)
	}

	// 3. Ping the primary database to verify the connection is alive and working.
	//    This helps catch issues even if `Connect` didn't return an error immediately.
	err = client.Ping(ctx , readpref.Primary())
	if err != nil{
		// If ping fails, log a fatal error and exit the application.
		log.Fatalf("MongoDB ping error: %v", err)
	}

	// 4. If connection and ping are successful, assign the client and the desired database
	//    to the global variables. 
	Client = client
	DB = client.Database("chat-db") // Make sure "chat-db" matches your database name

	fmt.Println("MongoDB connected successfully!")
}

// DisconnectDB closes the MongoDB connection gracefully.
// This function should be called when your application is shutting down.
func DisconnectDB(){
	// 1. Create a new context for the disconnection with a timeout.
	ctx, cancel :=context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // ensure the context is cancelled

	// 2. Check if the client is not nil before attempting to disconnect.
	if Client == nil{
		log.Println("MongoDB client is already nil, nothing to disconnect.")
		return
	}

	// 3. Disconnect the global MongoDB client.
	err := Client.Disconnect(ctx)
	if err != nil{
		// Log the error but don't fatally exit, as this is part of a graceful shutdown.
		log.Printf("Error disconnecting from MongoDB: %v", err)
		return
	}
	fmt.Println("MongoDB disconnected successfully.")
}