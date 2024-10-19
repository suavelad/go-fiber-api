package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/suavelad/go-fibre-api/database"
	"github.com/suavelad/go-fibre-api/routes"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to my app")
}

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/hello", welcome)

	//Users
	app.Post("api/v1/users", routes.CreateUser)
	app.Get("api/v1/users", routes.GetUsers)
	app.Get("api/v1/users/:id", routes.GetUser)
	app.Put("api/v1/users/:id", routes.UpdateUser)
	app.Delete("api/v1/users/:id", routes.DeleteUser)

	//Products
	app.Post("api/v1/products", routes.CreateProduct)
	app.Get("api/v1/products", routes.GetProducts)
	app.Get("api/v1/products/:id", routes.GetProduct)
	app.Put("api/v1/products/:id", routes.UpdateProduct)
	app.Delete("api/v1/products/:id", routes.DeleteProduct)

	//Orders
	app.Post("api/v1/orders", routes.CreateOrder)
	app.Get("api/v1/orders", routes.GetOrders)
	app.Get("api/v1/orders/:id", routes.GetOrder)
	app.Put("api/v1/orders/:id", routes.UpdateOrder)
	app.Delete("api/v1/orders/:id", routes.DeleteOrder)
}

func main() {
	database.ConnectDb()

	app := fiber.New()

	setupRoutes(app)

	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://gofiber.io, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/ready",
	}))

	app.Listen(":3000")
}
