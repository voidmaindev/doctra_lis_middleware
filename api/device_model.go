package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

const deviceModelAPIPath = "/device_models"

func (api *API) initDeviceModelAPI() {
	api.DeviceModels = api.APIRoot.Group(deviceModelAPIPath)

	api.DeviceModels.Get("/", getDeviceModels)
	api.DeviceModels.Get("/:id", getDeviceModel)
	api.DeviceModels.Post("/", createDeviceModel)
	api.DeviceModels.Put("/:id", updateDeviceModel)
	api.DeviceModels.Delete("/:id", deleteDeviceModel)
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

func getDeviceModel(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		app.Logger.Err(err, "failed to parse the ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the ID")
	}

	deviceModel, err := app.Store.DeviceModelStore.GetByID(uint(id))
	if err != nil {
		app.Logger.Err(err, "failed to get the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device model")
	}

	return apiResponseData(c, fiber.StatusOK, deviceModel)
}

func createDeviceModel(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	deviceModel := &model.DeviceModel{}
	if err := c.BodyParser(deviceModel); err != nil {
		app.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := app.Store.DeviceModelStore.Create(deviceModel); err != nil {
		app.Logger.Err(err, "failed to create the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create the device model")
	}

	return apiResponseData(c, fiber.StatusOK, deviceModel)
}

func updateDeviceModel(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		app.Logger.Err(err, "failed to parse the ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the ID")
	}

	deviceModel, err := app.Store.DeviceModelStore.GetByID(uint(id))
	if err != nil {
		app.Logger.Err(err, "failed to get the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device model")
	}

	if err := c.BodyParser(deviceModel); err != nil {
		app.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := app.Store.DeviceModelStore.Update(deviceModel); err != nil {
		app.Logger.Err(err, "failed to update the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to update the device model")
	}

	return apiResponseData(c, fiber.StatusOK, deviceModel)
}

func deleteDeviceModel(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		app.Logger.Err(err, "failed to parse the ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the ID")
	}

	deviceModel, err := app.Store.DeviceModelStore.GetByID(uint(id))
	if err != nil {
		app.Logger.Err(err, "failed to get the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device model")
	}

	if err := app.Store.DeviceModelStore.Delete(deviceModel); err != nil {
		app.Logger.Err(err, "failed to delete the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to delete the device model")
	}

	return apiResponseData(c, fiber.StatusOK, deviceModel)
}
