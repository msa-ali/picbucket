package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my awesome website</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Starting the server at port :8080")
	http.ListenAndServe(":8080", nil)
}
