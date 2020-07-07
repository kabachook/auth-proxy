package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var (
	username  string
	duration  time.Duration
	jwtKeyStr string
)

func init() {
	flag.StringVar(&username, "username", "jdoe", "token username")
	flag.DurationVar(&duration, "duration", time.Minute*10, "token duration")
	flag.StringVar(&jwtKeyStr, "jwt_key", "qwerty", "JWT key")
	flag.Parse()
}

func main() {
	jwtKey := []byte(jwtKeyStr)
	claims := &Claims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(tokenStr)
}
