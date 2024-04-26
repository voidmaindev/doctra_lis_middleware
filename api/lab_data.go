package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

// labDataAPIPath is the path for the lab data API.
const labDataAPIPath = "/lab_data"

// initLabDataAPI initializes the lab data API.
func (api *API) initLabDataAPI() {
	api.LabDatas = api.APIRoot.Group(labDataAPIPath)

	api.LabDatas.Use(isAuthorized)

	api.LabDatas.Get("/", getLabDatas)
	api.LabDatas.Get("/:id", getLabData)
	api.LabDatas.Get("/barcode/:barcode", getLabDataByBarcode)
	api.LabDatas.Get("/device/:device_id", getLabDataByDeviceID)
	api.LabDatas.Get("/device/:device_id/barcode/:barcode", getLabDataByDeviceIDAndBarcode)
	api.LabDatas.Post("/", createLabData)
	api.LabDatas.Put("/:id", updateLabData)
	api.LabDatas.Delete("/:id", deleteLabData)
}

// getLabDatas gets all lab datas.
func getLabDatas(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	labDatas, err := api.Store.LabDataStore.GetAll()
	if err != nil {
		api.Logger.Err(err, "failed to get lab datas")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get lab datas")
	}

	return apiResponseData(c, fiber.StatusOK, labDatas)
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

	return apiResponseData(c, fiber.StatusOK, labData)
}

// getLabDataByBarcode gets a lab data by barcode.
func getLabDataByBarcode(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	barcode := c.Params("barcode")
	labDatas, err := api.Store.LabDataStore.GetByBarcode(barcode)
	if err != nil {
		api.Logger.Err(err, "failed to get the lab data by barcode")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data by barcode")
	}

	return apiResponseData(c, fiber.StatusOK, labDatas)
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

	labDatas, err := api.Store.LabDataStore.GetByDeviceID(uint(deviceID))
	if err != nil {
		api.Logger.Err(err, "failed to get the lab data by device ID")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data by device ID")
	}

	return apiResponseData(c, fiber.StatusOK, labDatas)
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

	return apiResponseData(c, fiber.StatusOK, labData)
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

	return apiResponseData(c, fiber.StatusOK, labData)
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

	return apiResponseData(c, fiber.StatusOK, labData)
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

	return apiResponseData(c, fiber.StatusOK, labData)
}
