package api

import "github.com/gofiber/fiber/v2"

const deviceModelAPIPath = "/device_models"

func (api *API) initDeviceModelAPI() {
	api.DeviceModels = api.APIRoot.Group(deviceModelAPIPath)

	api.DeviceModels.Get("/", getDeviceModels)
}

func getDeviceModels(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	deviceModels, err := app.Store.DeviceModelStore.GetAll()
	if err != nil {
		app.Logger.Err(err, "failed to get device models")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get device models")
	}

	return apiResponseData(c, fiber.StatusOK, deviceModels)
}
