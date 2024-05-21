package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/payments/microservices/controllers"
)

func SetupUserRoutes(app *fiber.App) {
	app.Post("/api/v1/payment", controllers.HandlerCreatePayment)
}
