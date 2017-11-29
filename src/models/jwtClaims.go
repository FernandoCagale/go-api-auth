package models

import jwt "github.com/dgrijalva/jwt-go"

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}
