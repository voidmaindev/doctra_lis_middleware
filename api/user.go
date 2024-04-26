package api

import (
	"github.com/gofiber/fiber/v2"
)

const userAPIPath = "/users"

func (api *API) initUserAPI() {
	api.Users = api.APIRoot.Group(userAPIPath)

	api.Users.Get("/", getUsers)
}

func getUsers(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	users, err := app.Store.UserStore.GetAll()
	if err != nil {
		app.Logger.Err(err, "failed to get users")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get users")
	}

	return apiResponseData(c, fiber.StatusOK, users)
}
