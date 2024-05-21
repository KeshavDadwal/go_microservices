package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orders/microservices/controllers"
)

func SetupUserRoutes(app *fiber.App) {
	app.Post("/api/v1/order", controllers.HandlerCreateOrder)
}
