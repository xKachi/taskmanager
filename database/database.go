package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	MongoDBURI        string
	MongoDBName       string
	MongoDBCollection string
}

var config *Config
var DB *mongo.Database

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	config = &Config{
		MongoDBURI:        os.Getenv("MONGODB_URI"),
		MongoDBName:       os.Getenv("MONGODB_DB_NAME"),
		MongoDBCollection: os.Getenv("MONGODB_COLLECTION"),
	}

}

func GetDBCollection() *mongo.Collection {
	return DB.Collection(config.MongoDBCollection)
}

func NewDBInstance() error {
	client, err := mongo.Connect(options.Client().ApplyURI(config.MongoDBURI))
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database(config.MongoDBName)
	return nil
}