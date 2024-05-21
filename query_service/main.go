package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/query_service/microservices/db"
	"github.com/query_service/microservices/models"
	"github.com/query_service/microservices/routes"
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var wg sync.WaitGroup

	// Set up order event listener
	wg.Add(1)
	go handleOrderCreatedEvents(rdb, &wg)

	// Set up payment event listener
	wg.Add(1)
	go handlePaymentCreatedEvents(rdb, &wg)

	setupServer(app)

	wg.Wait()
}

func handleOrderCreatedEvents(rdb *redis.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	pubsub := rdb.Subscribe(context.Background(), "order_created_event")
	defer pubsub.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-stop:
			return
		case msg := <-pubsub.Channel():
			fmt.Printf("Received order_created_event message: %s\n", msg.Payload)

			var newOrder models.Order
			if err := json.Unmarshal([]byte(msg.Payload), &newOrder); err != nil {
				fmt.Printf("Error unmarshalling JSON: %v\n", err)
				continue
			}

			query := `INSERT INTO orders (order_id, user_id, product_id, quantity, total_price, status) VALUES ($1, $2, $3, $4, $5, $6)`
			if _, err := DB.Exec(context.Background(), query, newOrder.OrderID, newOrder.UserID, newOrder.ProductID, newOrder.Quantity, newOrder.TotalPrice, newOrder.Status); err != nil {
				log.Printf("Error creating orders: %v\n", err)
				continue
			}

			log.Println("Order created successfully")
		}
	}
}

func handlePaymentCreatedEvents(rdb *redis.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	pubsub := rdb.Subscribe(context.Background(), "payment_created_event")
	defer pubsub.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-stop:
			return
		case msg := <-pubsub.Channel():
			fmt.Printf("Received payment_created_event message: %s\n", msg.Payload)

			var payment models.Payment
			if err := json.Unmarshal([]byte(msg.Payload), &payment); err != nil {
				fmt.Printf("Error unmarshalling JSON: %v\n", err)
				continue
			}

			query := `INSERT INTO payments (payment_id, order_id, amount, status) VALUES ($1, $2, $3, $4) RETURNING payment_id`
			if err := DB.QueryRow(context.Background(), query, payment.PaymentID, payment.OrderID, payment.Amount, payment.Status).Scan(&payment.PaymentID); err != nil {
				log.Printf("Error processing payment: %v\n", err)
				continue
			}

			query1 := `UPDATE orders SET status = $1 WHERE order_id = $2`
			if _, err := DB.Exec(context.Background(), query1, "approved", payment.OrderID); err != nil {
				log.Printf("Error updating order status: %v\n", err)
				continue
			}

			log.Println("Payment processed Successfully")
		}
	}
}

func setupServer(app *fiber.App) {
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
