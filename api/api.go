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

	return api, nil
}

func getAppFromContext(c *fiber.Ctx) (*API, error) {
	app, ok := c.Locals("app").(*API)
	if !ok {
		return nil, errors.New("failed to get the app from context")
	}

	return app, nil
}

func apiResponce(c *fiber.Ctx, status int, message string, data interface{}) error {
	success, msg := true, ""
	if status != fiber.StatusOK {
		success = false
		msg = message
	}

	fm := fiber.Map{}
	if data == nil {
		fm = fiber.Map{"success": success, "message": msg}
	} else {
		fm = fiber.Map{"success": success, "message": msg, "data": data}
	}

	return c.Status(status).JSON(fm)
}

func apiResponceError(c *fiber.Ctx, status int, message string) error {
	return apiResponce(c, status, message, nil)
}

func apiResponceData(c *fiber.Ctx, status int, data interface{}) error {
	return apiResponce(c, status, "", data)
}
