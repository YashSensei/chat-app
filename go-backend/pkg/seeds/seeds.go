package seeds

import (
	"context" // For context with MongoDB operations
	//"fmt"     // For formatted output - REMOVED: Not used in this file
	"log"     // For logging messages
	"time"    // For timestamps

	"go-backend/config" // Import config for MongoDB URI
	"go-backend/internal/models" // Import models for User struct
	"go-backend/pkg/db" // Import db for MongoDB connection

	"go.mongodb.org/mongo-driver/bson" // For MongoDB queries
	"go.mongodb.org/mongo-driver/bson/primitive" // For ObjectID
	"go.mongodb.org/mongo-driver/mongo" // For MongoDB client operations
	"golang.org/x/crypto/bcrypt" // For password hashing
)

// SeedUsers defines the initial user data to be inserted.
// This mirrors the `seedUsers` array in your Node.js `user.seed.js`.
var SeedUsers = []struct {
	Email      string
	FullName   string
	Password   string
	ProfilePic string
}{
	// Female Users
	{
		Email:      "emma.thompson@example.com",
		FullName:   "Emma Thompson",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/women/1.jpg",
	},
	{
		Email:      "olivia.miller@example.com",
		FullName:   "Olivia Miller",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/women/2.jpg",
	},
	{
		Email:      "sophia.davis@example.com",
		FullName:   "Sophia Davis",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/women/3.jpg",
	},
	{
		Email:      "ava.wilson@example.com",
		FullName:   "Ava Wilson",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/women/4.jpg",
	},
	{
		Email:      "isabella.brown@example.com",
		FullName:   "Isabella Brown",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/women/5.jpg",
	},
	{
		Email:      "mia.johnson@example.com",
		FullName:   "Mia Johnson",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/women/6.jpg",
	},
	{
		Email:      "charlotte.williams@example.com",
		FullName:   "Charlotte Williams",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/women/7.jpg",
	},
	{
		Email:      "amelia.garcia@example.com",
		FullName:   "Amelia Garcia",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/women/8.jpg",
	},

	// Male Users
	{
		Email:      "james.anderson@example.com",
		FullName:   "James Anderson",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/men/1.jpg",
	},
	{
		Email:      "william.clark@example.com",
		FullName:   "William Clark",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/men/2.jpg",
	},
	{
		Email:      "benjamin.taylor@example.com",
		FullName:   "Benjamin Taylor",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/men/3.jpg",
	},
	{
		Email:      "lucas.moore@example.com",
		FullName:   "Lucas Moore",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/men/4.jpg",
	},
	{
		Email:      "henry.jackson@example.com",
		FullName:   "Henry Jackson",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/men/5.jpg",
	},
	{
		Email:      "alexander.martin@example.com",
		FullName:   "Alexander Martin",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/men/6.jpg",
	},
	{
		Email:      "daniel.rodriguez@example.com",
		FullName:   "Daniel Rodriguez",
		Password:   "123456",
		ProfilePic: "https://randomuser.me/api/portraits/men/7.jpg",
	},
}

// SeedDatabase connects to MongoDB and inserts the predefined users.
// This function mirrors the `seedDatabase` function in your Node.js `user.seed.js`.
func SeedDatabase() {
	// Load configuration (needed for MongoDB URI)
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("Failed to load configuration for seeding.")
	}

	// Connect to MongoDB
	db.ConnectDB(cfg)
	defer db.DisconnectDB() // Ensure disconnection on exit

	usersCollection := db.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Starting database seeding...")

	// Iterate through the seed users and insert them
	for _, seedUser := range SeedUsers {
		// Check if user already exists by email to prevent duplicates
		var existingUser models.User
		err := usersCollection.FindOne(ctx, bson.M{"email": seedUser.Email}).Decode(&existingUser)
		if err == nil {
			log.Printf("User with email %s already exists, skipping.", seedUser.Email)
			continue // Skip if user already exists
		}
		if err != mongo.ErrNoDocuments {
			log.Printf("Error checking for existing user %s: %v", seedUser.Email, err)
			continue // Log error and continue to next user
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password for %s: %v", seedUser.Email, err)
			continue // Log error and continue to next user
		}

		// Create a new User model instance
		newUser := models.User{
			ID:         primitive.NewObjectID(),
			FullName:   seedUser.FullName,
			Email:      seedUser.Email,
			Password:   string(hashedPassword), // Store hashed password as string
			ProfilePic: seedUser.ProfilePic,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		// Insert the new user into the database
		_, err = usersCollection.InsertOne(ctx, newUser)
		if err != nil {
			log.Printf("Error inserting user %s: %v", seedUser.Email, err)
			continue // Log error and continue to next user
		}
		log.Printf("Successfully seeded user: %s", newUser.Email)
	}

	log.Println("Database seeding completed.")
}

// main function for standalone execution of seeding.
// This is typically run once via `go run pkg/seeds/seeds.go`.
func init() {
    // This `init` function will run automatically when this package is imported.
    // However, for a standalone seeding script, you'd typically call SeedDatabase
    // from a `main` function if this were its own executable.
    // For our structure, we'll create a separate `cmd/seed/main.go` for execution.
}
