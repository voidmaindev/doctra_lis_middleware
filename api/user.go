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
	api.Users.Get("/:id", getUser)
	api.Users.Get("/username/:username", getUserByUsername)
	api.Users.Get("/me", getMe)
	api.Users.Put("/:id", updateUser)
	api.Users.Delete("/:id", deleteUser)
}

// registerUser registers a user.
func registerUser(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	authUser := &model.AuthUser{}
	if err := c.BodyParser(authUser); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	user, err := model.NewUserFromAuthUser(authUser)
	if err != nil {
		api.Logger.Err(err, "failed to create a new user from the auth user")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create a new user from the auth user")
	}

	user.SetDefaultRole()

	if err := api.Store.UserStore.Create(user); err != nil {
		api.Logger.Err(err, "failed to create the user")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to create the user")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("username", user.Username))
}

// token generates a JWT token.
func token(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	credentials := &model.AuthUser{}
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

	if !isAdmin(c) {
		return apiResponseError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	users, err := api.Store.UserStore.GetAll()
	if err != nil {
		api.Logger.Err(err, "failed to get users")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get users")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("users", users))
}

// getUser gets a user by ID.
func getUser(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the api from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the api from context")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		api.Logger.Err(err, "failed to parse the ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the ID")
	}

	if !isAdmin(c) {
		return apiResponseError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	user, err := api.Store.UserStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the user by ID")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the user by ID")
	}

	if user == nil {
		return apiResponseError(c, fiber.StatusNotFound, "user not found")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("user", user))
}

// getUserByUsername gets a user by username.
func getUserByUsername(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the api from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the api from context")
	}

	if !isAdmin(c) {
		return apiResponseError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	username := c.Params("username")
	user, err := api.Store.UserStore.GetByUsername(username)
	if err != nil {
		api.Logger.Err(err, "failed to get the user by username")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the user by username")
	}

	if user == nil {
		return apiResponseError(c, fiber.StatusNotFound, "user not found")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("user", user))
}

// getMe gets the current user.
func getMe(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the api from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the api from context")
	}

	username, ok := c.Locals("username").(string)
	if !ok {
		return apiResponseError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	user, err := api.Store.UserStore.GetByUsername(username)
	if err != nil {
		api.Logger.Err(err, "failed to get the user by username")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the user by username")
	}

	if user == nil {
		return apiResponseError(c, fiber.StatusNotFound, "user not found")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("user", user))
}

// updateUser updates a user.
func updateUser(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the api from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the api from context")
	}

	if !isAdmin(c) {
		return apiResponseError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		api.Logger.Err(err, "failed to parse the ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the ID")
	}

	user, err := api.Store.UserStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the user by ID")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the user by ID")
	}

	if user == nil {
		return apiResponseError(c, fiber.StatusNotFound, "user not found")
	}

	authUser := &model.AuthUser{}
	if err := c.BodyParser(authUser); err != nil {
		api.Logger.Err(err, "failed to parse the request body")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the request body")
	}

	err = user.UpdateFromAuthUser(authUser)
	if err != nil {
		api.Logger.Err(err, "failed to update the user from the auth user")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to update the user from the auth user")
	}

	if err := api.Store.UserStore.Update(user); err != nil {
		api.Logger.Err(err, "failed to update the user")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to update the user")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("username", user.Username))
}

// deleteUser deletes a user.
func deleteUser(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the api from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the api from context")
	}

	if !isAdmin(c) {
		return apiResponseError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		api.Logger.Err(err, "failed to parse the ID")
		return apiResponseError(c, fiber.StatusBadRequest, "failed to parse the ID")
	}

	user, err := api.Store.UserStore.GetByID(uint(id))
	if err != nil {
		api.Logger.Err(err, "failed to get the user by ID")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the user by ID")
	}

	if user == nil {
		return apiResponseError(c, fiber.StatusNotFound, "user not found")
	}

	if err := api.Store.UserStore.Delete(user); err != nil {
		api.Logger.Err(err, "failed to delete the user")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to delete the user")
	}

	return apiResponseData(c, fiber.StatusOK, NewAPIRV("username", user.Username))
}
