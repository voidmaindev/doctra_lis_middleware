package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

// labDataAPIPath is the path for the lab data API.
const labDataAPIPath = "/lab_data"

// initLabDataAPI initializes the lab data API.
func (api *API) initLabDataAPI() {
	api.LabData = api.APIRoot.Group(labDataAPIPath)

	api.LabData.Use(isAuthorized)

	api.LabData.Get("/", getLabDatas)
	api.LabData.Get("/:id", getLabData)
	api.LabData.Get("/barcode/:barcode", getLabDataByBarcode)
	api.LabData.Get("/device/:device_id", getLabDataByDeviceID)
	api.LabData.Get("/device/:device_id/barcode/:barcode", getLabDataByDeviceIDAndBarcode)
	api.LabData.Get("/serial/:serial", getLabDataBySerial)
	api.LabData.Get("/serial/:serial/barcode/:barcode", getLabDataBySerialAndBarcode)
	api.LabData.Post("/", createLabData)
	api.LabData.Put("/:id", updateLabData)
	api.LabData.Delete("/:id", deleteLabData)
}

// getLabDatas gets all lab datas.
func getLabDatas(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	labData, err := api.Store.LabDataStore.GetAll()
	if err != nil {
		api.Logger.Err(err, "failed to get lab datas")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get lab datas")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("lab_data", labData))
}

// getLabData gets a lab data by ID.
func getLabData(c *fiber.Ctx) error {
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

	labData, err := api.Store.LabDataStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("lab_data", labData))
}

// getLabDataByBarcode gets a lab data by barcode.
func getLabDataByBarcode(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	barcode := c.Params("barcode")
	labData, err := api.Store.LabDataStore.GetByBarcode(barcode)
	if err != nil {
		api.Logger.Err(err, "failed to get the lab data by barcode")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data by barcode")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("lab_data", labData))
}

// getLabDataByDeviceID gets lab data by device ID.
func getLabDataByDeviceID(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	deviceID, err := c.ParamsInt("device_id")
	if err != nil {
		api.Logger.Err(err, "failed to parse the device ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the device ID")
	}

	labData, err := api.Store.LabDataStore.GetByDeviceID(uint(deviceID))
	if err != nil {
		api.Logger.Err(err, "failed to get the lab data by device ID")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data by device ID")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("lab_data", labData))
}

// getLabDataByDeviceIDAndBarcode gets a lab data by device ID and barcode.
func getLabDataByDeviceIDAndBarcode(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	deviceID, err := c.ParamsInt("device_id")
	if err != nil {
		api.Logger.Err(err, "failed to parse the device ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the device ID")
	}

	barcode := c.Params("barcode")
	labData, err := api.Store.LabDataStore.GetByDeviceIDAndBarcode(uint(deviceID), barcode)
	if err != nil {
		api.Logger.Err(err, "failed to get the lab data by device ID and barcode")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data by device ID and barcode")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("lab_data", labData))
}

// getLabDataBySerial gets a lab data by serial.
func getLabDataBySerial(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	serial := c.Params("serial")
	labData, err := api.Store.LabDataStore.GetBySerial(serial)
	if err != nil {
		api.Logger.Err(err, "failed to get the lab data by serial")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data by serial")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("lab_data", labData))
}

// getLabDataBySerialAndBarcode gets a lab data by serial and barcode.
func getLabDataBySerialAndBarcode(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	serial := c.Params("serial")
	barcode := c.Params("barcode")
	labData, err := api.Store.LabDataStore.GetBySerialAndBarcode(serial, barcode)
	if err != nil {
		api.Logger.Err(err, "failed to get the lab data by serial and barcode")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data by serial and barcode")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("lab_data", labData))
}

// createLabData creates a new lab data.
func createLabData(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	labData := &model.LabData{}
	if err := c.BodyParser(labData); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := api.Store.LabDataStore.Create(labData); err != nil {
		api.Logger.Err(err, "failed to create the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create the lab data")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("id", labData.ID))
}

// updateLabData updates a lab data.
func updateLabData(c *fiber.Ctx) error {
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

	labData, err := api.Store.LabDataStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data")
	}

	if err := c.BodyParser(labData); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := api.Store.LabDataStore.Update(labData); err != nil {
		api.Logger.Err(err, "failed to update the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to update the lab data")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("id", labData.ID))
}

// deleteLabData deletes a lab data.
func deleteLabData(c *fiber.Ctx) error {
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

	labData, err := api.Store.LabDataStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data")
	}

	if err := api.Store.LabDataStore.Delete(labData); err != nil {
		api.Logger.Err(err, "failed to delete the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to delete the lab data")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("id", labData.ID))
}
