package config

import (
	"context"
	"fmt"
	"log"
	"main/models"

	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var MongoClient *mongo.Client
var InventoryCollection *mongo.Collection

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading in .env file.....")
	}
}

// func ConnectDatabase() (*mongo.Client, error) {
// 	LoadEnv()
// 	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Ping the database
// 	if err := client.Ping(context.TODO(), nil); err != nil {
// 		return nil, fmt.Errorf("could not connect to MongoDB: %w", err)
// 	}

// 	log.Println("Connected to MongoDB!")
// 	return client, nil
// }

func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	InventoryCollection = MongoClient.Database("inventoryDB").Collection("inventories")
}

var PG *gorm.DB

func PostgresConnect() {
	// Load environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve PostgreSQL credentials from environment variables
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")

	// Construct the PostgreSQL connection string
	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresHost, postgresPort, postgresUser, postgresPassword, postgresDB)

	// Open a connection to the PostgreSQL database
	PGDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening connection to PostgreSQL:", err)
	}

	// Assign the DB connection to the global PG variable
	PG = PGDB

	err = PG.AutoMigrate(&models.Inventory{}) // Add all the models you want to migrate here
	if err != nil {
		log.Fatal("Error migrating models:", err)
	}

	log.Println("PostgreSQL connected and tables migrated successfully!")
}
