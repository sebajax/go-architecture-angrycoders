package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/sebajax/go-vertical-slice-architecture/internal/user/handler"
	"github.com/sebajax/go-vertical-slice-architecture/pkg/injection"
	"github.com/sebajax/go-vertical-slice-architecture/pkg/middleware"
)

func main() {
	// prepare all components for dependency injection
	injection.ProvideComponents()

	// initiate service components with dependency injection
	if err := injection.InitComponents(); err != nil {
		panic(err)
	}

	// create fiber
	app := fiber.New()

	// add fiber middlewares
	app.Use(cors.New())
	app.Use(helmet.New())

	// custom middlewares
	app.Use(middleware.ErrorHandler)

	// create health end point
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Status ok - api running")
	})

	// add api group for users
	api := app.Group("/api")       // /api
	userApi := api.Group("/users") // /api/user
	handler.UserRouter(userApi, injection.UserServiceProvider)

	// listen in port 8080
	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("API_PORT"))))
}
