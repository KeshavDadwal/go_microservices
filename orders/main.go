package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/orders/microservices/db"
	"github.com/orders/microservices/routes"
)

func serveStatic(app *fiber.App) {
	app.Static("/", "./views")
}

func main() {
	dbPool, _ := db.ConnectToDB()
	defer dbPool.Close()

	app := fiber.New()
	serveStatic(app)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:4200",
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
