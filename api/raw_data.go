package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

// rawDataAPIPath is the path for the raw data API.
const rawDataAPIPath = "/raw_data"

// initRawDataAPI initializes the raw data API.
func (api *API) initRawDataAPI() {
	api.RawData = api.APIRoot.Group(rawDataAPIPath)

	api.RawData.Use(isAuthorized)

	api.RawData.Get("/", getRawDatas)
	api.RawData.Get("/:id", getRawData)
	api.RawData.Get("/device/:device_id", getRawDataByDeviceID)
	api.RawData.Post("/", createRawData)
	api.RawData.Put("/:id", updateRawData)
	api.RawData.Delete("/:id", deleteRawData)
}

// getRawDatas gets all raw datas.
func getRawDatas(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	rawData, err := api.Store.RawDataStore.GetAll()
	if err != nil {
		api.Logger.Err(err, "failed to get raw datas")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get raw datas")
	}

	rawDataAPIs := make([]*model.RawDataApi, 0, len(rawData))
	for _, r := range rawData {
		rawDataAPIs = append(rawDataAPIs, model.NewRawDataApi(r))
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("raw_data", rawDataAPIs))
}

// getRawData gets a raw data by ID.
func getRawData(c *fiber.Ctx) error {
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

	rawData, err := api.Store.RawDataStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the raw data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the raw data")
	}

	rawDataAPI := model.NewRawDataApi(rawData)

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("raw_data", rawDataAPI))
}

// getRawDataByDeviceID gets a raw data by device ID.
func getRawDataByDeviceID(c *fiber.Ctx) error {
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

	rawData, err := api.Store.RawDataStore.GetByDeviceID(uint(deviceID))
	if err != nil {
		api.Logger.Err(err, "failed to get the raw data by device ID")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the raw data by device ID")
	}

	rawDataAPIs := make([]*model.RawDataApi, 0, len(rawData))
	for _, r := range rawData {
		rawDataAPIs = append(rawDataAPIs, model.NewRawDataApi(r))
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("raw_data", rawDataAPIs))
}

// createRawData creates a raw data.
func createRawData(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	rawData := &model.RawData{}
	if err := c.BodyParser(rawData); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := api.Store.RawDataStore.Create(rawData); err != nil {
		api.Logger.Err(err, "failed to create the raw data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create the raw data")
	}

	return apiResponseData(c, fiber.StatusCreated, NewAPIRV("id", rawData.ID))
}

// updateRawData updates a raw data.
func updateRawData(c *fiber.Ctx) error {
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

	rawData, err := api.Store.RawDataStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the raw data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the raw data")
	}

	if err := c.BodyParser(rawData); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := api.Store.RawDataStore.Update(rawData); err != nil {
		api.Logger.Err(err, "failed to update the raw data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to update the raw data")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("id", rawData.ID))
}

// deleteRawData deletes a raw data.
func deleteRawData(c *fiber.Ctx) error {
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

	rawData, err := api.Store.RawDataStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the raw data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the raw data")
	}

	if err := api.Store.RawDataStore.Delete(uint(id)); err != nil {
		api.Logger.Err(err, "failed to delete the raw data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to delete the raw data")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("id", rawData.ID))
}
