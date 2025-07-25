package config
import(
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds all application configurations
type Config struct{
	Port string
	MongoDBURI           string
	JWTSecret            string
	CloudinaryCloudName  string
	CloudinaryAPIKey     string
	CloudinaryAPISecret  string
	NodeEnv              string
}

// LoadConfig reads environment variables and returns a Config struct   
func LoadConfig() *Config{
	// Load .env file. It returns an error if the file doesn't exist,
	// but we log it as info because in production, env vars might be set directly.
	err := godotenv.Load()
	if err != nil{
		log.Println("No .env file found, assuming environment variables are set directly in the environment.")
	}
	return &Config{
		Port:                 getEnv("PORT", "5000"), // Default to 5000 if not set
		MongoDBURI:           getEnv("MONGODB_URI", "mongodb://localhost:27017/chat-app"), // Default URI
		JWTSecret:            getEnv("JWT_SECRET", "supersecretjwtkeyforlocaldevonly"), // IMPORTANT: Change this default in production, better to ensure it's always set in .env
		CloudinaryCloudName:  getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:     getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret:  getEnv("CLOUDINARY_API_SECRET", ""),
		NodeEnv:              getEnv("NODE_ENV", "development"),
	}
}
// Helper function to get environment variable with a fallback default value
func getEnv(key string , defaultvalue string) string{
	if value, exists := os.LookupEnv(key); exists{
		return value
	}
	return defaultvalue
}