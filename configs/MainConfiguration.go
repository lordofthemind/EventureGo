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
	ServerHost  string //
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

	// Logging
	LoggingLevel string

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

	// Fetch environment and database type from the config file
	Environment = viper.GetString("application.environment")
	DatabaseType = viper.GetString("application.database_type")

	// Load environment-specific configurations
	viper.Set("active_environment", Environment)
	viper.Set("active_database", DatabaseType)

	// Conditional loading based on environment and database type
	switch Environment {
	case "development", "testing", "production", "staging":
		PostgresURL = viper.GetString(fmt.Sprintf("environments.%s.database.postgres.url", Environment))
		MongoDBURI = viper.GetString(fmt.Sprintf("environments.%s.database.mongodb.uri", Environment))
		CORSAllowedOrigins = viper.GetStringSlice(fmt.Sprintf("environments.%s.cors.allowed_origins", Environment))
		CORSAllowedMethods = viper.GetStringSlice(fmt.Sprintf("environments.%s.cors.allowed_methods", Environment))
		CORSAllowedHeaders = viper.GetStringSlice(fmt.Sprintf("environments.%s.cors.allowed_headers", Environment))
		CORSExposedHeaders = viper.GetStringSlice(fmt.Sprintf("environments.%s.cors.exposed_headers", Environment))
		CORSAllowCredentials = viper.GetBool(fmt.Sprintf("environments.%s.cors.allow_credentials", Environment))
		TLSCertFile = viper.GetString(fmt.Sprintf("environments.%s.cert_file", Environment))
		TLSKeyFile = viper.GetString(fmt.Sprintf("environments.%s.key_file", Environment))
	default:
		return fmt.Errorf("unknown environment: %s", Environment)
	}

	// Conditional loading based on database type
	switch DatabaseType {
	case "postgres":
		PostgresURL = viper.GetString(fmt.Sprintf("environments.%s.database.postgres.url", Environment))
	case "mongodb":
		MongoDBURI = viper.GetString(fmt.Sprintf("environments.%s.database.mongodb.uri", Environment))
	case "inmemory":
		log.Println("Using in-memory database")
	default:
		return fmt.Errorf("unknown database type: %s", DatabaseType)
	}

	ServerHost = viper.GetString("server.host")
	ServerPort = viper.GetInt("server.port")
	EnableTLS = viper.GetBool("server.use_tls")

	LoggingLevel = viper.GetString("logging.level")

	TokenType = viper.GetString("token.type")
	TokenSymmetricKey = viper.GetString("token.symmetric_key")
	TokenExpiryDuration = viper.GetDuration("token.access_duration")

	log.Println("Configuration loaded for environment:", Environment)
	log.Println("Configuration loaded for database:", DatabaseType)
	return nil
}
