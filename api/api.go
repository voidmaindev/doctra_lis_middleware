package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/log"
)

type API struct {
	Logger  *log.Logger
	Root    *fiber.App
	APIRoot *fiber.App
}

// func NewAPI(logger *log.Logger) *API {
// 	api := &API{

// 	return &API{
// 		Logger: logger,
// 		Root: fiber.New(),
// 		APIRoot: fiber.New(),
// 	}
// }
