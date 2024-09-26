package cmd

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
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
	"github.com/lordofthemind/mygopher/gophertoken"
	"github.com/lordofthemind/mygopher/mygopherlogger"
)

func GinServer() {
	// Set up logger
	logFile, err := mygopherlogger.SetUpLoggerFile("ginServer.log")
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

	switch configs.DatabaseType {
	case "postgres":
		// Check if GORM PostgreSQL connection is initialized
		if configs.GormDB == nil {
			log.Fatalf("Postgres connection was not initialized")
		}
		// Initialize PostgreSQL repository
		superUserRepository = postgresdb.NewPostgresSuperUserRepository(configs.GormDB)

	case "mongodb":
		// Check if MongoDB connection is initialized
		if configs.MongoClient == nil {
			log.Fatalf("MongoDB client was not initialized")
		}
		// Initialize MongoDB repository
		superUserDB := gophermongo.GetDatabase(configs.MongoClient, "superuser")
		superUserRepository = mongodb.NewMongoSuperUserRepository(superUserDB)

		// Similarly, if you need to set up another repository with a different database:
		// eventDB := gophermongo.GetDatabase(configs.MongoClient, "events")
		// eventRepository = mongodb.NewMongoEventRepository(eventDB) // Example

	case "inmemory":
		// Initialize In-memory repository
		superUserRepository = inmemory.NewInMemorySuperUserRepository()

	default:
		log.Fatalf("Invalid database configuration: %s", configs.DatabaseType)
	}

	// Use the new NewTokenManager function
	tokenManager, err := gophertoken.NewTokenManager(configs.TokenType, configs.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("Failed to initiate token: %v", err)
	}

	// Initialize service and handler
	superUserService := services.NewSuperUserService(superUserRepository)
	superUserHandler := handlers.NewSuperUserGinHandler(superUserService, tokenManager)

	// Set up Gin routes
	router := gin.Default()
	routes.SetupSuperUserGinRoutes(router, superUserHandler, tokenManager)

	// Dynamically fetch server address and port from the configuration
	serverAddress := fmt.Sprintf(":%d", configs.ServerPort)
	if err := router.Run(serverAddress); err != nil {
		log.Fatalf("Failed to start Gin server on %s: %v", serverAddress, err)
	}
}
