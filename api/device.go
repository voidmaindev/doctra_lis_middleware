package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

const deviceAPIPath = "/devices"

func (api *API) initDeviceAPI() {
	api.Devices = api.APIRoot.Group(deviceAPIPath)

	api.Devices.Get("/", getDevices)
	api.Devices.Get("/:id", getDevice)
	api.Devices.Post("/", createDevice)
	api.Devices.Put("/:id", updateDevice)
	api.Devices.Delete("/:id", deleteDevice)
}

func getDevices(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	devices, err := app.Store.DeviceStore.GetAll()
	if err != nil {
		app.Logger.Err(err, "failed to get devices")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get devices")
	}

	return apiResponseData(c, fiber.StatusOK, devices)
}

func getDevice(c *fiber.Ctx) error {
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

	device, err := app.Store.DeviceStore.GetByID(uint(id))
	if err != nil {
		app.Logger.Err(err, "failed to get the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device")
	}

	return apiResponseData(c, fiber.StatusOK, device)
}

func createDevice(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	device := &model.Device{}
	if err := c.BodyParser(device); err != nil {
		app.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := app.Store.DeviceStore.Create(device); err != nil {
		app.Logger.Err(err, "failed to create the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create the device")
	}

	return apiResponseData(c, fiber.StatusCreated, device)
}

func updateDevice(c *fiber.Ctx) error {
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

	device, err := app.Store.DeviceStore.GetByID(uint(id))
	if err != nil {
		app.Logger.Err(err, "failed to get the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device")
	}

	if err := c.BodyParser(device); err != nil {
		app.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := app.Store.DeviceStore.Update(device); err != nil {
		app.Logger.Err(err, "failed to update the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to update the device")
	}

	return apiResponseData(c, fiber.StatusOK, device)
}

func deleteDevice(c *fiber.Ctx) error {
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

	device, err := app.Store.DeviceStore.GetByID(uint(id))
	if err != nil {
		app.Logger.Err(err, "failed to get the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device")
	}

	if err := app.Store.DeviceStore.Delete(device); err != nil {
		app.Logger.Err(err, "failed to delete the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to delete the device")
	}

	return apiResponseData(c, fiber.StatusOK, device)
}