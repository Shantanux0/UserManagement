package routes

import (
	"user-management/internal/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers the HTTP routes on the Fiber app.
func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	app.Post("/users", userHandler.CreateUser)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
	app.Get("/users", userHandler.ListUsers)
}
