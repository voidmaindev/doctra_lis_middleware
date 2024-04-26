package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

const userAPIPath = "/users"

func (api *API) initUserAPI() {
	api.Users = api.APIRoot.Group(userAPIPath)

	api.Users.Post("/register", register)
	api.Users.Post("/token", token)
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

func register(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	user := &model.User{}
	if err := c.BodyParser(user); err != nil {
		app.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	if err := app.Store.UserStore.Create(user); err != nil {
		app.Logger.Err(err, "failed to create the user")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create the user")
	}

	return apiResponseData(c, fiber.StatusOK, user)
}

func token(c *fiber.Ctx) error {
	app, err := getAppFromContext(c)
	if err != nil {
		app.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	credentials := &authUser{}
	if err := c.BodyParser(credentials); err != nil {
		app.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	user, err := app.Store.UserStore.GetByUsername(credentials.Username)
	if err != nil {
		app.Logger.Err(err, "failed to get the user by email")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the user by email")
	}

	// if user == nil || !user.ComparePassword(credentials.Password) {
	// 	return apiResponseError(c, fiber.StatusUnauthorized, "invalid email or password")
	// }

	// token, err := user.GenerateToken()
	// if err != nil {
	// 	app.Logger.Err(err, "failed to generate the token")
	// 	return apiResponseError(c, fiber.StatusInternalServerError, "failed to generate the token")
	// }
	token := user.ID

	return apiResponseData(c, fiber.StatusOK, token)
}