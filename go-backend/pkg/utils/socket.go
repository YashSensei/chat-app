package utils

import (
	"encoding/json" // For marshaling/unmarshaling JSON messages
	"log"           // For logging messages
	"net/http"      // For HTTP status codes and upgrading HTTP to WebSocket
	"sync"          // For mutex to protect concurrent map access

	"go-backend/internal/models" // Import models for Message struct

	"github.com/gin-gonic/gin" // Gin context for handling WebSocket upgrade
	"github.com/gorilla/websocket" // WebSocket library for Go
	"go.mongodb.org/mongo-driver/bson/primitive" // For handling ObjectID
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
// CheckOrigin: allows cross-origin requests. In production, you'd want to restrict this.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow requests from your frontend origin.
		return r.Header.Get("Origin") == "http://localhost:5173"
	},
}

// Client represents a single WebSocket connection.
type Client struct {
	Conn *websocket.Conn
	UserID primitive.ObjectID // The ID of the user associated with this connection
}

// WebSocketMessage defines the generic structure for messages sent over WebSocket.
// This allows the frontend to identify the type of event.
type WebSocketMessage struct {
	Event   string      `json:"event"`   // e.g., "getOnlineUsers", "newMessage"
	Payload interface{} `json:"payload"` // The actual data for the event
}

// Hub manages the WebSocket clients (connections) and broadcasting.
// This is the Go equivalent of Socket.IO's server instance and userSocketMap.
type Hub struct {
	clients    map[primitive.ObjectID]*Client // Registered clients: {userID: *Client}
	broadcast  chan models.Message            // Channel for incoming messages from clients
	register   chan *Client                   // Channel for clients to register
	unregister chan *Client                   // Channel for clients to unregister
	mu         sync.Mutex                     // Mutex to protect concurrent access to `clients` map
}

// NewHub creates and returns a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[primitive.ObjectID]*Client),
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the Hub's goroutines to manage clients and broadcast messages.
// This should be run as a goroutine in your main function.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// A new client wants to register.
			h.mu.Lock() // Protect map access
			h.clients[client.UserID] = client
			h.mu.Unlock()
			h.sendOnlineUsers() // Notify all clients about updated online users
			log.Printf("User %s connected. Total online: %d", client.UserID.Hex(), len(h.clients))

		case client := <-h.unregister:
			// A client wants to unregister (disconnect).
			h.mu.Lock() // Protect map access
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				client.Conn.Close() // Close the WebSocket connection
			}
			h.mu.Unlock()
			h.sendOnlineUsers() // Notify all clients about updated online users
			log.Printf("User %s disconnected. Total online: %d", client.UserID.Hex(), len(h.clients))

		case message := <-h.broadcast:
			// A message needs to be broadcasted to the receiver.
			h.mu.Lock() // Protect map access
			receiverClient, ok := h.clients[message.ReceiverID]
			h.mu.Unlock()

			if ok {
				// Wrap the message in our generic WebSocketMessage structure.
				wsMessage := WebSocketMessage{
					Event:   "newMessage", // The event name the frontend expects
					Payload: message,      // The actual message data
				}
				msgJSON, err := json.Marshal(wsMessage) // Marshal the wrapped message
				if err != nil {
					log.Printf("Error marshaling message for receiver %s: %v", message.ReceiverID.Hex(), err)
					continue
				}
				if err := receiverClient.Conn.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
					log.Printf("Error sending message to receiver %s: %v", message.ReceiverID.Hex(), err)
					// Consider unregistering client if write fails consistently
				}
			} else {
				log.Printf("Receiver %s is offline. Message not sent via WebSocket.", message.ReceiverID.Hex())
				// In a real app, you might queue this message for offline delivery or push notifications.
			}
		}
	}
}

// sendOnlineUsers sends the list of currently online user IDs to all connected clients.
func (h *Hub) sendOnlineUsers() {
	h.mu.Lock()
	defer h.mu.Unlock()

	onlineUserIDs := make([]string, 0, len(h.clients))
	for userID := range h.clients {
		onlineUserIDs = append(onlineUserIDs, userID.Hex())
	}

	// Create a structured message for online users, similar to Socket.IO's event.
	// The frontend will expect an event like "getOnlineUsers".
	// Now using the generic WebSocketMessage struct.
	onlineUsersMessage := WebSocketMessage{
		Event:   "getOnlineUsers",
		Payload: onlineUserIDs, // The list of user IDs
	}

	msgJSON, err := json.Marshal(onlineUsersMessage)
	if err != nil {
		log.Printf("Error marshaling online users message: %v", err)
		return
	}

	// Iterate over all clients and send the online users list.
	for _, client := range h.clients {
		if err := client.Conn.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
			log.Printf("Error sending online users to client %s: %v", client.UserID.Hex(), err)
			// Potentially unregister this client if write fails
		}
	}
}

// WebSocketHandler upgrades the HTTP connection to a WebSocket connection.
// It registers the new client with the Hub.
// This will be used as a Gin route handler.
func WebSocketHandler(c *gin.Context, hub *Hub) {
	// Get the authenticated user from the context (set by AuthMiddleware)
	userAny, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized - User not found in context"})
		return
	}
	loggedInUser := userAny.(models.User)

	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to establish WebSocket connection"})
		return
	}

	// Create a new Client instance and register it with the Hub.
	client := &Client{Conn: conn, UserID: loggedInUser.ID}
	hub.register <- client // Send client to the register channel

	// Start a goroutine to continuously read messages from the WebSocket connection.
	// This loop keeps the connection alive and handles incoming messages (if any, though chat is outbound).
	go func() {
		defer func() {
			hub.unregister <- client // Ensure client is unregistered on exit
			conn.Close()
		}()

		for {
			// ReadMessage blocks until a message is received or an error occurs.
			// We primarily send messages from server to client, but this keeps the connection open.
			// If clients were sending messages to the server, this is where they'd be processed.
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error for user %s: %v", loggedInUser.ID.Hex(), err)
				}
				break // Exit the loop on error (e.g., client disconnected)
			}
			// If a message was read, you could process it here if your frontend sends messages
			// via this same WebSocket connection for other purposes.
		}
	}()
}

// EmitNewMessage is a public function to send a new message via the Hub's broadcast channel.
// This will be called from your chat handler (SendMessage) to send real-time updates.
var currentHub *Hub // Global reference to the Hub

// InitWebSocketHub initializes the global Hub. Call this once in main.go.
func InitWebSocketHub() *Hub {
	currentHub = NewHub()
	go currentHub.Run() // Start the Hub's goroutine
	return currentHub
}

// GetHub returns the initialized global Hub instance.
func GetHub() *Hub {
	return currentHub
}

// EmitNewMessage sends a message to the broadcast channel of the global Hub.
// This is the function that will be called from `chat.handler.go`'s `SendMessage` method.
func EmitNewMessage(message models.Message) {
	if currentHub != nil {
		currentHub.broadcast <- message
	} else {
		log.Println("WebSocket Hub not initialized. Cannot emit message.")
	}
}
