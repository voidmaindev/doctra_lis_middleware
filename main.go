package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/voidmaindev/doctra_lis_middleware/inits"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

var addr string

func init() {
	inits.LoadEnvVars()
	addr = os.Getenv("ADDRESS") + ":" + os.Getenv("PORT")

	store.ConnectToDB()
}

func main() {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
			return c.SendString("Hello, World!")
	})

	app.Listen(addr)
}