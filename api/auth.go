package api

type authUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateToken(username string, password string) (string, error) {
	return "", nil
}