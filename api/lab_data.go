package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

const labDataAPIPath = "/lab_data"

func (api *API) initLabDataAPI() {
	api.LabDatas = api.APIRoot.Group(labDataAPIPath)

	api.LabDatas.Get("/", getLabDatas)
	api.LabDatas.Get("/:id", getLabData)
	api.LabDatas.Get("/barcode/:barcode", getLabDataByBarcode)
	api.LabDatas.Get("/device/:device_id", getLabDataByDeviceID)
	api.LabDatas.Get("/device/:device_id/barcode/:barcode", getLabDataByDeviceIDAndBarcode)
	api.LabDatas.Post("/", createLabData)
	api.LabDatas.Put("/:id", updateLabData)
	api.LabDatas.Delete("/:id", deleteLabData)
}

func getLabDatas(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	labDatas, err := app.Store.LabDataStore.GetAll()
	if err != nil {
		app.Logger.Err(err, "failed to get lab datas")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get lab datas")
	}

	return apiResponseData(c, fiber.StatusOK, labDatas)
}

func getLabData(c *fiber.Ctx) error {
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

	labData, err := app.Store.LabDataStore.GetByID(uint(id))
	if err != nil {
		app.Logger.Err(err, "failed to get the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data")
	}

	return apiResponseData(c, fiber.StatusOK, labData)
}

func getLabDataByBarcode(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	barcode := c.Params("barcode")
	labDatas, err := app.Store.LabDataStore.GetByBarcode(barcode)
	if err != nil {
		app.Logger.Err(err, "failed to get the lab data by barcode")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data by barcode")
	}

	return apiResponseData(c, fiber.StatusOK, labDatas)
}

func getLabDataByDeviceID(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	deviceID, err := c.ParamsInt("device_id")
	if err != nil {
		app.Logger.Err(err, "failed to parse the device ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the device ID")
	}

	labDatas, err := app.Store.LabDataStore.GetByDeviceID(uint(deviceID))
	if err != nil {
		app.Logger.Err(err, "failed to get the lab data by device ID")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data by device ID")
	}

	return apiResponseData(c, fiber.StatusOK, labDatas)
}

func getLabDataByDeviceIDAndBarcode(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	deviceID, err := c.ParamsInt("device_id")
	if err != nil {
		app.Logger.Err(err, "failed to parse the device ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the device ID")
	}

	barcode := c.Params("barcode")
	labData, err := app.Store.LabDataStore.GetByDeviceIDAndBarcode(uint(deviceID), barcode)
	if err != nil {
		app.Logger.Err(err, "failed to get the lab data by device ID and barcode")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data by device ID and barcode")
	}

	return apiResponseData(c, fiber.StatusOK, labData)
}

func createLabData(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	labData := &model.LabData{}
	if err := c.BodyParser(labData); err != nil {
		app.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := app.Store.LabDataStore.Create(labData); err != nil {
		app.Logger.Err(err, "failed to create the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create the lab data")
	}

	return apiResponseData(c, fiber.StatusOK, labData)
}

func updateLabData(c *fiber.Ctx) error {
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

	labData, err := app.Store.LabDataStore.GetByID(uint(id))
	if err != nil {
		app.Logger.Err(err, "failed to get the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data")
	}

	if err := c.BodyParser(labData); err != nil {
		app.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := app.Store.LabDataStore.Update(labData); err != nil {
		app.Logger.Err(err, "failed to update the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to update the lab data")
	}

	return apiResponseData(c, fiber.StatusOK, labData)
}

func deleteLabData(c *fiber.Ctx) error {
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

	labData, err := app.Store.LabDataStore.GetByID(uint(id))
	if err != nil {
		app.Logger.Err(err, "failed to get the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the lab data")
	}

	if err := app.Store.LabDataStore.Delete(labData); err != nil {
		app.Logger.Err(err, "failed to delete the lab data")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to delete the lab data")
	}

	return apiResponseData(c, fiber.StatusOK, labData)
}