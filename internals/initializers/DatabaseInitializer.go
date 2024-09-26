package initializers

import (
	"context"
	"log"
	"time"

	"github.com/lordofthemind/EventureGo/configs"
	"github.com/lordofthemind/EventureGo/internals/types"

	"github.com/lordofthemind/mygopher/gophermongo"
	"github.com/lordofthemind/mygopher/gopherpostgres"
)

func DatabaseInitializer() {
	ctx := context.Background()

	if configs.DatabaseType == "postgres" {
		// Initialize PostgreSQL
		gormDB, err := gopherpostgres.ConnectToPostgresGORM(ctx, configs.PostgresURL, 10*time.Second, 3)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL using GORM: %v", err)
		}

		err = gopherpostgres.CheckAndEnableUUIDExtension(gormDB)
		if err != nil {
			log.Fatalf("Failed to confirm UUID extension: %v", err)
		}

		// Auto migrate for GORM (Postgres)
		if err := gormDB.AutoMigrate(&types.SuperUserType{}); err != nil {
			log.Fatalf("Failed to migrate Postgres database: %v", err)
		}

		// Set global GormDB
		configs.GormDB = gormDB
	}

	if configs.DatabaseType == "mongodb" {
		// Initialize MongoDB client
		mongoClient, err := gophermongo.ConnectToMongoDB(ctx, configs.MongoDBURI, 10*time.Second, 3)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Set global MongoClient
		configs.MongoClient = mongoClient
	}
}
