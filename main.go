package main

import (
	"fmt"
	"net/http"
)

func homehandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to picbucket</h1>")
}

func notFoundhandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.WriteHeader(http.StatusNotFound)
	// fmt.Fprint(w, "<h1>404 - page not found</h1>")
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
		<h1>Contact Page</h1>
		<p>
			You can get in touch with me at 
				<a 
					href="mailto:altamashattari786@gmail.com"
				>
				altamashattari786@gmail.com
			</a>
		</p>
	`)
}

func pathhandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homehandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		notFoundhandler(w, r)
	}
}

// type Router struct{}

// func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":
// 		homehandler(w, r)
// 	case "/contact":
// 		contactHandler(w, r)
// 	default:
// 		notFoundhandler(w, r)
// 	}
// }

func main() {
	// http.HandleFunc("/", homehandler)
	// http.HandleFunc("/contact", contactHandler)
	// http.HandleFunc("/", pathhandler)
	// var router Router
	fmt.Println("Starting the server at port :8080")
	http.ListenAndServe(":8080", http.HandlerFunc(pathhandler))
}
