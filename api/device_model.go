package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/driver"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

// deviceModelAPIPath is the path for the device model API.
const deviceModelAPIPath = "/device_models"

// initDeviceModelAPI initializes the device model API.
func (api *API) initDeviceModelAPI() {
	api.DeviceModels = api.APIRoot.Group(deviceModelAPIPath)

	api.DeviceModels.Use(isAuthorized)

	api.DeviceModels.Get("/", getDeviceModels)
	api.DeviceModels.Get("/drivers", getDrivers)
	api.DeviceModels.Get("/:id", getDeviceModel)
	api.DeviceModels.Post("/", createDeviceModel)
	api.DeviceModels.Put("/:id", updateDeviceModel)
	api.DeviceModels.Delete("/:id", deleteDeviceModel)
}

// getDeviceModels gets all device models.
func getDeviceModels(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	deviceModels, err := api.Store.DeviceModelStore.GetAll()
	if err != nil {
		api.Logger.Err(err, "failed to get device models")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get device models")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("device_models", deviceModels))
}

// getDeviceModel gets a device model by ID.
func getDeviceModel(c *fiber.Ctx) error {
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

	deviceModel, err := api.Store.DeviceModelStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device model")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("device_model", deviceModel))
}

// createDeviceModel creates a new device model.
func createDeviceModel(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	deviceModel := &model.DeviceModel{}
	if err := c.BodyParser(deviceModel); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := api.Store.DeviceModelStore.Create(deviceModel); err != nil {
		api.Logger.Err(err, "failed to create the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create the device model")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("id", deviceModel.ID))
}

// updateDeviceModel updates a device model.
func updateDeviceModel(c *fiber.Ctx) error {
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

	deviceModel, err := api.Store.DeviceModelStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device model")
	}

	if err := c.BodyParser(deviceModel); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := api.Store.DeviceModelStore.Update(deviceModel); err != nil {
		api.Logger.Err(err, "failed to update the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to update the device model")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("id", deviceModel.ID))
}

// deleteDeviceModel deletes a device model.
func deleteDeviceModel(c *fiber.Ctx) error {
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

	deviceModel, err := api.Store.DeviceModelStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the device model")
	}

	if err := api.Store.DeviceModelStore.Delete(deviceModel); err != nil {
		api.Logger.Err(err, "failed to delete the device model")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to delete the device model")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("id", deviceModel.ID))
}

// getDrivers gets all drivers.
func getDrivers(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	drivers := driver.Drivers()

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("drivers", drivers))
}
