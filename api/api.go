// Package api provides the API for the application.
package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

// API is the API for the application.
const apiRootPath = "/api/v1"

// API is the API struct for the application.
type API struct {
	Logger *log.Logger
	Root   *fiber.App
	Store  *store.Store

	APIRoot      fiber.Router
	Users        fiber.Router
	DeviceModels fiber.Router
	Devices      fiber.Router
	LabDatas     fiber.Router
}

// NewAPI creates a new API.
func NewAPI(logger *log.Logger, router *fiber.App, store *store.Store) (*API, error) {
	api := &API{
		Logger: logger,
		Root:   router,
		Store:  store,
	}

	api.Root.Use(func(c *fiber.Ctx) error {
		c.Locals("api", api)
		return c.Next()
	})

	api.APIRoot = api.Root.Group(apiRootPath)

	api.initUserAPI()
	api.initDeviceModelAPI()
	api.initDeviceAPI()
	api.initLabDataAPI()

	api.addNoRoute()

	return api, nil
}

// getApiFromContext gets the API from the context.
func getApiFromContext(c *fiber.Ctx) (*API, error) {
	api, ok := c.Locals("api").(*API)
	if !ok {
		return nil, errors.New("failed to get the app from context")
	}

	return api, nil
}

// apiResponse sends a response.
func apiResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	success, msg := true, "success"
	if status != fiber.StatusOK {
		success = false
		msg = message
	}

	fm := fiber.Map{"success": success, "message": msg}
	if data != nil {
		fm = fiber.Map{"success": success, "message": msg, "data": data}
	}

	return c.Status(status).JSON(fm)
}

// apiResponseError sends an error response.
func apiResponseError(c *fiber.Ctx, status int, message string) error {
	return apiResponse(c, status, message, nil)
}

// apiResponseData sends a data response.
func apiResponseData(c *fiber.Ctx, status int, data interface{}) error {
	return apiResponse(c, status, "", data)
}

// addNoRoute adds a no route handler.
func (api *API) addNoRoute() {
	api.Root.Use(func(c *fiber.Ctx) error {
		return apiResponseError(c, fiber.StatusNotFound, "not found")
	})
}
