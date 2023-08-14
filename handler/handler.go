package handler

import (
	"fmt"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "sing up")
	//TODO: implement me!
}

func HandleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from server!")
}
