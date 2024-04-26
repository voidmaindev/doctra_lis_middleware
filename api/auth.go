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
)

type authUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type jwtCustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

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

func isAuthorized(c *fiber.Ctx) error {
	api, err := getApiFromContext(c)
	if err != nil {
		api.Logger.Err(err, "failed to get the app from context")
		return apiResponseError(c, fiber.StatusInternalServerError, "failed to get the app from context")
	}

	token := c.Get("Authorization")
	if token == "" {
		api.Logger.Info("no token provided")
		return apiResponseError(c, fiber.StatusUnauthorized, "no token provided")
	}

	claims, err := parseJWTToken(token)
	if err != nil {
		api.Logger.Err(err, "failed to parse the token")
		return apiResponseError(c, fiber.StatusUnauthorized, "incorrect token")
	}

	c.Locals("username", claims.Username)
	c.Locals("role", claims.Role)

	return c.Next()
}
