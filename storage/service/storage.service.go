package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"storage/model"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func loadEnvVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error while loading .env file")
	}
}

func InitializeStore() {
	loadEnvVariables()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable.")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+uri))
	if err != nil {
		log.Fatalf("Failed to connect to the MongoDB database: %v", err)
	}

	collection = client.Database("nvd_scrapper").Collection("vulnerabilities")
}

func SaveVulnerability(vulnerability *model.Vulnerability) bool {
	result, err := collection.InsertOne(context.TODO(), vulnerability)
	if err != nil {
		fmt.Printf("Insertion failed due to some reasons: %v", err)
		return false
	}

	fmt.Printf("Inserted vulnerability with _id: %v\n", result.InsertedID)
	return true

}
