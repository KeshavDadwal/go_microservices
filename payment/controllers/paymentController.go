package controllers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/payments/microservices/models"
)

var ctx = context.Background()

type Payment struct {
	PaymentID int    `json:"payment_id"`
	Status    string `json:"status"`
}

func HandlerCreatePayment(c *fiber.Ctx) error {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	order, err := getOrderFromRedis(rdb)
	if err != nil {
		log.Printf("Failed there is no order will be there: %v", err)
		return err
	}

	var req_payment Payment
	if err := c.BodyParser(&req_payment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Invalid request body",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	payment := models.Payment{
		OrderID: order.OrderID,
		Amount:  order.TotalPrice,
		Status:  req_payment.Status,
	}

	query := `INSERT INTO payments (order_id, amount, status) VALUES ($1, $2, $3) RETURNING payment_id`
	err = DB.QueryRow(ctx, query, payment.OrderID, payment.Amount, payment.Status).Scan(&payment.PaymentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error creating payment",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	err = rdb.Del(ctx, "Order_created").Err()
	if err != nil {
		log.Printf("Failed to delete the Order_created: %v", err)
		return err
	}

	err = publishPaymentToRedis(rdb, payment)
	if err != nil {
		log.Printf("Failed to publish payment to Redis channel: %v", err)
	}

	log.Println("Payment created successfully")
	return c.Status(fiber.StatusCreated).JSON(models.JsonResponse{
		Status:     true,
		Message:    "Payment created successfully",
		Data:       payment,
		StatusCode: fiber.StatusCreated,
	})
}

func getOrderFromRedis(rdb *redis.Client) (models.Order, error) {

	orderJSON, err := rdb.Get(ctx, "Order_created").Result()
	if err != nil {
		return models.Order{}, err
	}

	var order models.Order
	err = json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func publishPaymentToRedis(rdb *redis.Client, payment models.Payment) error {
	// Convert order struct to JSON
	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, "payment_created", paymentJSON, 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	err = rdb.Publish(ctx, "payment_created_event", paymentJSON).Err()
	if err != nil {
		log.Printf("Failed to publish order creation event: %v", err)
		return err
	}

	return nil
}
