package config

import (
	"errors"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type Token struct {
	Hash         string `json:"token"`
	Expire       int64  `mapstructure:"JWT_EXPIRE" json:"expires_in" yaml:"expires_in"`
	AppJwtSecret string `mapstructure:"APP_JWT_SECRET" yaml:"app_jwt_secret"`
	ApiJwtSecret string `mapstructure:"API_JWT_SECRET" yaml:"api_jwt_secret"`
}

//CreateToken authenticates the user
func (t *Token) CreateToken(c *fiber.Ctx, userID uint, secret string, expire ...int64) (*Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	if len(expire) > 0 {
		t.Expire = expire[0]
	} else {
		t.Expire = 3600
	}
	expiresIn := time.Now().Add(time.Duration(t.Expire) * time.Second).Unix()
	claims["exp"] = expiresIn

	tokenHash, err := token.SignedString([]byte(secret))

	if err != nil {
		return nil, err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "Verify-Rest-Token",
		Value:    tokenHash,
		Secure:   false,
		HTTPOnly: true,
	})
	t.Hash = tokenHash
	t.Expire = expiresIn
	return t, nil
}

//ParseToken returns the users id or error
func (t *Token) ParseToken(c *fiber.Ctx, secret string) (uint, error) {
	tokenString := c.Cookies("Verify-Rest-Token")

	if tokenString == "" {
		return 0, errors.New("Empty auth cookie")
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	//Checks if the token is valid if it is not then it deletes it
	err2 := claims.Valid()

	if err2 != nil {
		t.DeleteToken(c)
		return 0, err2
	}

	return uint(claims["id"].(float64)), nil
}

//DeleteToken deletes the jwt token
func (t *Token) DeleteToken(c *fiber.Ctx) {
	c.ClearCookie("Verify-Rest-Token")
}

//RefreshToken refreshes the token
func (t *Token) RefreshToken(c *fiber.Ctx, secret string) (*Token, error) {
	u, err := t.ParseToken(c, secret)

	if err != nil {
		return nil, nil
	}

	return t.CreateToken(c, u, secret)
}
