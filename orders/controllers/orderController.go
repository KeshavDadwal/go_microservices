package controllers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/orders/microservices/models"
)

var ctx = context.Background()

func HandlerCreateOrder(c *fiber.Ctx) error {
	var newOrder models.Order
	if err := c.BodyParser(&newOrder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	query := `INSERT INTO orders (user_id, product_id, quantity, total_price, status) VALUES ($1, $2, $3, $4, $5)`
	_, err := DB.Exec(ctx, query, newOrder.UserID, newOrder.ProductID, newOrder.Quantity, newOrder.TotalPrice, newOrder.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error creating orders",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	err = DB.QueryRow(ctx, "SELECT lastval()").Scan(&newOrder.OrderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error retrieving orders",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	// Publish order information to Redis channel
	err = publishOrderToRedis(newOrder)
	if err != nil {
		log.Printf("Failed to publish order to Redis channel: %v", err)
	}

	log.Println( "Order created successfully")
	return c.Status(fiber.StatusCreated).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Order created successfully",
		Data:       newOrder,
		StatusCode: fiber.StatusCreated,
	})
}

func publishOrderToRedis(order models.Order) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, "Order_created", orderJSON, 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	err = rdb.Publish(ctx, "order_created_event", orderJSON).Err()
	if err != nil {
		log.Printf("Failed to publish order creation event: %v", err)
		return err
	}

	return nil
}
