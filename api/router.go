package api

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
)

func NewRouter(logger *zerolog.Logger) *fiber.App {
	r := fiber.New()

	r.Use(fiberzerolog.New(fiberzerolog.Config{Logger: logger}))
	r.Use(recover.New())
	r.Use(cors.New())

	return r
}