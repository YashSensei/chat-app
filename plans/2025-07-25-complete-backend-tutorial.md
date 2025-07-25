# Complete Go Backend Tutorial: Building a Real-Time Chat Application

## Table of Contents
1. [Introduction to Go and Project Structure](#introduction)
2. [Configuration Management](#configuration)
3. [Database Layer - MongoDB Integration](#database)
4. [Data Models and BSON Mapping](#models)
5. [Gin Framework and Server Setup](#gin-server)
6. [Authentication System - JWT Implementation](#authentication)
7. [Authentication Handlers and Endpoints](#auth-handlers)
8. [WebSocket Implementation](#websockets)
9. [Chat Functionality](#chat-functionality)
10. [Image Upload Integration](#image-upload)
11. [Frontend-Backend Communication](#frontend-backend)
12. [Application Startup Flow](#startup-flow)
13. [Security Considerations](#security)

---

## 1. Introduction to Go and Project Structure {#introduction}

### What is Go?
Go (also called Golang) is a programming language developed by Google. It's designed for building scalable, concurrent applications. Key features include:
- **Static typing**: Variables have fixed types
- **Garbage collection**: Automatic memory management
- **Goroutines**: Lightweight threads for concurrency
- **Simple syntax**: Easy to read and write
- **Fast compilation**: Quick build times

### Understanding Go Modules
Your project uses Go modules for dependency management. The `go.mod` file defines:

```go
module go-backend  // This is your module name
go 1.23.4         // Go version requirement
```

When you import packages in your code like `"go-backend/config"`, Go looks for the `config` directory relative to your module root.

### Project Structure Explained
```
go-backend/
├── cmd/api/main.go           # Application entry point
├── config/config.go          # Configuration management
├── internal/                 # Private application code
│   ├── auth/                # Authentication logic
│   ├── chat/                # Chat functionality
│   ├── models/              # Data structures
│   └── server/              # Server setup
├── pkg/                     # Public/reusable packages
│   ├── db/                  # Database connection
│   ├── utils/               # Utility functions
│   └── seeds/               # Database seeding
└── go.mod                   # Module definition
```

**Why this structure?**
- `cmd/`: Contains main applications (entry points)
- `internal/`: Private code that can't be imported by other projects
- `pkg/`: Public code that could be reused
- This follows Go community conventions

---

## 2. Configuration Management {#configuration}

### Understanding Environment Variables
Your application needs configuration values like database URLs and API keys. Instead of hardcoding these, you use environment variables.

**File: `config/config.go`**

```go
package config

import(
    "log"
    "os"
    "github.com/joho/godotenv"
)
```

**Package Declaration**: Every Go file starts with `package`. Files in the same directory should have the same package name.

**Imports**: Go's way of including external libraries:
- `"log"`: Built-in logging
- `"os"`: Operating system interface
- `"github.com/joho/godotenv"`: Third-party library for loading .env files

### The Config Struct
```go
type Config struct{
    Port                 string
    MongoDBURI           string
    JWTSecret            string
    CloudinaryCloudName  string
    CloudinaryAPIKey     string
    CloudinaryAPISecret  string
    NodeEnv              string
}
```

**Struct**: Go's way of grouping related data. Think of it like a class in other languages, but simpler.

### Loading Configuration
```go
func LoadConfig() *Config{
    err := godotenv.Load()
    if err != nil{
        log.Println("No .env file found...")
    }
    return &Config{
        Port: getEnv("PORT", "5000"),
        // ... other fields
    }
}
```

**Function Syntax**: `func FunctionName() ReturnType`
- `*Config` means "pointer to Config" - more efficient for large structs
- `&Config{}` creates a new Config instance and returns its memory address

**Error Handling**: Go uses explicit error handling:
```go
err := godotenv.Load()
if err != nil {
    // Handle error
}
```

### Helper Function
```go
func getEnv(key string, defaultvalue string) string{
    if value, exists := os.LookupEnv(key); exists{
        return value
    }
    return defaultvalue
}
```

**Multiple Return Values**: Go functions can return multiple values. `LookupEnv` returns the value and a boolean indicating if it exists.

---

## 3. Database Layer - MongoDB Integration {#database}

### Understanding the MongoDB Driver
Your app uses the official MongoDB Go driver. It provides:
- Connection management
- Query building
- Document mapping

**File: `pkg/db/mongodb.go`**

### Global Variables
```go
var(
    Client *mongo.Client
    DB *mongo.Database
)
```

**Global Variables**: These are accessible throughout your application. The `*` indicates pointers.

### Connection Function
```go
func ConnectDB(cfg *config.Config){
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBURI))
    if err != nil{
        log.Fatalf("MongoDB connection error: %v", err)
    }
}
```

**Context**: Go's way of handling timeouts and cancellation:
- `context.WithTimeout()` creates a context that expires after 10 seconds
- `defer cancel()` ensures cleanup happens when function exits
- `defer` runs the function when the current function returns

**Connection Process**:
1. Create a timeout context
2. Attempt to connect to MongoDB
3. If successful, ping to verify connection
4. Store client and database references globally

### Why This Pattern?
- **Global Access**: Any part of your app can access the database
- **Single Connection**: Reuse one connection pool throughout the app
- **Graceful Shutdown**: Properly close connections when app stops

---

## 4. Data Models and BSON Mapping {#models}

### Understanding BSON
BSON (Binary JSON) is MongoDB's data format. Go structs map to BSON documents using struct tags.

**File: `internal/models/user.go`**

### User Model
```go
type User struct {
    ID primitive.ObjectID `bson:"_id,omitempty"`
    Email string `bson:"email"`
    FullName string `bson:"fullName"`
    Password string `bson:"password"`
    ProfilePic string `bson:"profilePic,omitempty"`
    CreatedAt time.Time `bson:"createdAt"`
    UpdatedAt time.Time `bson:"updatedAt"`
}
```

**Struct Tags**: The backtick strings are struct tags that tell the MongoDB driver how to map fields:
- `bson:"_id,omitempty"`: Maps to MongoDB's `_id` field, omits if empty
- `bson:"email"`: Maps to `email` field in MongoDB
- `omitempty`: Don't include field if it's the zero value

**ObjectID**: MongoDB's unique identifier type. In Go, it's `primitive.ObjectID`.

### Message Model
```go
type Message struct {
    ID primitive.ObjectID `bson:"_id,omitempty"`
    SenderID primitive.ObjectID `bson:"senderId"`
    ReceiverID primitive.ObjectID `bson:"receiverId"`
    Text string `bson:"text,omitempty"`
    Image string `bson:"image,omitempty"`
    CreatedAt time.Time `bson:"createdAt"`
    UpdatedAt time.Time `bson:"updatedAt"`
}
```

**Relationships**: Instead of foreign keys, MongoDB uses ObjectIDs to reference other documents:
- `SenderID` and `ReceiverID` reference User documents
- This is similar to foreign keys in SQL databases

---

## 5. Gin Framework and Server Setup {#gin-server}

### What is Gin?
Gin is a web framework for Go. It provides:
- HTTP routing
- Middleware support
- JSON binding
- Template rendering

**File: `internal/server/server.go`**

### Server Struct
```go
type Server struct {
    Engine *gin.Engine
    Config *config.Config
}
```

**Composition**: Go doesn't have inheritance, but uses composition. The Server struct contains a Gin engine and configuration.

### Creating a New Server
```go
func NewServer(cfg *config.Config) *Server {
    if cfg.NodeEnv == "production" {
        gin.SetMode(gin.ReleaseMode)
    } else {
        gin.SetMode(gin.DebugMode)
    }
    
    engine := gin.Default()
    
    return &Server{
        Engine: engine,
        Config: cfg,
    }
}
```

**Constructor Pattern**: Go doesn't have constructors, but uses functions that return initialized structs.

**Gin Modes**:
- `DebugMode`: Verbose logging, helpful for development
- `ReleaseMode`: Minimal logging, optimized for production

### Setting Up Routes
```go
func (s *Server) SetupRoutes(hub *utils.Hub) {
    // CORS middleware
    s.Engine.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
}
```

**Method Receiver**: `(s *Server)` makes this a method on the Server struct. `s` is the receiver variable.

**CORS (Cross-Origin Resource Sharing)**: Allows your frontend (running on port 5173) to make requests to your backend (running on port 5000).

### Route Groups
```go
api := s.Engine.Group("/api")
{
    authRoutes := api.Group("/auth")
    {
        authRoutes.POST("/signup", authHandler.Signup)
        authRoutes.POST("/login", authHandler.Login)
    }
}
```

**Route Groups**: Organize related routes under common prefixes:
- All routes start with `/api`
- Auth routes are under `/api/auth`
- This creates endpoints like `POST /api/auth/signup`

---

## 6. Authentication System - JWT Implementation {#authentication}

### What is JWT?
JSON Web Token (JWT) is a way to securely transmit information between parties. It consists of three parts:
1. **Header**: Token type and signing algorithm
2. **Payload**: Claims (user data)
3. **Signature**: Ensures token hasn't been tampered with

**File: `pkg/utils/jwt.go`**

### JWT Claims Structure
```go
type Claims struct {
    UserID primitive.ObjectID `json:"userId"`
    jwt.RegisteredClaims
}
```

**Embedded Struct**: `jwt.RegisteredClaims` is embedded, meaning Claims inherits all its fields (like ExpiresAt, IssuedAt).

### Token Generation
```go
func GenerateToken(userID primitive.ObjectID, c *gin.Context, cfg *config.Config) error {
    expirationTime := time.Now().Add(7 * 24 * time.Hour)
    
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   userID.Hex(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
}
```

**Token Creation Process**:
1. Set expiration time (7 days from now)
2. Create claims with user ID and standard fields
3. Create token with HMAC SHA256 signing
4. Sign token with secret key

### Setting HTTP-Only Cookie
```go
c.SetCookie(
    "jwt",                           // Cookie name
    signedToken,                     // Cookie value
    int(7*24*time.Hour/time.Second), // Max age in seconds
    "/",                             // Path
    "",                              // Domain
    cfg.NodeEnv == "production",     // Secure flag
    true,                            // HttpOnly flag
)
```

**Security Features**:
- `HttpOnly`: Prevents JavaScript access (XSS protection)
- `Secure`: Only sent over HTTPS in production
- `MaxAge`: Cookie expires after 7 days

---

## 7. Authentication Handlers and Endpoints {#auth-handlers}

### Handler Structure
```go
type AuthHandler struct {
    Config          *config.Config
    CloudinaryService *utils.CloudinaryService
}
```

**Dependency Injection**: The handler receives its dependencies (config, services) when created.

**File: `internal/auth/handler.go`**

### Signup Endpoint
```go
func (h *AuthHandler) Signup(c *gin.Context) {
    var req SignupRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "All fields are required"})
        return
    }
}
```

**Request Binding**: Gin automatically parses JSON request body into Go struct:
```go
type SignupRequest struct {
    FullName string `json:"fullName" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

**Validation Tags**: 
- `binding:"required"`: Field must be present
- `binding:"email"`: Must be valid email format
- `binding:"min=6"`: Minimum 6 characters

### Password Hashing
```go
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
```

**bcrypt**: Industry-standard password hashing:
- **Salt**: Random data added to password before hashing
- **Cost**: Number of hashing rounds (higher = more secure but slower)
- **One-way**: Can't reverse the hash to get original password

### Database Operations
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := db.DB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
if err == nil {
    c.JSON(http.StatusBadRequest, gin.H{"message": "Email already exists"})
    return
}
```

**MongoDB Query Pattern**:
1. Create timeout context
2. Query collection with filter
3. Decode result into struct
4. Handle different error types

### Authentication Middleware
**File: `internal/auth/middleware.go`**

```go
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString, err := c.Cookie("jwt")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
            c.Abort()
            return
        }
        
        claims := &utils.Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(cfg.JWTSecret), nil
        })
    }
}
```

**Middleware Pattern**: Function that runs before route handlers:
1. Extract JWT from cookie
2. Validate token signature and expiration
3. Extract user ID from claims
4. Fetch user from database
5. Add user to context for handlers to use
6. Call `c.Next()` to continue to next handler

---

## 8. WebSocket Implementation for Real-Time Communication {#websockets}

### Understanding WebSockets
WebSockets provide full-duplex communication between client and server:
- **HTTP**: Request-response pattern
- **WebSocket**: Persistent connection, both sides can send messages

**File: `pkg/utils/socket.go`**

### WebSocket Upgrader
```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return r.Header.Get("Origin") == "http://localhost:5173"
    },
}
```

**Upgrader**: Converts HTTP connection to WebSocket:
- **Buffer sizes**: Memory allocated for reading/writing
- **CheckOrigin**: Security check to prevent unauthorized connections

### Client Structure
```go
type Client struct {
    Conn *websocket.Conn
    UserID primitive.ObjectID
}
```

**Client**: Represents a single WebSocket connection with associated user.

### Hub Pattern
```go
type Hub struct {
    clients    map[primitive.ObjectID]*Client
    broadcast  chan models.Message
    register   chan *Client
    unregister chan *Client
    mu         sync.Mutex
}
```

**Hub**: Central manager for all WebSocket connections:
- **clients**: Map of user ID to client connection
- **channels**: Go's way of communication between goroutines
- **mutex**: Prevents race conditions when multiple goroutines access the map

### Hub Operations
```go
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client.UserID] = client
            h.mu.Unlock()
            h.sendOnlineUsers()
            
        case client := <-h.unregister:
            h.mu.Lock()
            delete(h.clients, client.UserID)
            h.mu.Unlock()
            h.sendOnlineUsers()
            
        case message := <-h.broadcast:
            h.mu.Lock()
            receiverClient, ok := h.clients[message.ReceiverID]
            h.mu.Unlock()
            
            if ok {
                msgJSON, _ := json.Marshal(message)
                receiverClient.Conn.WriteMessage(websocket.TextMessage, msgJSON)
            }
        }
    }
}
```

**Select Statement**: Go's way of handling multiple channel operations:
- **register**: Add new client to hub
- **unregister**: Remove client from hub
- **broadcast**: Send message to specific client

**Concurrency Safety**: 
- `sync.Mutex` prevents race conditions
- `Lock()` before map access, `Unlock()` after

### WebSocket Handler
```go
func WebSocketHandler(c *gin.Context, hub *Hub) {
    userAny, exists := c.Get("user")
    loggedInUser := userAny.(models.User)
    
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    
    client := &Client{Conn: conn, UserID: loggedInUser.ID}
    hub.register <- client
    
    go func() {
        defer func() {
            hub.unregister <- client
            conn.Close()
        }()
        
        for {
            _, _, err := conn.ReadMessage()
            if err != nil {
                break
            }
        }
    }()
}
```

**Upgrade Process**:
1. Get authenticated user from context
2. Upgrade HTTP connection to WebSocket
3. Create client and register with hub
4. Start goroutine to keep connection alive
5. Unregister client when connection closes

---

## 9. Chat Functionality and Message Handling {#chat-functionality}

**File: `internal/chat/handler.go`**

### Chat Handler Structure
```go
type ChatHandler struct {
    CloudinaryService *utils.CloudinaryService
}
```

### Getting Users for Sidebar
```go
func (h *ChatHandler) GetUsersForSidebar(c *gin.Context) {
    userAny, exists := c.Get("user")
    loggedInUser := userAny.(models.User)
    
    var users []models.User
    usersCollection := db.DB.Collection("users")
    
    cursor, err := usersCollection.Find(ctx, 
        bson.M{"_id": bson.M{"$ne": loggedInUser.ID}}, 
        options.Find().SetProjection(bson.M{"password": 0}))
}
```

**MongoDB Query Breakdown**:
- `bson.M{"_id": bson.M{"$ne": loggedInUser.ID}}`: Find users where ID is not equal to current user
- `SetProjection(bson.M{"password": 0})`: Exclude password field from results
- `cursor.All(ctx, &users)`: Decode all results into slice

### Getting Messages Between Users
```go
func (h *ChatHandler) GetMessages(c *gin.Context) {
    receiverIDParam := c.Param("id")
    receiverID, err := primitive.ObjectIDFromHex(receiverIDParam)
    
    filter := bson.M{
        "$or": []bson.M{
            {"senderId": myID, "receiverId": receiverID},
            {"senderId": receiverID, "receiverId": myID},
        },
    }
    
    findOptions := options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}})
}
```

**Complex Query**:
- `$or`: MongoDB operator for "either condition"
- Gets messages where current user is sender OR receiver
- Sorts by creation time (oldest first)

### Sending Messages
```go
func (h *ChatHandler) SendMessage(c *gin.Context) {
    var req SendMessageRequest
    c.ShouldBindJSON(&req)
    
    if req.Text == "" && req.Image == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Message text or image is required"})
        return
    }
    
    newMessage := models.Message{
        ID:         primitive.NewObjectID(),
        SenderID:   senderID,
        ReceiverID: receiverID,
        Text:       req.Text,
        Image:      imageUrl,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    
    _, err = messagesCollection.InsertOne(ctx, newMessage)
    utils.EmitNewMessage(newMessage)
}
```

**Message Flow**:
1. Parse request body
2. Validate at least text or image provided
3. Upload image to Cloudinary if present
4. Create message document
5. Save to database
6. Emit via WebSocket for real-time delivery

---

## 10. Image Upload Integration with Cloudinary {#image-upload}

### Cloudinary Service
**File: `pkg/utils/cloudinary.go`**

```go
type CloudinaryService struct {
    cld *cloudinary.Cloudinary
}

func NewCloudinaryService(cfg *config.Config) *CloudinaryService {
    cld, _ := cloudinary.NewFromParams(
        cfg.CloudinaryCloudName,
        cfg.CloudinaryAPIKey,
        cfg.CloudinaryAPISecret)
    
    return &CloudinaryService{cld: cld}
}
```

### Image Upload Process
```go
func (cs *CloudinaryService) UploadImage(base64Image string) (string, error) {
    ctx := context.Background()
    
    uploadResult, err := cs.cld.Upload.Upload(ctx, base64Image, uploader.UploadParams{
        Folder: "chat-app",
    })
    
    return uploadResult.SecureURL, nil
}
```

**Base64 Images**: Frontend sends images as base64 encoded strings:
- User selects image file
- JavaScript converts to base64
- Sent in JSON request
- Backend uploads to Cloudinary
- Returns public URL

---

## 11. Frontend-Backend Communication Flow {#frontend-backend}

### API Endpoints Overview

**Authentication Routes** (`/api/auth/`):
- `POST /signup`: Create new user account
- `POST /login`: Authenticate user
- `POST /logout`: Clear authentication cookie
- `PUT /update-profile`: Update user profile picture
- `GET /check`: Verify current authentication status

**Message Routes** (`/api/messages/`):
- `GET /users`: Get all users for sidebar
- `GET /:id`: Get message history with specific user
- `POST /send/:id`: Send message to specific user

**WebSocket Route**:
- `GET /ws`: Upgrade to WebSocket connection

### Request/Response Flow

**Example: Sending a Message**
1. Frontend makes POST request to `/api/messages/send/USER_ID`
2. Auth middleware validates JWT cookie
3. Chat handler processes request
4. Message saved to MongoDB
5. Message emitted via WebSocket to receiver
6. Response sent back to sender
7. Both users see message in real-time

### Frontend Axios Configuration
```javascript
export const axiosInstance = axios.create({
  baseURL: "http://localhost:5001/api",
  withCredentials: true,
});
```

**Key Settings**:
- `baseURL`: All requests prefixed with this URL
- `withCredentials: true`: Include cookies in requests (for JWT)

---

## 12. Application Startup and Flow {#startup-flow}

### Main Function Breakdown
**File: `cmd/api/main.go`**

```go
func main() {
    // 1. Load configuration
    cfg := config.LoadConfig()
    
    // 2. Connect to MongoDB
    db.ConnectDB(cfg)
    defer db.DisconnectDB()
    
    // 3. Initialize WebSocket Hub
    hub := utils.InitWebSocketHub()
    
    // 4. Create and setup server
    appServer := server.NewServer(cfg)
    appServer.SetupRoutes(hub)
    
    // 5. Start server in goroutine
    go func() {
        appServer.Run()
    }()
    
    // 6. Wait for shutdown signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
}
```

### Startup Sequence
1. **Configuration**: Load environment variables
2. **Database**: Establish MongoDB connection
3. **WebSocket Hub**: Start goroutine for managing connections
4. **Server**: Configure Gin with routes and middleware
5. **Graceful Startup**: Server runs in background goroutine
6. **Signal Handling**: Wait for Ctrl+C or termination signal
7. **Graceful Shutdown**: Clean up resources before exit

### Goroutines in Action
- **Main goroutine**: Handles startup and waits for signals
- **Hub goroutine**: Manages WebSocket connections
- **Server goroutine**: Handles HTTP requests
- **Client goroutines**: One per WebSocket connection

---

## 13. Security Considerations and Best Practices {#security}

### Authentication Security
1. **JWT in HTTP-Only Cookies**: Prevents XSS attacks
2. **bcrypt Password Hashing**: Industry standard with salt
3. **Token Expiration**: 7-day expiry limits exposure
4. **Secure Cookies**: Only sent over HTTPS in production

### CORS Configuration
```go
cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowCredentials: true,
}
```
- Restricts which domains can make requests
- Allows credentials (cookies) to be sent

### Input Validation
```go
type SignupRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```
- Server-side validation prevents malicious input
- Email format validation
- Password length requirements

### Database Security
- **Projection**: Exclude sensitive fields like passwords
- **Context Timeouts**: Prevent hanging database operations
- **Error Handling**: Don't expose internal errors to clients

### WebSocket Security
- **Origin Checking**: Verify requests come from expected domain
- **Authentication Required**: Must be logged in to connect
- **Connection Limits**: Hub pattern naturally limits connections

---

## Conclusion

This chat application demonstrates modern Go web development practices:

1. **Clean Architecture**: Separation of concerns with clear layers
2. **Dependency Injection**: Services receive dependencies explicitly
3. **Error Handling**: Explicit error checking throughout
4. **Concurrency**: Goroutines for WebSocket management
5. **Security**: JWT authentication with HTTP-only cookies
6. **Real-time Features**: Custom WebSocket implementation
7. **Third-party Integration**: Cloudinary for image uploads

The combination of Go's simplicity, Gin's flexibility, and MongoDB's document model creates a robust foundation for real-time applications. The WebSocket implementation using the Hub pattern efficiently manages multiple concurrent connections while maintaining type safety and performance.

Key Go concepts demonstrated:
- **Structs and Methods**: Object-oriented patterns
- **Interfaces**: Implicit interface satisfaction
- **Goroutines and Channels**: Concurrent programming
- **Error Handling**: Explicit error checking
- **Package Management**: Module system and imports
- **HTTP Handling**: Web server implementation
- **Database Integration**: MongoDB operations
- **JSON Processing**: Request/response handling

This architecture scales well and follows Go idioms, making it maintainable and extensible for future features.