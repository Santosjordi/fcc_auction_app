package main

import (
	"context"
	"log"

	mongo "github.com/Santosjordi/fcc_auction_app/configs/database/mongodb"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	dataBaseClient, err := mongo.NewMongoDBConection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
