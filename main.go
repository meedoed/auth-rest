package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/meedoed/auth-rest/handler"
	"github.com/meedoed/auth-rest/repository"
)

func main() {

	db, err := repository.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = db.Disconnect(context.TODO())
	}()

	r := mux.NewRouter()

	r.HandleFunc("/hello", handler.HandleHello)
	r.HandleFunc("/signup", handler.Signup)

	http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r))
}
