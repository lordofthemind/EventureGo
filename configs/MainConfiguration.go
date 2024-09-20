package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	Port         int
	PostgresURL  string
	MongodbURI   string
	GormDB       *gorm.DB
	MongoDB      *mongo.Database
	MongoClient  *mongo.Client
	DatabaseType string
)

func LoadMainConfiguration(configFile string) error {
	viper.SetConfigFile(configFile)

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config.yaml file: %w", err)
	}

	PostgresURL = viper.GetString("postgres_url")
	MongodbURI = viper.GetString("mongodb_uri")
	DatabaseType = viper.GetString("database_type")

	log.Println("Main Configuration Done!!")

	return nil
}
