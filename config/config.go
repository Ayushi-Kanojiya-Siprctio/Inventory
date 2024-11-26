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

	err = PG.AutoMigrate(&models.Inventory{}) // Add all the models you want to migrate here
	if err != nil {
		log.Fatal("Error migrating models:", err)
	}

	log.Println("PostgreSQL connected and tables migrated successfully!")
}
