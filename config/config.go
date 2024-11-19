package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
    _ "github.com/lib/pq" 
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var InventoryCollection *mongo.Collection

func InitMongoDB() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatalf("MongoDB connection failed: %v", err)
    }

    InventoryCollection = MongoClient.Database("inventory").Collection("inventories")
}

var PG *sql.DB

func PostgresConnect() {
    // Load environment variables from .env file
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Retrieve PostgreSQL connection details from environment variables
    postgresHost := os.Getenv("POSTGRES_HOST")
    postgresPort := os.Getenv("POSTGRES_PORT")
    postgresUser := os.Getenv("POSTGRES_USER")
    postgresPassword := os.Getenv("POSTGRES_PASSWORD")
    postgresDB := os.Getenv("POSTGRES_DB")

    // Format the connection string (DSN) for PostgreSQL
    postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        postgresHost, postgresPort, postgresUser, postgresPassword, postgresDB)

    // Open the database connection
    PGDB, err := sql.Open("postgres", postgresDSN)
    if err != nil {
        log.Fatal("Error opening connection to PostgreSQL:", err)
    }

    // Ping the database to ensure the connection is working
    err = PGDB.Ping()
    if err != nil {
        log.Fatal("Error pinging PostgreSQL:", err)
    }

    // Assign the open DB connection to the global PG variable
    PG = PGDB
    log.Println("Connected to PostgreSQL successfully!")
}

 
