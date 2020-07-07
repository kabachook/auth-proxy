package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	flag "github.com/spf13/pflag"
)

var (
	addr string
)

func init() {
	flag.StringVarP(&addr, "addr", "l", ":8080", "bind address")
	flag.Parse()
}

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var username string
		usernameArr, ok := r.Header["X-Username"]
		if ok {
			username = usernameArr[0]
		}
		resp := fmt.Sprintf(
			"Hello!\nYou called %s\nHeaders: %+v\nUsername: %+v",
			r.URL.EscapedPath(),
			r.Header,
			username)
		w.Write([]byte(resp))
	})

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(addr, handlers.LoggingHandler(os.Stdout, router)))
}
