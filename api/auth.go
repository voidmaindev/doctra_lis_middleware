package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/voidmaindev/doctra_lis_middleware/model"
)

const (
	jwtTTL    = 24 * time.Hour
	jwtSecret = "lksjfowejr!@#1ejk12ESLdKJHk12QW:Lsdfakl123"
	bearer    = "bearer "
)

// authUser is the structure for the authentication user.
type authUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// jwtCustomClaims is the custom claims for the JWT token.
type jwtCustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWTToken generates a JWT token.
func GenerateJWTToken(user *model.User) (string, error) {
	claims := &jwtCustomClaims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the token string
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// parseJWTToken parses a JWT token.
func parseJWTToken(token string) (*jwtCustomClaims, error) {
	claims := &jwtCustomClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !jwtToken.Valid {
		return nil, err
	}

	return claims, nil
}

// isAuthorized is a middleware to check if the user is authorized.
func isAuthorized(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	token := c.Get("Authorization")
	if token == "" || len(token) <= len(bearer) {
		api.Logger.Info("no token provided")
		return apiResponseError(c, fiber.StatusUnauthorized, "no token provided")
	}

	token = token[len(bearer):]

	claims, err := parseJWTToken(token)
	if err != nil {
		api.Logger.Err(err, "failed to parse the token")
		return apiResponseError(c, fiber.StatusUnauthorized, "incorrect token")
	}

	c.Locals("username", claims.Username)
	c.Locals("role", claims.Role)

	return c.Next()
}

// isAdmin is a middleware to check if the user is an admin.
func isAdmin(c *fiber.Ctx) bool {
	role, ok := c.Locals("role").(string)
	if !ok {
		return false
	}

	return role == model.RoleAdmin
}

// isUser is a middleware to check if the user is a user.
func isUser(c *fiber.Ctx) bool {
	role, ok := c.Locals("role").(string)
	if !ok {
		return false
	}

	return role == model.RoleUser
}
