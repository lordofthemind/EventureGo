package cmd

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
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
	"github.com/lordofthemind/mygopher/gopherlogger"
	"github.com/lordofthemind/mygopher/gophermongo"
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

	// Initialize service and handler
	superUserService := services.NewSuperUserService(superUserRepository, tokenManager)
	superUserHandler := handlers.NewSuperUserFiberHandler(superUserService, tokenManager)

	// Set up Fiber routes
	app := fiber.New()

	// Apply middleware globally or for specific routes
	app.Use(middlewares.RequestIDFiberMiddleware())

	routes.SetupSuperUserFiberRoutes(app, superUserHandler, tokenManager)

	// Dynamically fetch server address and port from the configuration
	serverAddress := fmt.Sprintf(":%d", configs.ServerPort)
	if err := app.Listen(serverAddress); err != nil {
		log.Fatalf("Failed to start Fiber server: %v", err)
	}
}
