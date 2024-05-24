package api

import (
	"strings"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
)

// NewRouter creates a new router.
func NewRouter(logger *zerolog.Logger) *fiber.App {
	r := fiber.New()

	r.Use(fiberzerolog.New(fiberzerolog.Config{Logger: logger}))
	r.Use(recover.New())
	r.Use(newCors())

	return r
}

// newCors creates a new CORS middleware.
func newCors() fiber.Handler {
	corsConfig := cors.Config{
		AllowHeaders: "Origin,Content-Type,Cache-Control,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins: "*",
		// AllowCredentials: true,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		}, ","),
	}

	return cors.New(corsConfig)
}
