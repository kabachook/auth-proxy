package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/kabachook/auth-proxy/pkg/middlewares/jwt"
)

func main() {
	username := os.Args[1]
	jwtKeyStr, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		jwtKeyStr = "qwerty"
	}
	jwtKey := []byte(jwtKeyStr)
	claims := &jwtmiddleware.Claims{
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(tokenStr)
}
