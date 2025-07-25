package server

import (
	"fmt"      // For formatted output (e.g., server start message)
	"log"      // For logging errors
	//"net/http" // For HTTP status codes and constants (e.g., http.StatusUnauthorized)
	"time"     // For time-related operations (e.g., MaxAge duration)

	"go-backend/config" // Import your config package for application settings
	"go-backend/internal/auth" // Import auth package for handlers and middleware
	"go-backend/internal/chat" // Import chat package for handlers
	"go-backend/pkg/utils" // Import utils for CloudinaryService and Hub

	"github.com/gin-contrib/cors" // Gin middleware for CORS
	"github.com/gin-gonic/gin"    // The Gin web framework
)

// Server struct holds the Gin engine and application configuration.
// This allows us to pass dependencies (like config) to the server.
type Server struct {
	Engine *gin.Engine
	Config *config.Config
}

// NewServer creates and initializes a new Gin server instance.
// It sets up the Gin mode (release/debug) and returns a pointer to the Server struct.
func NewServer(cfg *config.Config) *Server {
	// Set Gin mode based on NodeEnv from config.
	// In production, Gin runs in release mode, which disables debug output.
	if cfg.NodeEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize the Gin default engine.
	// Default includes Logger and Recovery middleware.
	engine := gin.Default()

	return &Server{
		Engine: engine,
		Config: cfg,
	}
}

// SetupRoutes configures all API endpoints and applies middleware.
// MODIFIED: Accepts the WebSocket Hub instance.
func (s *Server) SetupRoutes(hub *utils.Hub) {
	// Configure CORS middleware.
	s.Engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize Cloudinary Service.
	cloudinaryService := utils.NewCloudinaryService(s.Config)

	// Initialize authentication and chat handlers.
	authHandler := auth.NewAuthHandler(s.Config, cloudinaryService)
	chatHandler := chat.NewChatHandler(cloudinaryService)

	// Group API routes under "/api".
	api := s.Engine.Group("/api")
	{
		// Authentication Routes (no protection needed for signup/login)
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/signup", authHandler.Signup)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/logout", authHandler.Logout)

			// Protected Auth Routes (require authentication middleware)
			protectedAuthRoutes := authRoutes.Group("/")
			protectedAuthRoutes.Use(auth.AuthMiddleware(s.Config))
			{
				protectedAuthRoutes.PUT("/update-profile", authHandler.UpdateProfile)
				protectedAuthRoutes.GET("/check", authHandler.CheckAuth)
			}
		}

		// Message Routes (all protected)
		messageRoutes := api.Group("/messages")
		messageRoutes.Use(auth.AuthMiddleware(s.Config))
		{
			messageRoutes.GET("/users", chatHandler.GetUsersForSidebar)
			messageRoutes.GET("/:id", chatHandler.GetMessages)
			messageRoutes.POST("/send/:id", chatHandler.SendMessage)
		}
	}

	// WebSocket Route
	// This route will handle upgrading the HTTP connection to a WebSocket.
	// It uses the AuthMiddleware to ensure only authenticated users can establish a WebSocket connection.
	s.Engine.GET("/ws", auth.AuthMiddleware(s.Config), func(c *gin.Context) {
		utils.WebSocketHandler(c, hub) // Pass the hub to the WebSocket handler
	})

	// Serve static files for frontend in production.
	if s.Config.NodeEnv == "production" {
		s.Engine.Static("/static", "./frontend/dist/assets")
		s.Engine.StaticFile("/", "./frontend/dist/index.html")
		s.Engine.NoRoute(func(c *gin.Context) {
			c.File("./frontend/dist/index.html")
		})
	}
}

// Run starts the Gin HTTP server.
func (s *Server) Run() {
	port := s.Config.Port
	if port == "" {
		port = "5000" // Default port if not set in config
	}
	log.Printf("Server is running on PORT: %s", port)
	log.Fatal(s.Engine.Run(fmt.Sprintf(":%s", port)))
}