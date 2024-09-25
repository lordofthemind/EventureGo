package cmd

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/lordofthemind/EventureGo/configs"
	"github.com/lordofthemind/EventureGo/internals/handlers"
	"github.com/lordofthemind/EventureGo/internals/initializers"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/repositories/inmemory"
	"github.com/lordofthemind/EventureGo/internals/repositories/mongodb"
	"github.com/lordofthemind/EventureGo/internals/repositories/postgresdb"
	"github.com/lordofthemind/EventureGo/internals/routes"
	"github.com/lordofthemind/EventureGo/internals/services"
	"github.com/lordofthemind/mygopher/gophermongo"
	"github.com/lordofthemind/mygopher/mygopherlogger"
)

func FiberServer() {
	// Set up logger
	logFile, err := mygopherlogger.SetUpLoggerFile("fiberServer.log")
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
	case "postgres":
		if configs.GormDB == nil {
			log.Fatalf("Postgres connection was not initialized")
		}
		// Initialize Postgres repository (not shown in your example, but you can add it here)
		superUserRepository = postgresdb.NewPostgresSuperUserRepository(configs.GormDB) // Example
	case "inmemory":

		// Initialize Postgres repository (not shown in your example, but you can add it here)
		superUserRepository = inmemory.NewInMemorySuperUserRepository()

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

	// Initialize service and handler
	superUserService := services.NewSuperUserService(superUserRepository)
	superUserHandler := handlers.NewSuperUserFiberHandler(superUserService)

	// Set up Fiber routes
	app := fiber.New()

	routes.SetupSuperUserFiberRoutes(app, superUserHandler)

	// Start the Fiber server
	serverAddress := ":9090" // This can be configurable
	if err := app.Listen(serverAddress); err != nil {
		log.Fatalf("Failed to start Fiber server: %v", err)
	}
}
