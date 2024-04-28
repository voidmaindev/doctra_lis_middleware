package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

// userAPIPath is the path for the user API.
const userAPIPath = "/users"

// initUserAPI initializes the user API.
func (api *API) initUserAPI() {
	api.Users = api.APIRoot.Group(userAPIPath)

	api.Users.Post("/register", registerUser)
	api.Users.Post("/token", token)

	api.Users.Use(isAuthorized)

	api.Users.Get("/", getUsers)
}

// registerUser registers a user.
func registerUser(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	user := &model.User{}
	if err := c.BodyParser(user); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	user.SetDefaultRole()

	err = user.HashPassword()
	if err != nil {
		api.Logger.Err(err, "failed to hash the password")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to hash the password")
	}

	if err := api.Store.UserStore.Create(user); err != nil {
		api.Logger.Err(err, "failed to create the user")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create the user")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("user", user.Username))
}

// token generates a JWT token.
func token(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	credentials := &authUser{}
	if err := c.BodyParser(credentials); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	user, err := api.Store.UserStore.GetByUsername(credentials.Username)
	if err != nil {
		api.Logger.Err(err, "failed to get the user by username")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the user by username")
	}

	if user == nil || !user.CheckPassword(credentials.Password) {
		return apiResponseError(c, fiber.StatusUnauthorized, "invalid username or password")
	}

	jwtToken, err := GenerateJWTToken(user)
	if err != nil {
		api.Logger.Err(err, "failed to generate JWT token")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to generate JWT token")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("token", jwtToken))
}

// getUsers gets all users.
func getUsers(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the api from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the api from context")
	}

	users, err := api.Store.UserStore.GetAll()
	if err != nil {
		api.Logger.Err(err, "failed to get users")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get users")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("users", users))
}

