package config

import (
	"context"
	"fmt"
	"log"

	"main/models"

	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MongoConfig struct {
	MONGO_URI         string `env:"MONGO_URL" envDefault:"mongodb://localhost:27017"`
	MongoDBCollection string `env:"MONGODB_COLLECTION" envDefault:"inventories"`
	MongoDBName       string `env:"MONGODB_DB_NAME" envDefault:"inventoryDB"`
	MongoPort string `env:"MONGO_PORT" envDefault:"8080"`
}

var MongoClient *mongo.Client
var InventoryCollection *mongo.Collection

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading in .env file.....")
	}
}

func InitMongoDB() {
	var config MongoConfig
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(config.MONGO_URI))
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	InventoryCollection = MongoClient.Database(config.MongoDBName).Collection(config.MongoDBCollection)

	log.Println("MongoDB initialized successfully")
}

var PG *gorm.DB

func PostgresConnect() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")

	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresHost, postgresPort, postgresUser, postgresPassword, postgresDB)

	PGDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening connection to PostgreSQL:", err)
	}

	PG = PGDB

	err = PG.AutoMigrate(&models.Inventory{}) 
	if err != nil {
		log.Fatal("Error migrating models:", err)
	}

	log.Println("PostgreSQL connected and tables migrated successfully!")
}
