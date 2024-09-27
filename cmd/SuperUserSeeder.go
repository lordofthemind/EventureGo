package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lordofthemind/EventureGo/configs"
	"github.com/lordofthemind/EventureGo/internals/initializers"
	"github.com/lordofthemind/EventureGo/internals/repositories"
	"github.com/lordofthemind/EventureGo/internals/repositories/inmemory"
	"github.com/lordofthemind/EventureGo/internals/repositories/mongodb"
	"github.com/lordofthemind/EventureGo/internals/repositories/postgresdb"
	"github.com/lordofthemind/EventureGo/internals/services"
	"github.com/lordofthemind/EventureGo/internals/utils"
	"github.com/lordofthemind/mygopher/gopherlogger"
	"github.com/lordofthemind/mygopher/gophermongo"
	"github.com/lordofthemind/mygopher/gophertoken"
)

func SuperUserSeeder() {
	// Set up logger
	logFile, err := gopherlogger.SetUpLoggerFile("superUserSeeder.log")
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

	default:
		log.Fatalf("Invalid database configuration: %s", configs.DatabaseType)
	}

	// Use the new NewTokenManager function
	tokenManager, err := gophertoken.NewTokenManager(configs.TokenType, configs.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("Failed to initiate token: %v", err)
	}

	// Initialize service and handler
	superUserService := services.NewSuperUserService(superUserRepository, tokenManager)

	// Seed SuperUsers before starting the server
	seedSuperUsers(superUserService)

}

func seedSuperUsers(service services.SuperUserServiceInterface) {
	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	fmt.Println("Starting SuperUser Seeder...")

	for {
		req := &utils.RegisterSuperuserRequest{}

		// Get Email
		fmt.Print("Enter SuperUser Email: ")
		email, _ := reader.ReadString('\n')
		req.Email = strings.TrimSpace(email)

		// Get Full Name
		fmt.Print("Enter SuperUser Full Name: ")
		fullName, _ := reader.ReadString('\n')
		req.FullName = strings.TrimSpace(fullName)

		// Get Username
		fmt.Print("Enter SuperUser Username: ")
		username, _ := reader.ReadString('\n')
		req.Username = strings.TrimSpace(username)

		// Get Password
		fmt.Print("Enter SuperUser Password: ")
		password, _ := reader.ReadString('\n')
		req.Password = strings.TrimSpace(password)

		// Seed the SuperUser
		err := service.SeedSuperUser(ctx, req)
		if err != nil {
			fmt.Printf("Failed to seed superuser: %v\n", err)
		} else {
			fmt.Println("SuperUser seeded successfully.")
		}

		// Ask if user wants to seed another
		fmt.Print("Do you want to seed another SuperUser? (y/n): ")
		another, _ := reader.ReadString('\n')
		another = strings.TrimSpace(strings.ToLower(another))

		if another != "y" {
			break
		}
	}

	fmt.Println("Seeding process completed.")
}
