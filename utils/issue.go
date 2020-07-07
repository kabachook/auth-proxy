package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	flag "github.com/spf13/pflag"
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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Issue an JWT token\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVarP(&username, "username", "u", "jdoe", "token username")
	flag.DurationVarP(&duration, "duration", "d", time.Minute*10, "token duration")
	flag.StringVarP(&jwtKeyStr, "key", "k", "qwerty", "JWT key")
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
