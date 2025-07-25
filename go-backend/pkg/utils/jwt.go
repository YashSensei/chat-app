package utils

import (
	"fmt"        // For formatted error messages
	//"net/http"   // REQUIRED for http.SameSiteStrictMode and other HTTP constants
	"time"       // For token expiration

	"go-backend/config" // Import your config package to get JWT_SECRET. IMPORTANT: Replace "go-backend" with your actual Go module name from go.mod
	"github.com/gin-gonic/gin" // Gin context for setting cookies and responses
	"github.com/golang-jwt/jwt/v5" // JWT library for Go (version 5 is used here)
	"go.mongodb.org/mongo-driver/bson/primitive" // For handling ObjectID from user ID
)

// Claims defines the structure of our JWT claims.
// It embeds jwt.RegisteredClaims for standard JWT fields like Issuer, ExpiresAt, etc.
// UserID is a custom claim to store the user's MongoDB ObjectID.
type Claims struct {
	UserID primitive.ObjectID `json:"userId"` // Custom claim to store the user's ID
	jwt.RegisteredClaims     // Standard JWT claims (e.g., expiration, issued at, subject)
}

// GenerateToken creates a JWT and sets it as an HTTP-only cookie.
// This function mirrors your `generateToken` in Node.js.

// Parameters:
//   userID: The MongoDB ObjectID of the user for whom the token is being generated.
//   c: The Gin context, used to set the HTTP cookie in the response.
//   cfg: A pointer to the application's configuration, containing the JWT secret.

// Returns: An error if token generation or cookie setting fails, otherwise nil.
func GenerateToken(userID primitive.ObjectID, c *gin.Context, cfg *config.Config) error {
	// Define the expiration time for the token (7 days from now).
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// Create the JWT claims payload.
	// The `UserID` field of our custom `Claims` struct is populated with the provided `userID`.
	// `jwt.RegisteredClaims` are populated with standard JWT fields:
	//   - `ExpiresAt`: The time when the token becomes invalid. `jwt.NewNumericDate` converts `time.Time` to a numeric date.
	//   - `IssuedAt`: The time when the token was created.
	//   - `Subject`: A unique identifier for the subject of the token. Here, we use the hex string of the `userID`.
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.Hex(), // Use the hex string representation of the ObjectID
		},
	}

	// Create the token using the HS256 signing method (HMAC SHA256) and the defined claims.
	// HS256 is a symmetric algorithm, meaning the same secret key is used for both signing and verifying the token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your JWT secret key.
	// The secret key is retrieved from your application configuration (`cfg.JWTSecret`).
	// It must be converted to a byte slice `[]byte()`.
	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {  
		// If signing fails (e.g., secret key is invalid), return a wrapped error.
		return fmt.Errorf("failed to sign token: %w", err)
	}

	// Set the JWT as an HTTP-only cookie in the Gin context's response.
	// This is critical for security:
	//   - `httpOnly: true` prevents client-side JavaScript from accessing the cookie, mitigating XSS attacks.
	//
	// Parameters for `c.SetCookie`:
	//   - `name`: "jwt" (This must match the cookie name your frontend expects).
	//   - `value`: The `signedToken` string.
	//   - `maxAge`: The maximum age of the cookie in seconds. We convert 7 days duration to seconds.
	//   - `path`: "/" (The cookie is valid for all paths on the domain).
	//   - `domain`: "" (An empty string means the cookie is valid for the current host only).
	//   - `secure`: `cfg.NodeEnv == "production"` (The `Secure` flag ensures the cookie is only sent over HTTPS.
	//     It's `true` in production, `false` in development for easier local testing).
	//   - `httpOnly`: `true` (Makes the cookie inaccessible to JavaScript).
	//   - `sameSite`: `http.SameSiteStrictMode` (The strictest SameSite policy).
	// CORRECTED: Removed http.SameSiteStrictMode as it's not accepted by this Gin SetCookie signature.
	c.SetCookie(
		"jwt",
		signedToken,
		int(7*24*time.Hour/time.Second), // Convert 7 days duration to seconds
		"/",
		"",
		cfg.NodeEnv == "production", // Secure flag: true if in production, false otherwise
		true,                        // HttpOnly flag: true
		// http.SameSiteStrictMode,     // COMMENTED OUT: SameSite flag. This argument is causing the error.
	)

	return nil // Return nil if token generation and cookie setting were successful
}
