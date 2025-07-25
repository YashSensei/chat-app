package auth

import (
	"context"    // For context with MongoDB operations (e.g., timeouts)
	"fmt"        // For formatted error messages
	"net/http"   // For HTTP status codes (e.g., 401 Unauthorized, 404 Not Found)
	"strings"    // For string manipulation (e.g., checking if an error message contains "token is expired")
	"time"       // For time-related operations (e.g., checking token expiration)

	"go-backend/config" // Import your config package to access JWT_SECRET. IMPORTANT: Replace "go-backend" with your actual Go module name from go.mod
	"go-backend/internal/models" // Import models to use the User struct for database operations
	"go-backend/pkg/db" // Import db to access the global MongoDB client (db.DB)
	"go-backend/pkg/utils" // Import utils for the JWT Claims struct (defined in jwt.go)

	"github.com/gin-gonic/gin" // Gin context for handling HTTP requests and responses
	"github.com/golang-jwt/jwt/v5" // The JWT library for Go (version 5 is used here)
	"go.mongodb.org/mongo-driver/bson" // For constructing MongoDB queries (e.g., bson.M for map-like queries)
	//"go.mongodb.org/mongo-driver/bson/primitive" // For converting string IDs to MongoDB's ObjectID type
	"go.mongodb.org/mongo-driver/mongo" // The main MongoDB client type, used to check for specific errors like ErrNoDocuments
)

// AuthMiddleware creates a Gin middleware to protect routes.
// It performs the following steps:
// 1. Retrieves the JWT token from the "jwt" HTTP-only cookie.
// 2. Parses and validates the token's signature and expiration using the configured JWT secret.
// 3. Extracts the UserID from the token's claims.
// 4. Queries the MongoDB database to find the user corresponding to the UserID.
// 5. If the token is valid and the user is found, it attaches the user object to the Gin context.
// 6. Calls the next handler in the Gin chain.
// If any step fails (e.g., no token, invalid token, user not found), it aborts the request
// and sends an appropriate JSON error response.
// This function directly mirrors the functionality of your `protectRoute` middleware in Node.js.
//
// Parameters:
//   cfg: A pointer to the application's `Config` struct, which contains the `JWTSecret` needed for token validation.
//
// Returns:
//   A `gin.HandlerFunc`, which is the standard type for Gin middleware functions.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	// The returned function is the actual middleware that Gin will execute for protected routes.
	return func(c *gin.Context) {
		// 1. Get the JWT token string from the "jwt" cookie.
		// `c.Cookie("jwt")` attempts to read the cookie by its name.
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			// If the "jwt" cookie is not found (meaning no token was provided),
			// send a 401 Unauthorized response and abort the request.
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized - No Token Provided"})
			c.Abort() // Stop processing this request and don't call subsequent handlers
			return
		}

		// Initialize a new `utils.Claims` struct. This struct will be populated
		// with the claims extracted from the JWT after parsing.
		claims := &utils.Claims{}

		// Parse the token string using `jwt.ParseWithClaims`.
		// This function performs several critical steps:
		//   - Decodes the token string.
		//   - Validates its signature using the provided secret key.
		//   - Unmarshals the token's payload (claims) into the `claims` struct.
		// The `func(token *jwt.Token) (interface{}, error)` is a callback function
		// that provides the secret key used for signature verification.
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// A security check: ensure the signing method used in the token's header
			// is the expected HMAC SHA256 (`jwt.SigningMethodHS256`).
			// This prevents attackers from changing the algorithm to a weaker one.
			// CORRECTED LINE: Directly compare the method with the expected signing method constant.
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the JWT secret key (from your config) as a byte slice for verification.
			return []byte(cfg.JWTSecret), nil
		})

		// Handle any errors that occurred during token parsing or validation.
		if err != nil {
			// Differentiate between common JWT errors for more specific messages.
			if err == jwt.ErrSignatureInvalid {
				// If the token's signature is invalid (e.g., tampered or wrong secret).
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized - Invalid Token Signature"})
			} else if strings.Contains(err.Error(), "token is expired") {
				// If the token has expired. The `jwt.ParseWithClaims` will automatically check `exp`.
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized - Token Expired"})
			} else {
				// Catch-all for other parsing/validation errors.
				c.JSON(http.StatusUnauthorized, gin.H{"message": fmt.Sprintf("Unauthorized - Invalid Token: %v", err)})
			}
			c.Abort() // Abort the request
			return
		}

		// After parsing, explicitly check if the token is considered valid by the JWT library.
		// This checks overall validity including expiration (if not caught by string check above)
		// and other registered claims.
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized - Invalid Token"})
			c.Abort()
			return
		}

		// Although `jwt.ParseWithClaims` often handles expiration, an explicit check
		// provides clarity and can be useful for debugging or specific logic.
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized - Token Expired"})
			c.Abort()
			return
		}

		// 2. Find the user in the database using the UserID extracted from the claims.
		// The UserID from claims is already a `primitive.ObjectID`.
		userID := claims.UserID

		// Get a reference to the "users" collection in your MongoDB database.
		usersCollection := db.DB.Collection("users")

		var user models.User // Declare a variable of type `models.User` to hold the retrieved user data.

		// Create a context with a timeout for the database query.
		// This prevents the application from hanging indefinitely if the database is slow.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() // Ensure the context resources are released when the function exits.

		// Execute the MongoDB query: Find one document in the "users" collection
		// where the "_id" field matches the `userID` from the token claims.
		// `bson.M` is a convenient type for creating BSON documents (maps) for queries.
		// `.Decode(&user)` attempts to unmarshal the found MongoDB document into our `user` struct.
		err = usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			// Handle specific MongoDB errors.
			if err == mongo.ErrNoDocuments {
				// If no document was found for the given ID, even if the token was valid.
				c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			} else {
				// Catch-all for other database errors (e.g., connection issues).
				c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Internal server error fetching user: %v", err)})
			}
			c.Abort() // Abort the request if the user cannot be found or there's a DB error.
			return
		}

		// 3. If everything is successful (token valid, user found), attach the `user` object
		// to the Gin context. This makes the authenticated user's information easily
		// accessible to subsequent handlers in the request chain (e.g., controllers).
		// The key "user" is used to retrieve it later: `c.Get("user")`.
		c.Set("user", user)

		// Call the next handler in the Gin chain. If there are other middlewares, they run next.
		// If not, the final route handler will be executed.
		c.Next()
	}
}