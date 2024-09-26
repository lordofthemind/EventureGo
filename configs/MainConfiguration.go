package configs

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	// Server Configuration
	ServerPort  int    // The port the server listens on
	Environment string // The environment (e.g., "development", "production")

	// Database Configuration
	PostgresURL  string        // URL for connecting to the PostgreSQL database
	MongoDBURI   string        // URI for connecting to the MongoDB database
	DatabaseType string        // Type of the database being used (e.g., "postgres", "mongodb")
	GormDB       *gorm.DB      // GORM DB object for PostgreSQL
	MongoClient  *mongo.Client // MongoDB client connection

	// Security (TLS)
	EnableTLS   bool   // Flag to enable TLS
	TLSKeyFile  string // Path to the TLS private key file
	TLSCertFile string // Path to the TLS certificate file

	// Token & Authentication
	TokenType           string        // Type of token used (e.g., "Bearer")
	TokenSymmetricKey   string        // Symmetric key for token signing
	TokenExpiryDuration time.Duration // Duration before the token expires

	// CORS Configuration
	CORSAllowedOrigins   []string // List of allowed origins for CORS
	CORSAllowedMethods   []string // List of allowed methods for CORS (e.g., GET, POST)
	CORSAllowedHeaders   []string // List of allowed headers in CORS requests
	CORSExposedHeaders   []string // List of headers exposed to the browser
	CORSAllowCredentials bool     // Whether or not credentials are allowed in CORS requests

)

func LoadMainConfiguration(configFile string) error {
	viper.SetConfigFile(configFile)

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config.yaml file: %w", err)
	}

	PostgresURL = viper.GetString("postgres_url")
	MongoDBURI = viper.GetString("mongodb_uri")
	DatabaseType = viper.GetString("database_type")

	log.Println("Main Configuration Done!!")

	return nil
}
