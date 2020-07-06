package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <username>\nIssues JWT token for <username>\nJWT_KEY env var is used as secret", os.Args[0])
		os.Exit(1)
	}

	username := os.Args[1]
	jwtKeyStr, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		jwtKeyStr = "qwerty"
	}
	jwtKey := []byte(jwtKeyStr)
	claims := &Claims{
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
