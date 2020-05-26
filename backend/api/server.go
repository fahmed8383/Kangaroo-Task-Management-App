package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	mux.NewRouter()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "TEST Welcome to my website!")
	})
	http.ListenAndServe(":6060", nil)
}
