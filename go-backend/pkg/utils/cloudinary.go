package utils

import (
	"context" // For context with Cloudinary upload operations
	"fmt"     // For formatted error messages
	"log"     // For logging errors
	"time"    // For time-related operations (REQUIRED for context.WithTimeout)

	"go-backend/config" // Import your config package for Cloudinary credentials

	"github.com/cloudinary/cloudinary-go/v2" // The Cloudinary Go SDK
	"github.com/cloudinary/cloudinary-go/v2/api/uploader" // For upload specific functions
)

// CloudinaryService struct holds the Cloudinary client instance.
// This allows for dependency injection and easier testing.
type CloudinaryService struct {
	Client *cloudinary.Cloudinary
}

// NewCloudinaryService initializes and returns a new CloudinaryService.
// It takes the application configuration to get Cloudinary API credentials.
func NewCloudinaryService(cfg *config.Config) *CloudinaryService {
	// Create a new Cloudinary client instance using the credentials from your config.
	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)
	if err != nil {
		// If initialization fails, log a fatal error and exit the application,
		// as Cloudinary is a critical dependency for image handling.
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	return &CloudinaryService{Client: cld}
}

// UploadImage uploads a base64 encoded image string to Cloudinary.
// Mirrors backend/src/lib/cloudinary.js's upload functionality.
//
// Parameters:
//   base64Image: The base64 encoded image string (e.g., "data:image/jpeg;base64,...").
//
// Returns:
//   The secure URL of the uploaded image, or an error if the upload fails.
func (cs *CloudinaryService) UploadImage(base64Image string) (string, error) {
	// REVERTED TO RECOMMENDED APPROACH:
	// Create a context with a timeout for the upload operation.
	// This is good practice to prevent the application from hanging indefinitely
	// if the external API (Cloudinary) is slow or unresponsive.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Ensure the context is cancelled when the function exits

	// Define upload parameters.
	// `Folder`: Optional, but good for organizing uploads (e.g., "chat_app_profile_pics").
	// `PublicID`: Cloudinary will generate a unique public ID if not specified.
	// `ResourceType`: "image" is standard for image uploads.
	uploadParams := uploader.UploadParams{
		Folder: "chat_app_images", // You can customize this folder name
	}

	// Perform the upload.
	// The `base64Image` string is directly passed as the source.
	uploadResult, err := cs.Client.Upload.Upload(ctx, base64Image, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to Cloudinary: %w", err)
	}

	// Return the secure URL of the uploaded image.
	return uploadResult.SecureURL, nil
}