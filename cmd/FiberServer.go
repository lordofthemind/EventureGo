package cmd

import (
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lordofthemind/EventureGo/configs"
	"github.com/lordofthemind/EventureGo/internals/handlers"
	"github.com/lordofthemind/EventureGo/internals/initializers"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/repositories/inmemory"
	"github.com/lordofthemind/EventureGo/internals/repositories/mongodb"
	"github.com/lordofthemind/EventureGo/internals/repositories/postgresdb"
	"github.com/lordofthemind/EventureGo/internals/routes"
	"github.com/lordofthemind/EventureGo/internals/services"
	"github.com/lordofthemind/mygopher/gopherfiber"
	"github.com/lordofthemind/mygopher/gopherlogger"
	"github.com/lordofthemind/mygopher/gophermongo"
	"github.com/lordofthemind/mygopher/gophersmtp"
	"github.com/lordofthemind/mygopher/gophertoken"
)

func FiberServer() {
	// Set up logger
	logFile, err := gopherlogger.SetUpLoggerFile("fiberServer.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logFile.Close()

	// Load configuration
	err = configs.LoadMainConfiguration("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration file: %v", err)
	}

	// Initialize database (Postgres or MongoDB)
	initializers.DatabaseInitializer()

	// Setup repository and service based on the selected database
	var superUserRepository repositories.SuperUserRepositoryInterface

	switch configs.DatabaseType {
	case "inmemory":
		// Initialize Postgres repository (not shown in your example, but you can add it here)
		superUserRepository = inmemory.NewInMemorySuperUserRepository()

	case "postgres":
		// Check if GORM PostgreSQL connection is initialized
		if configs.GormDB == nil {
			log.Fatalf("Postgres connection was not initialized")
		}
		// Initialize Postgres repository (not shown in your example, but you can add it here)
		superUserRepository = postgresdb.NewPostgresSuperUserRepository(configs.GormDB) // Example

	case "mongodb":
		if configs.MongoClient == nil {
			log.Fatalf("MongoDB client was not initialized")
		}

		// Retrieve the specific database for the SuperUser repository
		superUserDB := gophermongo.GetDatabase(configs.MongoClient, "superuser")
		superUserRepository = mongodb.NewMongoSuperUserRepository(superUserDB)

		// Similarly, if you need to set up another repository with a different database:
		// eventDB := gophermongo.GetDatabase(configs.MongoClient, "events")
		// eventRepository = mongodb.NewMongoEventRepository(eventDB) // Example

	default:
		log.Fatalf("Invalid database configuration")
	}

	// Use the new NewTokenManager function
	tokenManager, err := gophertoken.NewTokenManager(configs.TokenType, configs.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("Failed to initiate token: %v", err)
	}

	// Initialize services
	emailService := gophersmtp.NewEmailService(
		configs.SMTPHost,
		configs.SMTPPort,
		configs.EmailUsername,
		configs.EmailPassword,
	)
	superUserService := services.NewSuperUserService(superUserRepository, tokenManager, emailService)

	superUserHandler := handlers.NewSuperUserFiberHandler(superUserService)

	// Create ServerConfig for gopherfiber
	serverConfig := gopherfiber.ServerConfig{
		Port: configs.ServerPort,
		// StaticPath:   configs.StaticPath, // Adjust paths as necessary
		// TemplatePath: configs.TemplatePath,
		UseTLS:      configs.EnableTLS,   // Set true if you want to use TLS
		TLSCertFile: configs.TLSCertFile, // TLS certificate file path
		TLSKeyFile:  configs.TLSKeyFile,  // TLS key file path
		UseCORS:     configs.EnableCors,
		CORSConfig: cors.Config{
			AllowOrigins: configs.CORSAllowedOrigins[0], // Modify as needed for your application
			// AllowMethods:     configs.CORSAllowedMethods[0],
			// AllowHeaders:     configs.CORSAllowedHeaders[0],
			// ExposeHeaders:    configs.CORSExposedHeaders[0],
			// AllowCredentials: configs.CORSAllowCredentials,
			// MaxAge:           12 * time.Hour,
		},
	}

	// Create a new Fiber server using gopherfiber
	fiberServer := gopherfiber.NewFiberServer(&gopherfiber.ServerSetupImpl{}, serverConfig)

	// Set up Fiber routes
	routes.SetupSuperUserFiberRoutes(fiberServer.GetRouter(), superUserHandler, tokenManager)

	// Start the Fiber server
	if err := fiberServer.Start(); err != nil {
		log.Fatalf("Failed to start Fiber server: %v", err)
	}

	// Graceful shutdown on interrupt signal
	fiberServer.GracefulShutdown()
}
