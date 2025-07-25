package chat

import (
	"context"    // For context with MongoDB operations
	"fmt"        // For formatted error messages
	//"log"        // For logging errors
	"net/http"   // For HTTP status codes
	"time"       // For handling timestamps

	"go-backend/internal/models" // Import models for User and Message structs
	"go-backend/pkg/db" // Import db to access MongoDB client
	"go-backend/pkg/utils" // Import utils for socket operations AND CloudinaryService

	"github.com/gin-gonic/gin" // Gin context for handling requests
	"go.mongodb.org/mongo-driver/bson" // For MongoDB queries
	"go.mongodb.org/mongo-driver/bson/primitive" // For ObjectID
	"go.mongodb.org/mongo-driver/mongo/options" // For MongoDB find options (e.g., sort)
)

// Struct for SendMessage request body
type SendMessageRequest struct {
	Text  string `json:"text,omitempty"`  // Message text, optional
	Image string `json:"image,omitempty"` // Base64 encoded image, optional
}

// ChatHandler struct holds dependencies for chat operations.
// ADDED: CloudinaryService dependency
type ChatHandler struct {
	CloudinaryService *utils.CloudinaryService // Add Cloudinary service
}

// NewChatHandler creates a new instance of ChatHandler.
// MODIFIED: Accepts CloudinaryService
func NewChatHandler(cldService *utils.CloudinaryService) *ChatHandler { // Changed signature
	return &ChatHandler{
		CloudinaryService: cldService,
	}
}

// GetUsersForSidebar retrieves a list of users for the sidebar, excluding the logged-in user.
// Mirrors backend/src/controllers/message.controller.js -> getUsersForSidebar
func (h *ChatHandler) GetUsersForSidebar(c *gin.Context) {
	// Get the authenticated user from the context (set by AuthMiddleware)
	userAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authenticated user not found in context"})
		return
	}
	loggedInUser := userAny.(models.User) // Type assertion to models.User

	var users []models.User // Slice to hold the retrieved users
	usersCollection := db.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find all users where _id is not equal to the logged-in user's ID.
	// The projection (options.Find().SetProjection) is used to exclude the password field.
	cursor, err := usersCollection.Find(ctx, bson.M{"_id": bson.M{"$ne": loggedInUser.ID}}, options.Find().SetProjection(bson.M{"password": 0}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error fetching users: %v", err)})
		return
	}
	defer cursor.Close(ctx) // Ensure the cursor is closed after use

	// Iterate through the cursor and decode each document into a models.User struct.
	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error decoding users: %v", err)})
		return
	}

	// Prepare response data to match frontend expectation (converting ObjectID to hex string)
	responseUsers := make([]gin.H, len(users))
	for i, user := range users {
		responseUsers[i] = gin.H{
			"_id":        user.ID.Hex(),
			"fullName":   user.FullName,
			"email":      user.Email,
			"profilePic": user.ProfilePic,
			"createdAt":  user.CreatedAt,
			"updatedAt":  user.UpdatedAt,
		}
}

	c.JSON(http.StatusOK, responseUsers)
}

// GetMessages retrieves messages between the logged-in user and a specific receiver.
// Mirrors backend/src/controllers/message.controller.js -> getMessages
func (h *ChatHandler) GetMessages(c *gin.Context) {
	// Get receiver ID from URL parameters
	receiverIDParam := c.Param("id")
	receiverID, err := primitive.ObjectIDFromHex(receiverIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID format"})
		return
	}

	// Get the authenticated user from the context
	userAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authenticated user not found in context"})
		return
	}
	loggedInUser := userAny.(models.User)
	myID := loggedInUser.ID

	var messages []models.Message // Slice to hold the retrieved messages
	messagesCollection := db.DB.Collection("messages")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Construct the query using $or to find messages where:
	// (senderId = myID AND receiverId = receiverID) OR
	// (senderId = receiverID AND receiverId = myID)
	filter := bson.M{
		"$or": []bson.M{
			{"senderId": myID, "receiverId": receiverID},
			{"senderId": receiverID, "receiverId": myID},
		},
	}

	// Sort messages by createdAt to ensure chronological order
	findOptions := options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}})

	cursor, err := messagesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error fetching messages: %v", err)})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error decoding messages: %v", err)})
		return
	}

	// Prepare response data (converting ObjectIDs to hex strings for frontend)
	responseMessages := make([]gin.H, len(messages))
	for i, msg := range messages {
		responseMessages[i] = gin.H{
			"_id":        msg.ID.Hex(),
			"senderId":   msg.SenderID.Hex(),
			"receiverId": msg.ReceiverID.Hex(),
			"text":       msg.Text,
			"image":      msg.Image,
			"createdAt":  msg.CreatedAt,
			"updatedAt":  msg.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, responseMessages)
}

// SendMessage handles sending a new message between two users.
// Mirrors backend/src/controllers/message.controller.js -> sendMessage
func (h *ChatHandler) SendMessage(c *gin.Context) {
	// Get receiver ID from URL parameters
	receiverIDParam := c.Param("id")
	receiverID, err := primitive.ObjectIDFromHex(receiverIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID format"})
		return
	}

	// Get the authenticated user from the context (sender)
	userAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authenticated user not found in context"})
		return
	}
	loggedInUser := userAny.(models.User)
	senderID := loggedInUser.ID

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body format"})
		return
	}

	// Ensure at least text or image is provided
	if req.Text == "" && req.Image == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message text or image is required"})
		return
	}

	var imageUrl string
	if req.Image != "" {
		// INTEGRATED CLOUDINARY: Upload the base64 image to Cloudinary
		uploadResultURL, err := h.CloudinaryService.UploadImage(req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error uploading image: %v", err)})
			return
		}
		imageUrl = uploadResultURL // Use the secure URL from Cloudinary
	}


	// Create new message
	newMessage := models.Message{
		ID:         primitive.NewObjectID(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Text:       req.Text,
		Image:      imageUrl,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	messagesCollection := db.DB.Collection("messages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert message into database
	_, err = messagesCollection.InsertOne(ctx, newMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error saving message: %v", err)})
		return
	}

	// UNCOMMENTED: Emit the new message via WebSocket for real-time update
	utils.EmitNewMessage(newMessage)

	// Respond with the newly created message
	c.JSON(http.StatusCreated, gin.H{
		"_id":        newMessage.ID.Hex(),
		"senderId":   newMessage.SenderID.Hex(),
		"receiverId": newMessage.ReceiverID.Hex(),
		"text":       newMessage.Text,
		"image":      newMessage.Image,
		"createdAt":  newMessage.CreatedAt,
		"updatedAt":  newMessage.UpdatedAt,
	})
}
