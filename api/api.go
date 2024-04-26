package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const apiRootPath = "/api/v1"

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

func NewAPI(logger *log.Logger, router *fiber.App, store *store.Store) (*API, error) {
	api := &API{
		Logger: logger,
		Root:   router,
		Store:  store,
	}

	api.APIRoot = api.Root.Group(apiRootPath)
	api.initUserAPI()
	api.initDeviceModelAPI()
	api.initDeviceAPI()
	api.initLabDataAPI()

	api.addNoRoute()

	return api, nil
}

func getAppFromContext(c *fiber.Ctx) (*API, error) {
	app, ok := c.Locals("app").(*API)
	if !ok {
		return nil, errors.New("failed to get the app from context")
	}

	return app, nil
}

func apiResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	success, msg := true, ""
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

func apiResponseError(c *fiber.Ctx, status int, message string) error {
	return apiResponse(c, status, message, nil)
}

func apiResponseData(c *fiber.Ctx, status int, data interface{}) error {
	return apiResponse(c, status, "", data)
}

func (api *API)addNoRoute() {
	api.Root.Use(func(c *fiber.Ctx) error {
		return apiResponseError(c, fiber.StatusNotFound, "not found")
	})
}

