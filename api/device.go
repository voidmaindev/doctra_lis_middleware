package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

// deviceAPIPath is the path for the device API.
const deviceAPIPath = "/devices"

// initDeviceAPI initializes the device API.
func (api *API) initDeviceAPI() {
	api.Devices = api.APIRoot.Group(deviceAPIPath)

	api.Devices.Use(isAuthorized)

	api.Devices.Get("/", getDevices)
	api.Devices.Get("/:id", getDevice)
	api.Devices.Post("/", createDevice)
	api.Devices.Put("/:id", updateDevice)
	api.Devices.Delete("/:id", deleteDevice)
}

// getDevices gets all devices.
func getDevices(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	devices, err := api.Store.DeviceStore.GetAll()
	if err != nil {
		api.Logger.Err(err, "failed to get devices")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get devices")
	}

	return apiResponseData(c, fiber.StatusOK, devices)
}

// getDevice gets a device by ID.
func getDevice(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		api.Logger.Err(err, "failed to parse the ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the ID")
	}

	device, err := api.Store.DeviceStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device")
	}

	return apiResponseData(c, fiber.StatusOK, device)
}

// createDevice creates a new device.
func createDevice(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	device := &model.Device{}
	if err := c.BodyParser(device); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := api.Store.DeviceStore.Create(device); err != nil {
		api.Logger.Err(err, "failed to create the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create the device")
	}

	return apiResponseData(c, fiber.StatusCreated, device.ID)
}

// updateDevice updates a device.
func updateDevice(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		api.Logger.Err(err, "failed to parse the ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the ID")
	}

	device, err := api.Store.DeviceStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device")
	}

	if err := c.BodyParser(device); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := api.Store.DeviceStore.Update(device); err != nil {
		api.Logger.Err(err, "failed to update the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to update the device")
	}

	return apiResponseData(c, fiber.StatusOK, device.ID)
}

// deleteDevice deletes a device.
func deleteDevice(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		api.Logger.Err(err, "failed to parse the ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the ID")
	}

	device, err := api.Store.DeviceStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device")
	}

	if err := api.Store.DeviceStore.Delete(device); err != nil {
		api.Logger.Err(err, "failed to delete the device")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to delete the device")
	}

	return apiResponseData(c, fiber.StatusOK, device.ID)
}
