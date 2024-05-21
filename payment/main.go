package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/payments/microservices/db"
	"github.com/payments/microservices/models"
	"github.com/payments/microservices/routes"
)


var ctx = context.Background()
var DB *pgxpool.Pool


func main() {
	db, err := db.ConnectToDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	DB = db

	app := fiber.New()

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Content-Type,X-Access-Token, Authorization",
	}))

	routes.SetupUserRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func getOrderFromRedis(orderJSON string) (models.Order, error) {
	var order models.Order
	err := json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		return models.Order{}, err
	}

	fmt.Println("Order data is there ",orderJSON)

	return order, nil
}




