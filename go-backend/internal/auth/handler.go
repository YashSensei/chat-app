package auth

import (
	"context"    // For context with MongoDB operations
	"fmt"        // For formatted error messages
	"net/http"   // For HTTP status codes
	"time"       // For handling timestamps

	"go-backend/config" // Import config for JWT secret and other settings
	"go-backend/internal/models" // Import models for User struct
	"go-backend/pkg/db" // Import db to access MongoDB client
	"go-backend/pkg/utils" // Import utils for JWT generation AND CloudinaryService

	"github.com/gin-gonic/gin" // Gin context for handling requests
	"go.mongodb.org/mongo-driver/bson" // For MongoDB queries
	"go.mongodb.org/mongo-driver/bson/primitive" // For ObjectID
	"go.mongodb.org/mongo-driver/mongo" // For MongoDB client operations and error checking
	"golang.org/x/crypto/bcrypt" // For password hashing
)

// Structs for request bodies (input validation)
type SignupRequest struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	ProfilePic string `json:"profilePic" binding:"required"` // This will be the base64 string
}

// AuthHandler struct holds dependencies for authentication operations.
// ADDED: CloudinaryService dependency
type AuthHandler struct {
	Config          *config.Config
	CloudinaryService *utils.CloudinaryService // Add Cloudinary service
}

// NewAuthHandler creates a new instance of AuthHandler.
// MODIFIED: Accepts CloudinaryService
func NewAuthHandler(cfg *config.Config, cldService *utils.CloudinaryService) *AuthHandler {
	return &AuthHandler{
		Config:          cfg,
		CloudinaryService: cldService,
	}
}

// Signup handles new user registration.
// Mirrors backend/src/controllers/auth.controller.js -> signup
func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "All fields are required or invalid format"})
		return
	}

	// Check if user already exists
	var existingUser models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.DB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email already exists"})
		return
	}
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Internal server error checking user: %v", err)})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}

	// Create new user
	newUser := models.User{
		ID:         primitive.NewObjectID(), // MongoDB will generate this, but good to set explicitly or omit
		FullName:   req.FullName,
		Email:      req.Email,
		Password:   string(hashedPassword),
		ProfilePic: "", // Default empty string
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Insert user into database
	_, err = db.DB.Collection("users").InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error saving user: %v", err)})
		return
	}

	// Generate JWT token and set cookie
	if err := utils.GenerateToken(newUser.ID, c, h.Config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error generating token: %v", err)})
		return
	}

	// Respond with user data (excluding password)
	c.JSON(http.StatusCreated, gin.H{
		"_id":        newUser.ID.Hex(), // Convert ObjectID to hex string for frontend
		"fullName":   newUser.FullName,
		"email":      newUser.Email,
		"profilePic": newUser.ProfilePic,
	})
}

// Login handles user authentication.
// Mirrors backend/src/controllers/auth.controller.js -> login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password format"})
		return
	}

	// Find user by email
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.DB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Internal server error finding user: %v", err)})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
		return
	}

	// Generate JWT token and set cookie
	if err := utils.GenerateToken(user.ID, c, h.Config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error generating token: %v", err)})
		return
	}

	// Respond with user data (excluding password)
	c.JSON(http.StatusOK, gin.H{
		"_id":        user.ID.Hex(),
		"fullName":   user.FullName,
		"email":      user.Email,
		"profilePic": user.ProfilePic,
	})
}

// Logout handles user logout by clearing the JWT cookie.
// Mirrors backend/src/controllers/auth.controller.js -> logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the "jwt" cookie by setting its maxAge to 0.
	// CORRECTED: Removed http.SameSiteStrictMode as it's not accepted by this Gin SetCookie signature.
	c.SetCookie("jwt", "", -1, "/", "", h.Config.NodeEnv == "production", true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// UpdateProfile handles updating the user's profile picture.
// Mirrors backend/src/controllers/auth.controller.js -> updateProfile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Get the authenticated user from the context (set by AuthMiddleware)
	userAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found in context"})
		return
	}
	user := userAny.(models.User) // Type assertion

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Profile pic is required"})
		return
	}

	// INTEGRATED CLOUDINARY: Upload the base64 image to Cloudinary
	uploadResultURL, err := h.CloudinaryService.UploadImage(req.ProfilePic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error uploading profile picture: %v", err)})
		return
	}

	newProfilePicURL := uploadResultURL // Use the secure URL from Cloudinary

	// Update user in database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Define the update operation using bson.M for a map-like update document
	update := bson.M{
		"$set": bson.M{
			"profilePic": newProfilePicURL,
			"updatedAt":  time.Now(), // Manually update updatedAt
		},
	}

	// Find and update the user by their ID
	_, err = db.DB.Collection("users").UpdateByID(ctx, user.ID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error updating profile: %v", err)})
		return
	}

	// Fetch the updated user to return the latest data
	var updatedUser models.User
	err = db.DB.Collection("users").FindOne(ctx, bson.M{"_id": user.ID}).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error fetching updated user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"_id":        updatedUser.ID.Hex(),
		"fullName":   updatedUser.FullName,
		"email":      updatedUser.Email,
		"profilePic": updatedUser.ProfilePic,
	})
}

// CheckAuth returns the currently authenticated user's data.
// Mirrors backend/src/controllers/auth.controller.js -> checkAuth
func (h *AuthHandler) CheckAuth(c *gin.Context) {
	// Get the authenticated user from the context (set by AuthMiddleware)
	userAny, exists := c.Get("user")
	if !exists {
		// This case should ideally not be hit if middleware works correctly,
		// but it's a good safeguard.
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}
	user := userAny.(models.User) // Type assertion

	// Respond with user data (excluding password)
	c.JSON(http.StatusOK, gin.H{
		"_id":        user.ID.Hex(),
		"fullName":   user.FullName,
		"email":      user.Email,
		"profilePic": user.ProfilePic,
	})
}
