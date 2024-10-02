package cmd

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/lordofthemind/EventureGo/configs"
	"github.com/lordofthemind/EventureGo/internals/handlers"
	"github.com/lordofthemind/EventureGo/internals/initializers"
	"github.com/lordofthemind/EventureGo/internals/middlewares"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/repositories/inmemory"
	"github.com/lordofthemind/EventureGo/internals/repositories/mongodb"
	"github.com/lordofthemind/EventureGo/internals/repositories/postgresdb"
	"github.com/lordofthemind/EventureGo/internals/routes"
	"github.com/lordofthemind/EventureGo/internals/services"
	"github.com/lordofthemind/mygopher/gophergin"
	"github.com/lordofthemind/mygopher/gopherlogger"
	"github.com/lordofthemind/mygopher/gophermongo"
	"github.com/lordofthemind/mygopher/gophersmtp"
	"github.com/lordofthemind/mygopher/gophertoken"
)

func GinServer() {
	// Set up logger
	logFile, err := gopherlogger.SetUpLoggerFile("ginServer.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logFile.Close()

	// Load configuration
	err = configs.LoadMainConfiguration("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration file: %v", err)
	}

	// Initialize the database (Postgres, MongoDB, or in-memory) based on the loaded config
	initializers.DatabaseInitializer()

	// Set up repository and service based on the selected database type from the config
	var superUserRepository repositories.SuperUserRepositoryInterface
	var eventRepository repositories.EventRepositoryInterface

	switch configs.DatabaseType {
	case "inmemory":
		// Initialize In-memory repository
		superUserRepository = inmemory.NewInMemorySuperUserRepository()

	case "postgres":
		// Check if GORM PostgreSQL connection is initialized
		if configs.GormDB == nil {
			log.Fatalf("Postgres connection was not initialized")
		}
		// Initialize PostgreSQL repository
		superUserRepository = postgresdb.NewPostgresSuperUserRepository(configs.GormDB)
		eventRepository = postgresdb.NewPostgresEventRepository(configs.GormDB)

	case "mongodb":
		// Check if MongoDB connection is initialized
		if configs.MongoClient == nil {
			log.Fatalf("MongoDB client was not initialized")
		}
		// Initialize MongoDB repository
		eventureGoDatabase := gophermongo.GetDatabase(configs.MongoClient, "EventureGo")
		superUserRepository = mongodb.NewMongoSuperUserRepository(eventureGoDatabase)
		eventRepository = mongodb.NewMongoEventRepository(eventureGoDatabase)

		// Similarly, if you need to set up another repository with a different database:
		// eventDB := gophermongo.GetDatabase(configs.MongoClient, "events")
		// eventRepository = mongodb.NewMongoEventRepository(eventDB) // Example

	default:
		log.Fatalf("Invalid database configuration: %s", configs.DatabaseType)
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
	eventService := services.NewEventService(eventRepository)

	// Initialize handler
	superUserHandler := handlers.NewSuperUserGinHandler(superUserService)
	eventHandler := handlers.NewEventGinHandler(eventService)

	// Use gophergin to set up the server
	serverConfig := gophergin.ServerConfig{
		Port:        configs.ServerPort,
		UseTLS:      configs.EnableTLS,
		TLSCertFile: configs.TLSCertFile,
		TLSKeyFile:  configs.TLSKeyFile,
		UseCORS:     configs.EnableCors,
		CORSConfig: cors.Config{
			AllowOrigins:     configs.CORSAllowedOrigins,
			AllowMethods:     configs.CORSAllowedMethods,
			AllowHeaders:     configs.CORSAllowedHeaders,
			ExposeHeaders:    configs.CORSExposedHeaders,
			AllowCredentials: configs.CORSAllowCredentials,
			MaxAge:           12 * time.Hour,
		},
	}

	ginServer := gophergin.NewGinServer(&gophergin.ServerSetupImpl{}, serverConfig)

	// Get the router from the server
	router := ginServer.GetRouter()

	// Middleware
	router.Use(middlewares.RequestIDGinMiddleware())

	// Set up routes
	routes.SetupSuperUserGinRoutes(router, superUserHandler, tokenManager)
	routes.SetupEventGinRoutes(router, eventHandler, tokenManager)

	// Start server (with or without TLS)
	if err := ginServer.Start(); err != nil {
		log.Fatalf("Failed to start Gin server on port %d: %v", serverConfig.Port, err)
	}

	// Graceful shutdown handling
	ginServer.GracefulShutdown()
}
