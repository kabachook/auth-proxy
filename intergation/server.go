package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func loggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("%s %s", r.Method, r.URL.EscapedPath())
	}

	return http.HandlerFunc(fn)
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func jwtKey() []byte {
	return []byte(os.Getenv("JWT_TOKEN"))
}

type HTTPError struct {
	Code    int
	Message string
	Err     error
}

func (e *HTTPError) Error() string { return string(e.Code) + ": " + e.Message }
func (e *HTTPError) Unwrap() error { return e.Err }

func extractJWTFromHeader(r *http.Request) string {
	h := r.Header.Get("Authorization")
	strArr := strings.Split(h, " ")
	if len(strArr) == 2 && strArr[0] == "Bearer" {
		return strArr[1]
	}
	return ""
}
func parseJWTToken(t string) (*Claims, *HTTPError) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, &HTTPError{
				Message: "Invalid signature",
				Code:    http.StatusUnauthorized,
				Err:     err,
			}
		}
		return nil, &HTTPError{
			Message: "Bad request",
			Code:    http.StatusBadRequest,
			Err:     err,
		}
	}
	if !tkn.Valid {
		return nil, &HTTPError{
			Message: "Token is not valid",
			Code:    http.StatusUnauthorized,
			Err:     err,
		}
	}
	return claims, nil
}

func checkAuth(r *http.Request) (*Claims, *HTTPError) {
	token := extractJWTFromHeader(r)
	if token == "" {
		return nil, &HTTPError{
			Message: "Can't extract token from header",
			Code:    http.StatusUnauthorized,
		}
	}
	return parseJWTToken(token)
}

func main() {
	_, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		os.Setenv("JWT_KEY", "qwerty")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Hello!\nYou called %s", r.URL.EscapedPath())))
	})
	mux.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		claims, err := checkAuth(r)
		if err != nil {
			w.WriteHeader(err.Code)
			w.Write([]byte(err.Message))
			return
		}

		w.Write([]byte(fmt.Sprintf("Welcome, %s!", claims.Username)))
	})
	mux.HandleFunc("/issue", func(w http.ResponseWriter, r *http.Request) {
		usernames, ok := r.URL.Query()["username"]
		if !ok || len(usernames[0]) < 1 {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("username is not provided"))
		}

		username := usernames[0]
		claims := &Claims{
			Username: username,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString(jwtKey())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(tokenStr))
	})

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", loggingMiddleware(mux)))
}
