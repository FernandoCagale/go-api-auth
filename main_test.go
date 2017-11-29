package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var (
	MainJson = `{"message":"informations"}`
	Token    = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiamFjayIsImV4cCI6MTUwNjE5NTU1MSwianRpIjoibWFpbl91c2VyX2lkIn0.KbtO8K2R78_z1ZxU4IRf5b1_zmqKwonTcRdriXN1gJlFhBjh-q675tlDHjJjDmuUoBUc9Ap8umpfcR3hGbU-XQ"
)

type Login struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func TestLogin(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(`{"username":"jack","password":"1234"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		login := Login{}

		err := json.Unmarshal([]byte(rec.Body.String()), &login)

		assert.Nil(t, err)
		assert.NotNil(t, login.AccessToken)
		assert.Equal(t, login.ExpiresIn, "21600")
		assert.Equal(t, login.TokenType, "Bearer")
	}
}

func TestMainUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/user", nil)
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", Token))
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, mainUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, MainJson, rec.Body.String())
	}
}
