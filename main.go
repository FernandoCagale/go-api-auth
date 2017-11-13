package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func mainUser(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	fmt.Println("User Name:", claims["name"], "User ID:", claims["jti"])

	return c.JSON(http.StatusOK, map[string]string{
		"message": "informations",
	})
}

func login(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Bad Request",
		})
	}

	if u.Username == "jack" && u.Password == "1234" {

		token, err := createJwtToken()
		if err != nil {
			return c.String(http.StatusInternalServerError, "something went wrong")
		}

		return c.JSON(http.StatusOK, map[string]string{
			"access_token": token,
			"expires_in":   "21600",
			"token_type":   "Bearer",
		})
	}
	return c.JSON(http.StatusUnauthorized, map[string]string{
		"message": "Your username or password were wrong",
	})
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"jack",
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func main() {
	e := echo.New()

	userGroup := e.Group("/user")

	userGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("secret"),
	}))

	userGroup.GET("", mainUser)

	e.POST("/login", login)

	e.Logger.Fatal(e.Start(":8000"))
}
